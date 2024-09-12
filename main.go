package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/yuri-ncs/novel-chapter-scraper/database"
	"github.com/yuri-ncs/novel-chapter-scraper/scraper"
	"github.com/yuri-ncs/novel-chapter-scraper/telegram"
)

func main() {

	godotenv.Load()

	database.Connect()

	go telegram.Start()

	c := cron.New()

	c.AddFunc("* * * * *", func() {
		scraper.Work()
	})

	fmt.Println("Cron job started!")
	c.Start()

	select {}

}
