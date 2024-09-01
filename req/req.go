package req

import (
	"fmt"
	"github.com/yuri-ncs/novel-chapter-scraper/database"
	"io/ioutil"
	"net/http"
)

func MakeRequest(url string) (string, error) {

	if !database.VerifySupportedSite(url) {
		return "", fmt.Errorf("site not supported")
	}

	// Get the response from the URL
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}

	return string(data), nil
}
