package main

import (
	"github.com/joho/godotenv"
	"github.com/yuri-ncs/novel-chapter-scraper/database"
	"github.com/yuri-ncs/novel-chapter-scraper/scraper"
)

func main() {

	godotenv.Load()

	database.Connect()

	scraper.Work()

}
