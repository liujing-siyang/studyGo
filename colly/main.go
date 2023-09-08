package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

// 选择器示例https://cloud.tencent.com/developer/article/1196783

func main() {
	test5()
}

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
	c.OnHTML(".Box .Box-row", func(e *colly.HTMLElement) {
		repo := &Repository{}

		// author & repository name
		authorRepoName := e.ChildText("h2.h3 > a")
		parts := strings.Split(authorRepoName, "/")
		repo.Author = strings.TrimSpace(parts[0])
		repo.Name = strings.TrimSpace(parts[1])

		// link
		repo.Link = e.Request.AbsoluteURL(e.ChildAttr("h2.h3 >a", "href"))

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
		repo.Add, _ = strconv.Atoi(strings.Replace(strings.TrimSpace(parts[0]), ",", "", -1))

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

type Hot struct {
	Rank   string `selector:"a > div.index_1Ew5p"`
	Name   string `selector:"div.content_1YWBm > a.title_dIF3B" ` // > div.c-single-text-ellipsis
	Author string `selector:"div.content_1YWBm > div.intro_1l0wp:nth-child(2)"`
	Type   string `selector:"div.content_1YWBm > div.intro_1l0wp:nth-child(3)"`
	Desc   string `selector:"div.desc_3CTjT"`
}

func test3() {

	c := colly.NewCollector()
	hots := make([]*Hot, 0)
	c.OnHTML("div.category-wrap_iQLoo", func(e *colly.HTMLElement) {
		hot := &Hot{}

		err := e.Unmarshal(hot)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		hots = append(hots, hot)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting:", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response:", len(r.Body))
	})

	err := c.Visit("https://top.baidu.com/board?tab=novel")
	if err != nil {
		fmt.Println("Visit error:", err)
		return
	}

	fmt.Printf("%d hots\n", len(hots))
	for _, hot := range hots {
		fmt.Println("first hot:")
		fmt.Println("Rank:", hot.Rank)
		fmt.Println("Name:", hot.Name)
		fmt.Println("Author:", hot.Author)
		fmt.Println("Type:", hot.Type)
		fmt.Println("Desc:", hot.Desc)
		break
	}
}

func test4() {
	c1 := colly.NewCollector(
		colly.Async(true), // 异步
	)

	rp, err := proxy.RoundRobinProxySwitcher("http://127.0.0.1:7890")
	if err != nil {
		fmt.Println(err)
		// return
	}
	// 【设置代理IP】 ，这里使用的是轮询ip方式
	c1.SetProxyFunc(rp)

	c2 := c1.Clone()
	c3 := c1.Clone()

	c1.OnHTML("figure[itemProp] a[itemProp]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if href == "" {
			return
		}

		c2.Visit(e.Request.AbsoluteURL(href))
	})

	c2.OnHTML("div.MorZF > img[src]", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		if src == "" {
			return
		}

		c3.Visit(src)
	})

	c1.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c1.OnError(func(r *colly.Response, err error) {
		fmt.Println("Visiting", r.Request.URL, "failed:", err)
	})

	var count uint32
	c3.OnResponse(func(r *colly.Response) {
		fileName := fmt.Sprintf("images/img%d.jpg", atomic.AddUint32(&count, 1))
		err := r.Save(fileName)
		if err != nil {
			fmt.Printf("saving %s failed:%v\n", fileName, err)
		} else {
			fmt.Printf("saving %s success\n", fileName)
		}
	})

	c3.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	c1.Visit("https://unsplash.com/")
	// 异步需要等待结束
	c1.Wait()
	c2.Wait()
	c3.Wait()
}

type Item struct {
	Id     string
	Width  int
	Height int
	Links  Links
}

type Links struct {
	Download string
}

// 从动态加载的url来访问
func test4_v2() {
	c := colly.NewCollector(
		colly.Async(true),
	)
	c.SetRequestTimeout(time.Minute)
	rp, err := proxy.RoundRobinProxySwitcher("http://127.0.0.1:7890")
	if err != nil {
		fmt.Println(err)
		// return
	}
	// 设置代理IP
	c.SetProxyFunc(rp)
	d := c.Clone()
	c.OnResponse(func(r *colly.Response) {
		var items []*Item
		json.Unmarshal(r.Body, &items)
		for _, item := range items {
			d.Visit(item.Links.Download)
		}
	})

	var count uint32
	d.OnResponse(func(r *colly.Response) {
		fileName := fmt.Sprintf("images/img%d.jpg", atomic.AddUint32(&count, 1))
		err := r.Save(fileName)
		if err != nil {
			fmt.Printf("saving %s failed:%v\n", fileName, err)
		} else {
			fmt.Printf("saving %s success\n", fileName)
		}
	})

	d.OnRequest(func(r *colly.Request) {
		fmt.Println("d visiting", r.URL)
	})
	d.OnError(func(r *colly.Response, err error) {
		fmt.Println("d error:", err)
	})

	for page := 1; page <= 3; page++ {
		c.Visit(fmt.Sprintf("https://unsplash.com/napi/photos?page=%d&per_page=12", page))
	}
	c.Wait()
	d.Wait()
}

// type ProductInfo struct {
// 	Size                   string `selector:"a > div.index_1Ew5p"`
// 	VehicleServiceType     string
// 	Material               string
// 	Brand                  string
// 	FastenerMaterial       string
// 	ProductDimensions      string
// 	Manufacturer           string
// 	ItemWeight             string
// 	ItemModelNumber        string
// 	ManufacturerPartNumber string
// }

type ProductInfo struct {
	Details []Detail
}

type Detail struct {
	Key   string `selector:"th > .prodDetSectionEntry"`
	Value string `selector:"td > .prodDetAttrValue"`
}

func test5() {
	c := colly.NewCollector()
	details := ProductInfo{Details: make([]Detail, 0)}
	c.OnHTML("table.prodDetTable > tbody > tr", func(e *colly.HTMLElement) {
		hot := Detail{}

		hot.Key = e.ChildText("th.prodDetSectionEntry")
		hot.Value = strings.Replace(e.ChildText("td.prodDetAttrValue"),"‎","",-1) 
		details.Details = append(details.Details, hot)
	})
	// c.OnResponse(func(r *colly.Response) {

	// })
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("d visiting", r.URL)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("d error:", err)
	})
	c.Visit("https://www.amazon.com/dp/B09BFKBLQK")

	for _, v := range details.Details {
		fmt.Printf("%s : %s\n", v.Key, v.Value)
	}
}
