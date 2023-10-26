package brand

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/queue"
)

func test1() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.baidu.com"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Response %s: %d bytes\n", r.Request.URL, len(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error %s: %v\n", r.Request.URL, err)
	})

	c.Visit("http://www.baidu.com/")
}

type Repository struct {
	Author  string
	Name    string
	Link    string
	Desc    string
	Lang    string
	Stars   int
	Forks   int
	Add     int
	BuiltBy []string
}

func test2() {
	c := colly.NewCollector(
		colly.MaxDepth(1),
	)

	repos := make([]*Repository, 0, 15)
	c.OnHTML("article.Box-row", func(e *colly.HTMLElement) {
		repo := &Repository{}

		// author & repository name
		authorRepoName := e.ChildText("h2.h3.lh-condensed")
		parts := strings.Split(authorRepoName, "/")
		repo.Author = strings.TrimSpace(parts[0])
		repo.Name = strings.TrimSpace(parts[1])

		// link
		repo.Link = e.Request.AbsoluteURL(e.ChildAttr("h2.h3.lh-condensed", "href"))

		// description
		repo.Desc = e.ChildText("p.pr-4")

		// language
		repo.Lang = strings.TrimSpace(e.ChildText("div.mt-2 > span.mr-3 > span[itemprop]"))

		// star & fork
		starForkStr := e.ChildText("div.mt-2 > a.mr-3")
		starForkStr = strings.Replace(strings.TrimSpace(starForkStr), ",", "", -1)
		parts = strings.Split(starForkStr, "\n")
		repo.Stars, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
		repo.Forks, _ = strconv.Atoi(strings.TrimSpace(parts[len(parts)-1]))

		// add
		addStr := e.ChildText("div.mt-2 > span.float-sm-right")
		parts = strings.Split(addStr, " ")
		fmt.Println(strings.Replace(strings.TrimSpace(parts[0]), ",", "", -1))
		repo.Add, _ = strconv.Atoi(parts[0])

		// built by
		e.ForEach("div.mt-2 > span.mr-3  img[src]", func(index int, img *colly.HTMLElement) {
			repo.BuiltBy = append(repo.BuiltBy, img.Attr("src"))
		})

		repos = append(repos, repo)
	})

	c.Visit("https://github.com/trending")

	fmt.Printf("%d repositories\n", len(repos))
	fmt.Println("first repository:")
	for _, repo := range repos {
		// fmt.Println("Author:", repo.Author)
		// fmt.Println("Name:", repo.Name)
		fmt.Printf("%+v\n", repo)
		break
	}
}

var (
	FirstUrl  = "https://www.amazon.com/gp/bestsellers"
	secondUrl = "https://www.amazon.com/Best-Sellers-Amazon-Devices-Accessories/zgbs/amazon-devices/ref=zg_bs_nav_amazon-devices_0"
	thirdUrl  = "https://www.amazon.com/Instant-Pot-Electric-Pressure-Stainless/dp/B096G7XXN2/ref=zg_bs_g_amazon-renewed_sccl_1/137-0061650-5417576?th=1"
)

type Res struct {
	Count int    `json:"count"`
	Err   string `json:"errors"`
}

