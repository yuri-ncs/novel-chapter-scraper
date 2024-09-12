package req

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func MakeRequest(url string) (string, error) {
	///it makes more sense if we make this verification when this is added in the db
	//if !database.VerifySupportedSite(url) {
	//	return "", fmt.Errorf("site not supported")
	//}

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
