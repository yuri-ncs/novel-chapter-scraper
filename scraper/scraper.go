package scraper

import (
	"fmt"
	"github.com/yuri-ncs/novel-chapter-scraper/database"
	"github.com/yuri-ncs/novel-chapter-scraper/models"
	"github.com/yuri-ncs/novel-chapter-scraper/parser"
	"github.com/yuri-ncs/novel-chapter-scraper/req"
	"time"
)

func Work() {

	novels, err := database.GetNovelsToUpdate()

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for _, novel := range novels {
		res, err := req.MakeRequest(string(novel.URL))

		if err != nil {
			fmt.Println("Error making the request: ", err)
			continue
		}

		var chapter parser.Chapter

		err = parser.ParseHTML(res, &chapter)

		if err != nil {
			fmt.Println("Error parsing the html: ", err)
			continue
		}

		fmt.Println(chapter)

		var lastChapter *models.Chapter
		lastChapter, err = database.GetLastChapterByNovelID(novel.ID)

		if err != nil {
			fmt.Println("Error trying to get the last saved chapter: ", err)
			continue
		}

		if lastChapter != nil && lastChapter.Number >= chapter.Number {
			continue
		}

		newChapter := models.Chapter{
			NovelID: novel.ID,
			Number:  chapter.Number,
			Title:   chapter.Title,
			Href:    chapter.Href,
		}

		err = database.CreateChapter(&newChapter)

		if err != nil {
			fmt.Println("Error Trying to savechapter in db: ", err)
			continue
		}

		///tecnically here we should send a notification to the user
		///but for now we will just print the new chapter
		///and save in the database
		fmt.Println("New chapter: ", newChapter)

		time.Sleep(5 * time.Second)
	}

}
