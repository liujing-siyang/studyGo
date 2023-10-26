package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

var (
	cookieUrl   = "https://www.amazon.com.au/gp/bestsellers/?ref_=nav_cs_bestsellers"
	baseUrl     string
	isSecondary bool
	pageIndex   int
	auSelectUrl = "https://search.ipaustralia.gov.au/trademarks/search/count/quick?"
	auParseUrl  *url.URL
)

type Res struct {
	Count int    `json:"count"`
	Err   string `json:"errors"`
}

func init() {
	flag.StringVar(&baseUrl, "b", "https://www.amazon.com.au/gp/bestsellers/beauty/ref=zg_bs_nav_beauty_0", "输入需要爬取的Url")
	flag.BoolVar(&isSecondary, "s", false, "是否需要爬取二级类目下的所有商品")
	flag.IntVar(&pageIndex, "p", 5, "类目需要爬取的总页数，最大8页")
}

func main() {
	flag.Parse()
	validURL := govalidator.IsURL(baseUrl)
	if !validURL {
		fmt.Println("需要爬取的url无效")
		return
	}
	parseUrl, err := url.Parse(auSelectUrl)
	if err != nil {
		fmt.Println("澳洲品牌查询url解析失败")
	}
	auParseUrl = parseUrl
	fmt.Printf("爬取的url为: %s \n", baseUrl)
	if isSecondary {
		fmt.Println("将爬取一级类目所有二级类目下的所有商品")
		fmt.Println("5秒后开始爬取，如需终止，请按ctrl + c")
		time.Sleep(5 * time.Second)
		fmt.Println("开始爬取")
		SecondaryCategory()
	} else {
		fmt.Println("将爬取一级类目下的所有商品")
		fmt.Println("5秒后开始爬取，如需终止，请按ctrl + c")
		time.Sleep(5 * time.Second)
		fmt.Println("开始爬取")
		PrimaryCategory()
	}
	fmt.Println("爬取结束")
}

func PrimaryCategory() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PrimaryCategory", r)
		}
	}()
	c := colly.NewCollector(
		// 只爬取该域名下
		// colly.AllowedDomains("https://www.amazon.com"),
		// // 忽略爬虫规则
		colly.IgnoreRobotsTxt(),
		colly.MaxDepth(1),
	)
	// 设置代理和refer
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
	})
	cookie := c.Cookies(cookieUrl)
	// 自动设置cookie
	c.SetCookies("", cookie)
	secondC := c.Clone()
	pageC := c.Clone()
	thirdC := secondC.Clone()
	// c.OnHTML("div[role='treeitem']", func(e *colly.HTMLElement) {
	// 	dom := e.DOM
	// 	hrefSelect := dom.Find("a")
	// 	href, _ := hrefSelect.Attr("href")
	// 	// fmt.Printf("ret: %s\n", href)
	// 	if len(href) == 0 {
	// 		return
	// 	}
	// 	absoluteURL := e.Request.AbsoluteURL(href)
	// 	// fmt.Printf("First Url: %s\n", absoluteURL)
	// 	// que.AddURL(absoluteURL)
	// 	secondC.Visit(absoluteURL)

	// })

	secondC.OnHTML("div.a-text-center > .a-pagination > .a-selected ", func(e *colly.HTMLElement) {
		dom := e.DOM
		a := dom.Find("a")
		link, _ := a.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)
		// fmt.Printf("AbsoluteURL: %s\n", absoluteURL)
		base := absoluteURL[:len(absoluteURL)-1]
		for i := 1; i <= pageIndex; i++ {
			url := fmt.Sprintf("%s%d", base, i)
			// fmt.Println(url)
			// que.AddURL(url)
			pageC.Visit(url)
		}

	})

	pageC.OnHTML(".a-link-normal", func(e *colly.HTMLElement) {
		if len(e.Text) == 0 {
			return
		}
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)
		thirdC.Visit(absoluteURL)
	})

	// que.Run(thirdC)
	thirdC.OnHTML("div[id='dp-container']", func(e *colly.HTMLElement) {
		dom := e.DOM
		// span := dom.Find("a")
		span := dom.Find("a[id='bylineInfo']")
		// text := e.ChildText("span.a-size-base")
		brand := span.Text()
		brand = strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(brand, "Brand: "), "Visit the "), " Store")
		// fmt.Printf("brand: %s\n", brand)
		q := auParseUrl.Query()
		q.Add("q", brand)
		queryStr := q.Encode()
		url := fmt.Sprintf("%s%s", auSelectUrl, queryStr)
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(data))
		myres := Res{}
		err = json.Unmarshal(data, &myres)
		if err == nil && myres.Count == 0 {
			var asin string
			dom.Find("#prodDetails .prodDetAttrValue").Each(func(i int, s *goquery.Selection) {
				text := s.Text()
				if strings.Contains(text, "B0") {
					asin = strings.TrimLeft(strings.TrimSpace(text),"‎") 
					return
				}
			})
			fmt.Printf("brand: %s ; asin: %s \n", brand, asin)
		}
		// fmt.Printf("brand: %s\n", e.Text)
	})

	secondC.Visit(baseUrl)
}