// 一级
func test3() {
	c := colly.NewCollector(
	// colly.AllowedDomains("https://www.amazon.com"),
	// colly.IgnoreRobotsTxt(),
	)
	// 设置用户代理，模拟浏览器
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	cookie := c.Cookies(FirstUrl)
	// 自动设置cookie
	c.SetCookies("", cookie)
	que, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)
	// c.OnRequest(func(r *colly.Request) {
	// 	// 手动设置cookie
	// 	r.Headers.Add("cookie",`csm-sid=168-8407135-6138251; x-amz-captcha-1=1694877773730581; x-amz-captcha-2=af5gB0yUFjly5ibu35fP4g==; session-id=137-0061650-5417576; session-id-time=2082787201l; i18n-prefs=USD; sp-cdn="L5Z9:SG"; skin=noskin; ubid-main=131-6500085-9846561; ld=AZUSSOA-sell; s_pers=%20s_fid%3D54FD015DF5F8CE09-097B846A5CD3DC37%7C1852723626710%3B%20s_dl%3D1%7C1694872626711%3B%20gpv_page%3DUS%253AAZ%253ASOA-overview-sell%7C1694872626714%3B%20s_ev15%3D%255B%255B%2527AZUSSOA-sell%2527%252C%25271694870826715%2527%255D%255D%7C1852723626715%3B; s_sess=%20s_cc%3Dtrue%3B%20s_ppvl%3DUS%25253AAZ%25253ASOA-overview-sell%252C6%252C6%252C715%252C1536%252C715%252C1536%252C864%252C1.25%252CL%3B%20c_m%3DAZUSSOA-sellundefinedAmazon.comundefined%3B%20s_ppv%3DUS%25253AAZ%25253ASOA-overview-sell%252C6%252C6%252C715%252C1182%252C715%252C1536%252C864%252C1.25%252CL%3B; session-token=wxQmRNYNRS1CAfB76Xs9gyWvEryDgW+osvj4Mo5KKMD6Ik+ugR/BXQfCcBwYGIpSgBXMfAZ0aJz2a+1D3A5CZDAgD1BL1oQB6S/VyuoUBpcuaIlHx1zY2SW0H6JYyVOsKfEvyjFaLT8QWtlgr5Az5sdDvJWhIwIiuFuIAyqiDrZbb6D/q9LBh0TlCD0mg236m6Hkq7ljKhgYEGjfP5FkZLcndf0i54NJi82KYLirQOModLeeYuUhLxWFm7928w6VHuQsawrUAbrn1LlCPfs58QW6xK59UwdwmAimPGjRAaVfRRgfpVjlUhdNQN0xHsLdbbWuLke9Ot8w+BIiyDV060nXDsK16TFf; csm-hit=tb:5YN0SH522YBYS3T2YGV2+s-QBR5Q038NZHNTYMWRKKP|1694916106054&t:1694916106054&adb:adblk_yes`)
	// })
	c.OnHTML("div[role='treeitem']", func(e *colly.HTMLElement) {
		dom := e.DOM
		hrefSelect := dom.Find("a")
		href, _ := hrefSelect.Attr("href")
		fmt.Printf("ret: %s\n", href)
		if len(href) == 0 {
			return
		}
		absoluteURL := e.Request.AbsoluteURL(href)
		fmt.Printf("AbsoluteURL: %s\n", absoluteURL)
		// c.Visit(absoluteURL)
		// link := e.Attr("href")
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
		que.AddURL(absoluteURL)
	})

	// c.OnHTML("",func(h *colly.HTMLElement) {

	// })
	que.Run(c)
	c.Visit(FirstUrl)
}

// 二级商品
func test4() {
	c := colly.NewCollector()
	// 设置用户代理，模拟浏览器
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	cookie := c.Cookies(secondUrl)
	// 自动设置cookie
	c.SetCookies("", cookie)
	c.OnHTML(".a-link-normal", func(e *colly.HTMLElement) {
		if len(e.Text) == 0 {
			return
		}
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)
		fmt.Printf("AbsoluteURL: %s\n", absoluteURL)
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
	})
	c.Visit(secondUrl)
}

// 品牌名
func test5() {
	c := colly.NewCollector()
	// 设置用户代理，模拟浏览器
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	cookie := c.Cookies(thirdUrl)
	// 自动设置cookie
	c.SetCookies("", cookie)
	// c.OnHTML("div[id='bylineInfo_feature_div']", func(e *colly.HTMLElement) {
	// 	dom := e.DOM
	// 	span := dom.Find("a")
	// 	// text := e.ChildText("span.a-size-base")
	// 	fmt.Printf("brand: %s\n", span.Text())
	// 	fmt.Printf("brand: %s\n", e.Text)
	// 	// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	// 	// c.Visit(e.Request.AbsoluteURL(link))
	// })
	c.OnHTML("div[id='bylineInfo_feature_div']", func(e *colly.HTMLElement) {
		dom := e.DOM
		// span := dom.Find("a")
		span := dom.Find("a[id='bylineInfo']")
		// text := e.ChildText("span.a-size-base")
		brand := span.Text()
		brand = strings.TrimSuffix(strings.TrimPrefix(brand, "Visit the "), " Store")
		fmt.Printf("brand: %s\n", brand)
		// fmt.Printf("brand: %s\n", e.Text)

		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
	})
	c.Visit(thirdUrl)
}

