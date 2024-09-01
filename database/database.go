package database

import (
	"fmt"
	"github.com/yuri-ncs/novel-chapter-scraper/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
)

var db_global *gorm.DB

func Connect() {

	// Load environment variables (assuming they are set in the environment)
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	// Connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open the connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Site{}, &models.Novel{}, &models.Chapter{})

	fmt.Println("Successfully connected to PostgreSQL!")

	db_global = db

}

func PopulateDatabase() {
	// Creating a sample site

	if db_global == nil {
		fmt.Println("Database connection is nil!")
		return
	}
	site := models.Site{
		Name:       "Novel Bin",
		DefaultURL: "https://novelbin.me",
	}
	result := db_global.Create(&site)
	if result.Error != nil {
		fmt.Println("Error creating site:", result.Error)
		return
	}

	// Creating a sample novel for the site
	novel := models.Novel{
		SiteID:           site.ID,
		Name:             "Shadow Slave",
		NumberOfChapters: 1800,
	}
	result = db_global.Create(&novel)
	if result.Error != nil {
		fmt.Println("Error creating novel:", result.Error)
		return
	}

	// Creating sample chapters for the novel
	chapters := []models.Chapter{
		{NovelID: novel.ID, Number: 1, Title: "Chapter 1: The Beginning"},
		{NovelID: novel.ID, Number: 2, Title: "Chapter 2: The Journey"},
		{NovelID: novel.ID, Number: 3, Title: "Chapter 3: The Challenge"},
	}

	for _, chapter := range chapters {
		result = db_global.Create(&chapter)
		if result.Error != nil {
			fmt.Println("Error creating chapter:", result.Error)
			return
		}
	}

	fmt.Println("Database populated successfully!")
}

func VerifySupportedSite(url string) bool {
	// Check if the site is supported
	var sites []models.Site
	db_global.Find(&sites)

	for _, site := range sites {
		if strings.HasPrefix(url, site.DefaultURL) {
			return true
		}
	}

	return false
}