func SecondaryCategory() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in SecondaryCategory", r)
		}
	}()
	c := colly.NewCollector(
		// 只爬取该域名下
		// colly.AllowedDomains("https://www.amazon.com"),
		// // 忽略爬虫规则
		colly.IgnoreRobotsTxt(),
		// colly.Async(true),
		colly.MaxDepth(1),
	)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
	})
	cookie := c.Cookies(cookieUrl)
	// 自动设置cookie
	c.SetCookies("", cookie)
	secondC := c.Clone()
	pageC := c.Clone()
	thirdC := secondC.Clone()

	c.OnHTML("div[role='treeitem']", func(e *colly.HTMLElement) {
		dom := e.DOM
		hrefSelect := dom.Find("a")
		href, _ := hrefSelect.Attr("href")
		// fmt.Printf("ret: %s\n", href)
		if len(href) == 0 {
			return
		}
		absoluteURL := e.Request.AbsoluteURL(href)
		// fmt.Printf("AbsoluteURL: %s\n", absoluteURL)
		secondC.Visit(absoluteURL)
	})

	secondC.OnHTML("div.a-text-center > .a-pagination > .a-selected ", func(e *colly.HTMLElement) {
		dom := e.DOM
		a := dom.Find("a")
		link, _ := a.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)
		// fmt.Printf("AbsoluteURL: %s\n", absoluteURL)
		base := absoluteURL[:len(absoluteURL)-1]
		for i := 1; i <= pageIndex; i++ {
			url := fmt.Sprintf("%s%d", base, i)
			// fmt.Println(url)
			// que.AddURL(url)
			pageC.Visit(url)
		}

	})

	pageC.OnHTML(".a-link-normal", func(e *colly.HTMLElement) {
		if len(e.Text) == 0 {
			return
		}
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)
		// fmt.Printf("AbsoluteURL: %s\n", absoluteURL)
		thirdC.Visit(absoluteURL)
	})

	// que.Run(thirdC)
	thirdC.OnHTML("div[id='dp-container']", func(e *colly.HTMLElement) {
		dom := e.DOM
		// span := dom.Find("a")
		span := dom.Find("a[id='bylineInfo']")
		// text := e.ChildText("span.a-size-base")
		brand := span.Text()
		brand = strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(brand, "Brand: "), "Visit the "), " Store")
		// fmt.Printf("brand: %s\n", brand)
		q := auParseUrl.Query()
		q.Add("q", brand)
		queryStr := q.Encode()
		url := fmt.Sprintf("%s%s", auSelectUrl, queryStr)
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		myres := Res{}
		err = json.Unmarshal(data, &myres)
		if err == nil && myres.Count == 0 {
			var asin string
			dom.Find("#prodDetails .prodDetAttrValue").Each(func(i int, s *goquery.Selection) {
				text := s.Text()
				if strings.Contains(text, "B0") {
					asin = strings.TrimLeft(strings.TrimSpace(text),"‎") 
					return
				}
			})
			fmt.Printf("brand: %s ; asin: %s \n", brand, asin)
		}
		// fmt.Printf("brand: %s\n", e.Text)
	})
	c.Visit(baseUrl)
}