// 品牌注册
func test6() {
	url := "https://search.ipaustralia.gov.au/trademarks/search/count/quick?q=Dynamite"
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
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%+v", myres)
	}
}

// 二级分页url
func test7() {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	cookie := c.Cookies(secondUrl)
	// 自动设置cookie
	c.SetCookies("", cookie)
	page := c.Clone()
	page.OnHTML(".a-link-normal", func(e *colly.HTMLElement) {
		if len(e.Text) == 0 {
			return
		}
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)
		fmt.Printf("AbsoluteURL: %s\n", absoluteURL)
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
	})
	// que, _ := queue.New(
	// 	2, // Number of consumer threads
	// 	&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	// )
	c.OnHTML("div.a-text-center > .a-pagination > .a-selected ", func(e *colly.HTMLElement) {
		dom := e.DOM
		a := dom.Find("a")
		link, _ := a.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)
		// fmt.Printf("AbsoluteURL: %s\n", absoluteURL)
		base := absoluteURL[:len(absoluteURL)-1]
		for i := 1; i < 10; i++ {
			url := fmt.Sprintf("%s%d", base, i)
			// fmt.Println(url)
			// que.AddURL(url)
			page.Visit(url)
		}

	})

	// que.Run(page)
	c.Visit(secondUrl)

}

// <div id = div1 > </div>
// #div1
// <div class = "div1" > </div>
// .div1

func test8() {
	c := colly.NewCollector(
		// 只爬取该域名下
		// colly.AllowedDomains("https://www.amazon.com"),
		// // 忽略爬虫规则
		colly.IgnoreRobotsTxt(),
		colly.Async(true),
		colly.MaxDepth(1),
	)
	extensions.RandomUserAgent(c)
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	// 设置用户代理，模拟浏览器
	// c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	cookie := c.Cookies(FirstUrl)
	// 自动设置cookie
	c.SetCookies("", cookie)
	secondC := c.Clone()
	thirdC := secondC.Clone()
	// que, _ := queue.New(
	// 	3, // Number of consumer threads
	// 	&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	// )

	secondC.Limit(&colly.LimitRule{
		// DomainRegexp: `unsplash\.com`,
		RandomDelay: 500 * time.Millisecond,
		Parallelism: 10,
	})
	thirdC.Limit(&colly.LimitRule{
		// DomainRegexp: `unsplash\.com`,
		RandomDelay: 500 * time.Millisecond,
		Parallelism: 50,
	})
	c.OnHTML("div[role='treeitem']", func(e *colly.HTMLElement) {
		dom := e.DOM
		hrefSelect := dom.Find("a")
		href, _ := hrefSelect.Attr("href")
		// fmt.Printf("ret: %s\n", href)
		if len(href) == 0 {
			return
		}
		absoluteURL := e.Request.AbsoluteURL(href)
		// fmt.Printf("First Url: %s\n", absoluteURL)
		// que.AddURL(absoluteURL)
		secondC.Visit(absoluteURL)

	})

	secondC.OnHTML(".a-link-normal", func(e *colly.HTMLElement) {
		if len(e.Text) == 0 {
			return
		}
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)
		// fmt.Printf("Second URL: %s\n", absoluteURL)
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
		// que.AddURL(absoluteURL)
		thirdC.Visit(absoluteURL)
	})

	// que.Run(thirdC)
	thirdC.OnHTML("div[id='bylineInfo_feature_div']", func(e *colly.HTMLElement) {
		dom := e.DOM
		// span := dom.Find("a")
		span := dom.Find("a[id='bylineInfo']")
		// text := e.ChildText("span.a-size-base")
		brand := span.Text()
		brand = strings.TrimSuffix(strings.TrimPrefix(brand, "Visit the "), " Store")
		// fmt.Printf("brand: %s\n", brand)
		url := fmt.Sprintf("https://search.ipaustralia.gov.au/trademarks/search/count/quick?q=%s", brand)
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
			fmt.Printf("brand: %s\n", brand)
		}
		// fmt.Printf("brand: %s\n", e.Text)

		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
	})
	// que.Run(secondC)
	c.Visit(FirstUrl)
	c.Wait()
	secondC.Wait()
	thirdC.Wait()
}

