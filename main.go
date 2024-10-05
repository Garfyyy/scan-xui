package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Garfyyy/scan-xui/fofa"
	"github.com/Garfyyy/scan-xui/task"
	"github.com/Garfyyy/scan-xui/utils"
	"github.com/joho/godotenv"
)

func main() {
	f := flag.String("f", "", "file path")
	flag.Parse()

	if *f != "" {
		content, err := os.ReadFile(*f)
		if err != nil {
			panic(err)
		}
		links := strings.Split(string(content), "\n")
		RunScanXui(links, "custom")
		return
	}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	fofaKey := os.Getenv("FOFA_KEY")
	countrys := strings.Split(os.Getenv("COUNTRY"), ",")
	fids := strings.Split(os.Getenv("FID"), ",")

	// Create a new client
	fc := fofa.NewClient(fofaKey)

	for _, fid := range fids {
		for _, country := range countrys {
			links := GetLinksByFofa(fid, country, fc)
			RunScanXui(links, country)
		}
	}
}

func GetLinksByFofa(fid, country string, fc *fofa.Client) []string {
	fofaQuery := &fofa.SearchParams{
		Query:  fmt.Sprintf(`title=="登录" && fid="%s" && country="%s"`, fid, country),
		Size:   10000,
		Page:   1,
		Fields: "ip,port,country,region,link",
	}

	fmt.Printf("Query: %s\n", fofaQuery.Query)
	fofaRes, err := fc.Search(fofaQuery)
	if err != nil {
		panic(err)
	}

	links := make([]string, 0, len(fofaRes.Results))
	linksMap := make(map[string]bool)
	for _, res := range fofaRes.Results {
		link := res[len(res)-1]
		if !linksMap[link] {
			linksMap[link] = true
			links = append(links, res[len(res)-1])
		}
	}
	return links
}

func RunScanXui(links []string, country string) {
	res := task.ScanLinks(links)

	if len(res) > 0 {
		saveFile := fmt.Sprintf("result/%s.txt", country)

		ipMap := make(map[string]bool)
		var uniqueRes []string

		for _, link := range res {
			u, err := url.Parse(link)
			if err != nil {
				continue
			}
			ip := u.Hostname()
			if !ipMap[ip] {
				ipMap[ip] = true
				uniqueRes = append(uniqueRes, link)
			}
		}

		if err := utils.Write2File(saveFile, uniqueRes); err != nil {
			panic(err)
		} else {
			fmt.Printf("Write to file %s success, len: \033[32m%d\033[0m\n", saveFile, len(uniqueRes))
		}
	}
}
