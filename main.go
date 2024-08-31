package main

import (
	"fmt"
	"github.com/yuri-ncs/novel-chapter-scraper/parser"
	"github.com/yuri-ncs/novel-chapter-scraper/req"
)

func main() {
	url := "https://novelbin.me/novel-book/shadow-slave"

	res, err := req.MakeRequest(url)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	parser.ParseHTML(res)

}