func test9() {
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
	// 设置用户代理，模拟浏览器
	// c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	cookie := c.Cookies(FirstUrl)
	// 自动设置cookie
	c.SetCookies("", cookie)
	secondC := c.Clone()
	pageC := c.Clone()
	thirdC := secondC.Clone()
	// que, _ := queue.New(
	// 	3, // Number of consumer threads
	// 	&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	// )

	// c.Limit(&colly.LimitRule{
	// 	// DomainRegexp: `unsplash\.com`,
	// 	RandomDelay: 1000 * time.Millisecond,
	// 	Parallelism: 2,
	// })
	// secondC.Limit(&colly.LimitRule{
	// 	// DomainRegexp: `unsplash\.com`,
	// 	RandomDelay: 1000 * time.Millisecond,
	// 	Parallelism: 2,
	// })
	// pageC.Limit(&colly.LimitRule{
	// 	// DomainRegexp: `unsplash\.com`,
	// 	RandomDelay: 1000 * time.Millisecond,
	// 	Parallelism: 10,
	// })
	// thirdC.Limit(&colly.LimitRule{
	// 	// DomainRegexp: `unsplash\.com`,
	// 	RandomDelay: 1000 * time.Millisecond,
	// 	Parallelism: 50,
	// })
	c.OnHTML("div[role='treeitem']", func(e *colly.HTMLElement) {
		dom := e.DOM
		hrefSelect := dom.Find("a")
		href, _ := hrefSelect.Attr("href")
		// fmt.Printf("ret: %s\n", href)
		if len(href) == 0 {
			return
		}
		absoluteURL := e.Request.AbsoluteURL(href)
		// fmt.Printf("First Url: %s\n", absoluteURL)
		// que.AddURL(absoluteURL)
		secondC.Visit(absoluteURL)

	})

	secondC.OnHTML("div.a-text-center > .a-pagination > .a-selected ", func(e *colly.HTMLElement) {
		dom := e.DOM
		a := dom.Find("a")
		link, _ := a.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)
		// fmt.Printf("AbsoluteURL: %s\n", absoluteURL)
		base := absoluteURL[:len(absoluteURL)-1]
		for i := 8; i < 11; i++ {
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
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
		thirdC.Visit(absoluteURL)
	})

	// que.Run(thirdC)
	thirdC.OnHTML("div[id='bylineInfo_feature_div']", func(e *colly.HTMLElement) {
		dom := e.DOM
		// span := dom.Find("a")
		span := dom.Find("a[id='bylineInfo']")
		// text := e.ChildText("span.a-size-base")
		brand := span.Text()
		brand = strings.TrimSuffix(strings.TrimPrefix(brand, "Visit the "), " Store")
		// fmt.Printf("brand: %s\n", brand)
		url := fmt.Sprintf("https://search.ipaustralia.gov.au/trademarks/search/count/quick?q=%s", brand)
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
			fmt.Printf("brand: %s\n", brand)
		}
		// fmt.Printf("brand: %s\n", e.Text)

		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
	})
	// que.Run(secondC)
	c.Visit(FirstUrl)
	// c.Wait()
	// secondC.Wait()
	// pageC.Wait()
	// thirdC.Wait()
}

