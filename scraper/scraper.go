package scraper

import (
	"fmt"
	"github.com/yuri-ncs/novel-chapter-scraper/database"
	"github.com/yuri-ncs/novel-chapter-scraper/models"
	"github.com/yuri-ncs/novel-chapter-scraper/parser"
	"github.com/yuri-ncs/novel-chapter-scraper/req"
	"github.com/yuri-ncs/novel-chapter-scraper/telegram"
	"time"
)

func Work() {

	novels, err := database.GetNovelsToUpdate()

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Novels to update: ", novels)
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

		if err == nil {
			fmt.Println("Last chapter: ", lastChapter.Number)
			fmt.Println("Current chapter: ", chapter.Number)

			if lastChapter.Number >= chapter.Number {
				fmt.Println("No new chapters")
				novel.UpdatedAt = time.Now()
				err = database.UpdateNovel(&novel)
				if err != nil {
					fmt.Println("Error trying to update the novel: ", err)
				}
				continue
			}

		} else {
			fmt.Println("Error trying to get the last chapter: ", err)
			fmt.Println("Creating the first chapter")
		}

		newChapter := models.Chapter{
			NovelID: novel.ID,
			Number:  chapter.Number,
			Title:   chapter.Title,
			Href:    chapter.Href,
		}

		novel.NumberOfChapters = chapter.Number
		novel.UpdatedAt = time.Now()

		telegram.SendNotification(newChapter, novel.Name)

		fmt.Println("New chapter: ", newChapter)

		err = database.UpdateNovel(&novel)

		err = database.CreateChapter(&newChapter)

		if err != nil {
			fmt.Println("Error Trying to savechapter in db: ", err)
			continue
		}

	}

}
