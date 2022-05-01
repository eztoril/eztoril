package price

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const (
	SpotAddr = "https://www.elbruk.se/timpriser-se3-stockholm"
)

func getHtmlPage(webPage string) (string, error) {

	//	fmt.Printf("getHtmlPage %s", webPage)
	resp, err := http.Get(webPage)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func parseMonthlySpot(input string) float64 {
	// Find correct data portion
	labelIx := strings.Index(input, "label: '2022'")
	if labelIx != -1 {
		labelStr := input[labelIx:]
		dataStartIx := strings.Index(labelStr, "data: [")
		dataStartStr := labelStr[dataStartIx:]
		dataEndIx := strings.Index(dataStartStr, " ]")
		dataStr := dataStartStr[:dataEndIx+1]

		re := regexp.MustCompile("[ ]?([0-9]*[.])?[0-9]+ ")
		found := re.FindAllString(dataStr, -1)
		lastSpotStr := strings.TrimSpace(found[len(found)-1])
		f, _ := strconv.ParseFloat(lastSpotStr, 32)
		return f
	}
	return -1
}

func GetPrice(webPage string) float64 {

	data, err := getHtmlPage(webPage)
	if err != nil {
		log.Fatal(err)
	}
	tkn := html.NewTokenizer(strings.NewReader(data))
	var isTd bool

	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return -1
		case tt == html.StartTagToken:
			t := tkn.Token()
			isTd = t.Data == "script"
		case tt == html.TextToken:
			t := tkn.Token()
			if isTd {
				spot := parseMonthlySpot(t.Data)
				if spot != -1 {
					return spot
				}
			}
			isTd = false
		}
	}
}