// 指定一级类下所有的商品
func test10() {
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
	// 设置用户代理，模拟浏览器
	// c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	cookie := c.Cookies(FirstUrl)
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
		for i := 1; i < 9; i++ {
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
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
		thirdC.Visit(absoluteURL)
	})

	// que.Run(thirdC)
	thirdC.OnHTML("div[id='bylineInfo_feature_div']", func(e *colly.HTMLElement) {
		dom := e.DOM
		// span := dom.Find("a")
		span := dom.Find("a[id='bylineInfo']")
		// text := e.ChildText("span.a-size-base")
		brand := span.Text()
		brand = strings.TrimSuffix(strings.TrimPrefix(brand, "Visit the "), " Store")
		// fmt.Printf("brand: %s\n", brand)
		url := fmt.Sprintf("https://search.ipaustralia.gov.au/trademarks/search/count/quick?q=%s", brand)
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
			fmt.Printf("brand: %s\n", brand)
		}
		// fmt.Printf("brand: %s\n", e.Text)

		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
	})
	// que.Run(secondC)
	secondC.Visit("https://www.amazon.com.au/gp/bestsellers/toys/ref=zg_bs_nav_toys_0")
	// c.Wait()
	// secondC.Wait()
	// pageC.Wait()
	// thirdC.Wait()
}

func test11() {
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
	// 设置用户代理，模拟浏览器
	// c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	cookie := c.Cookies(FirstUrl)
	// 自动设置cookie
	c.SetCookies("", cookie)
	c.OnHTML("div[role='treeitem']", func(e *colly.HTMLElement) {
		dom := e.DOM
		hrefSelect := dom.Find("a")
		href, _ := hrefSelect.Attr("href")
		// fmt.Printf("ret: %s\n", href)
		if len(href) == 0 {
			return
		}
		absoluteURL := e.Request.AbsoluteURL(href)
		fmt.Printf("AbsoluteURL: %s\n", absoluteURL)
		// c.Visit(absoluteURL)
		// link := e.Attr("href")
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
	})
	c.Visit("https://www.amazon.com/Best-Sellers-Pet-Supplies/zgbs/pet-supplies/ref=zg_bs_nav_pet-supplies_0")
}

// 指定一级类目下所有的二级类目
func test12() {
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
	// 设置用户代理，模拟浏览器
	// c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	cookie := c.Cookies(FirstUrl)
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
		for i := 1; i < 9; i++ {
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
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
		thirdC.Visit(absoluteURL)
	})

	// que.Run(thirdC)
	thirdC.OnHTML("div[id='bylineInfo_feature_div']", func(e *colly.HTMLElement) {
		dom := e.DOM
		// span := dom.Find("a")
		span := dom.Find("a[id='bylineInfo']")
		// text := e.ChildText("span.a-size-base")
		brand := span.Text()
		brand = strings.TrimSuffix(strings.TrimPrefix(brand, "Visit the "), " Store")
		// fmt.Printf("brand: %s\n", brand)
		url := fmt.Sprintf("https://search.ipaustralia.gov.au/trademarks/search/count/quick?q=%s", brand)
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
			fmt.Printf("%s\n", brand)
		}
		// fmt.Printf("brand: %s\n", e.Text)

		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// c.Visit(e.Request.AbsoluteURL(link))
	})
	// que.Run(secondC)
	c.Visit("https://www.amazon.com.au/gp/bestsellers/beauty/ref=zg_bs_nav_beauty_0")
	// c.Wait()
	// secondC.Wait()
	// pageC.Wait()
	// thirdC.Wait()
}


func test13() {
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
	cookie := c.Cookies(FirstUrl)
	// 自动设置cookie
	c.SetCookies("", cookie)
	c.OnHTML("div[id='dp-container']", func(e *colly.HTMLElement) {
		dom := e.DOM
		// span := dom.Find("a")
		span := dom.Find("a[id='bylineInfo']")
		// text := e.ChildText("span.a-size-base")
		brand := span.Text()
		brand = strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(brand, "Brand: "), "Visit the "), " Store")
		fmt.Printf("brand: %s\n", brand)
		var asin string
		dom.Find("#productDetails_feature_div .prodDetAttrValue").Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if strings.Contains(text, "B0") {
				asin = text
			}
		})
		fmt.Printf("asin: %s\n", asin)
	})
	c.Visit("https://www.amazon.com.au/MamaBabyCo%C2%AE-Silver-Nursing-Cups-Breastfeeding/dp/B0BWDYVYNZ/ref=sr_1_1_sspa?keywords=MAMABABYCO&qid=1695948582&sr=8-1-spons&sp_csd=d2lkZ2V0TmFtZT1zcF9hdGY&th=1")
}