package stw

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pfandzelter/go-eat/pkg/food"
)

// Mensen:
// 321 - TU Hardenbergstr
// 631 - TU Veggie 2.0
type mensa struct {
	id int
}

var blacklist = [...]string{
	"kuchen",
	"creme",
	"torte",
	"Brownie",
}

// New creates a new service to pull the menu a STW Mensa based on an id.
func New(id int) *mensa {
	return &mensa{
		id: id,
	}
}

func checkBlacklist(name string) bool {
	for _, item := range blacklist {
		if strings.Contains(strings.ToUpper(name), strings.ToUpper(item)) {
			return true
		}
	}

	return false
}

func (m *mensa) GetFood(t time.Time) ([]food.Food, error) {
	// get today's date
	date := t.Format("2006-01-02")

	// download the correct website
	// should be something like:
	// $ curl 'https://www.stw.berlin/xhr/speiseplan-wochentag.html' -v --data 'resources_id=321&date=2020-02-21' --compressed
	data := []byte(fmt.Sprintf("resources_id=%d&date=%s", m.id, date))

	resp, err := http.Post("https://www.stw.berlin/xhr/speiseplan-wochentag.html",
		"application/x-www-form-urlencoded", bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// parse the results
	foodstuff := make(map[string]food.Food)

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	doc.Find(".splGroupWrapper").Each(func(i int, t *goquery.Selection) {
		if t.Find("div > .splGroup").Text() == "Aktionen" || t.Find("div > .splGroup").Text() == "Essen" {
			t.Find(".splMeal").Each(func(i int, s *goquery.Selection) {
				name := s.Find("div > .bold").Text()

				if checkBlacklist(name) {
					return
				}

				price := s.Find(".col-xs-6.col-md-3.text-right").Text()
				price = strings.Replace(price, "\n", "", -1)
				price = strings.Replace(price, " ", "", -1)
				price = strings.Replace(price, "€", "", -1)
				price = strings.Replace(price, "&euro;", "", -1)
				price = strings.Replace(price, ",", "", -1)

				prices := strings.Split(price, "/")

				studPrice, err := strconv.Atoi(prices[0])

				if err != nil {
					return
				}

				profPrice, err := strconv.Atoi(prices[1])

				if err != nil {
					return
				}

				vegetarian := false
				vegan := false
				fish := false

				s.Find("div > .splIcon").Each(func(i int, x *goquery.Selection) {
					src, ok := x.Attr("src")

					if !ok {
						return
					}

					if src == "/vendor/infomax/mensen/icons/15.png" {
						vegan = true
						return
					}

					if src == "/vendor/infomax/mensen/icons/1.png" {
						vegetarian = true
						return
					}

					if src == "/vendor/infomax/mensen/icons/38.png" {
						fish = true
						return
					}
				})

				foodstuff[name] = food.Food{
					Name:       name,
					StudPrice:  studPrice,
					ProfPrice:  profPrice,
					Vegan:      vegan,
					Vegetarian: vegetarian,
					Fish:       fish,
				}

			})
		}
	})

	// return stuff
	foodlist := make([]food.Food, len(foodstuff))
	i := 0

	for _, f := range foodstuff {
		foodlist[i] = f
		i++
	}

	return foodlist, nil

}
