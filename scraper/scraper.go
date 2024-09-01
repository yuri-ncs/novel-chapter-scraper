package scraper

import (
	"fmt"
	"github.com/yuri-ncs/novel-chapter-scraper/parser"
	"github.com/yuri-ncs/novel-chapter-scraper/req"
)

func Work() {
	url := "https://novel-next.com/novelnext/castle-of-black-iron-nov1559752588"

	res, err := req.MakeRequest(url)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	parser.ParseHTML(res)
}
