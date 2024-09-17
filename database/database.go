package database

import (
	"errors"
	"fmt"
	"github.com/yuri-ncs/novel-chapter-scraper/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	Url "net/url"
	"os"
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

	db.AutoMigrate(&models.Site{}, &models.Novel{}, &models.Chapter{}, &models.User{}, &models.UserNovel{})

	fmt.Println("Successfully connected to PostgreSQL!")

	db_global = db

}

func PopulateDatabase() {
	// Creating a sample site

	if db_global == nil {
		fmt.Println("Database connection is nil!")
		return
	}
	site := []models.Site{{
		Name:       "Novel Bin",
		DefaultURL: "https://novelbin.me"},
		{Name: "Novel Next",
			DefaultURL: "https://novel-next.com"},
		{Name: "All Novel Bin",
			DefaultURL: "https://allnovelbin.net"},
	}
	result := db_global.Create(&site)
	if result.Error != nil {
		fmt.Println("Error creating site:", result.Error)
		return
	}

	// Creating a sample novel for the site
	novel := models.Novel{
		SiteID:           site[0].ID,
		Name:             "Shadow Slave",
		URL:              "https://novelbin.me/novel-book/shadow-slave",
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

func VerifySupportedSite(url string) (bool, uint) {
	// Check if the site is supported
	u, err := Url.Parse(url)
	if err != nil {
		fmt.Println("Error parsing the URL:", err)
		return false, 0
	}

	var sites []models.Site
	db_global.Find(&sites)

	for _, site := range sites {

		s, _ := Url.Parse(site.DefaultURL)

		if s.Hostname() == u.Hostname() {
			return true, site.ID
		}
	}

	return false, 0
}

func CreateNovel(novel *models.Novel) error {
	result := db_global.Create(novel)
	if result.Error != nil {
		fmt.Println("Error creating novel:", result.Error)
		return result.Error
	}

	return nil
}

func GetAllNovels() []models.Novel {
	var novels []models.Novel
	db_global.Find(&novels)
	return novels
}

func GetActiveNovels() []models.Novel {
	var novels []models.Novel
	db_global.Where("deleted_at IS NULL").Find(&novels)
	return novels
}

func GetNovelsToUpdate() ([]models.Novel, error) {
	var novels []models.Novel

	result := db_global.Where("NOW() - updated_at >= INTERVAL '" + os.Getenv("SCRAPE_INTERVAL") + " minutes' AND deleted_at IS NULL").Find(&novels)
	if result.Error != nil {
		fmt.Println("Error getting novels to update:", result.Error)
		return nil, result.Error
	}

	return novels, nil
}

func UpdateNovel(novel *models.Novel) error {
	result := db_global.Save(novel)
	if result.Error != nil {
		fmt.Println("Error updating novel:", result.Error)
		return result.Error
	}

	return nil

}

func UpdateChapter(chapter *models.Chapter) error {
	result := db_global.Save(chapter)
	if result.Error != nil {
		fmt.Println("Error updating chapter:", result.Error)
		return result.Error
	}

	return nil
}

func GetLastChapterByNovelID(novelID uint) (*models.Chapter, error) {
	var chapter models.Chapter

	result := db_global.Where("novel_id = ?", novelID).Order("number DESC").First(&chapter)
	if result.Error != nil && !errors.Is(gorm.ErrRecordNotFound, result.Error) {
		fmt.Println("Error getting last chapter by novel ID:", result.Error)
		return nil, result.Error
	}

	return &chapter, nil
}

func CreateChapter(chapter *models.Chapter) error {
	result := db_global.Create(chapter)
	if result.Error != nil {
		fmt.Println("Error creating chapter:", result.Error)
		return result.Error
	}

	return nil
}

func GetSitesList() string {
	var sites []models.Site
	db_global.Find(&sites)

	var sitesList string
	for _, site := range sites {
		sitesList += fmt.Sprintf("ID [%d]. Site [%s]\n", site.ID, site.Name)
	}

	return sitesList
}

func FindUsersByNovelID(novelID uint) ([]models.User, error) {
	var users []models.User
	err := db_global.Preload("Novels").Where("novels.id = ?", novelID).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsersByNovelID(novelID uint) ([]models.User, error) {
	var users []models.User
	err := db_global.Joins("JOIN user_novels ON user_novels.user_id = users.id").
		Where("user_novels.novel_id = ?", novelID).
		Find(&users).Error

	if err != nil {
		fmt.Println("Error retrieving users:", err)
	} else {
		fmt.Printf("Found %d users associated with novel ID %d\n", len(users), novelID)
	}
	return users, nil
}

func CreateUser(user *models.User) error {
	result := db_global.Create(user)
	if result.Error != nil {
		fmt.Println("Error creating user:", result.Error)
		return result.Error
	}

	return nil
}

func GetUserByChatID(chatID int64) (*models.User, error) {
	var user models.User
	result := db_global.Where("chat_id = ?", chatID).First(&user)
	if result.Error != nil && !errors.Is(gorm.ErrRecordNotFound, result.Error) {
		fmt.Println("Error getting user by chat ID:", result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func GetNovelByName(name string) (*models.Novel, error) {
	var novel models.Novel
	result := db_global.Where("name = ?", name).First(&novel)
	if result.Error != nil && !errors.Is(gorm.ErrRecordNotFound, result.Error) {
		fmt.Println("Error getting novel by name:", result.Error)
		return nil, result.Error
	}

	return &novel, nil
}

func TrackNovel(user *models.User, novel *models.Novel) bool {

	db_global.Model(&user).Attrs(&models.User{Novels: []models.Novel{*novel}})

	return true
}
