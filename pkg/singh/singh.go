package singh

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pfandzelter/go-eat/pkg/food"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type singh struct{}

// New creates a new service to pull the menu from Personalkantine.
func New() *singh {
	return &singh{}
}

func (m *singh) GetFood(t time.Time) ([]food.Food, error) {
	// get today's date
	date := t.Weekday().String()

	switch date {
	case "Monday":
		date = "Montag"
	case "Tuesday":
		date = "Diensttag"
	case "Wednesday":
		date = "Mittwoch"
	case "Thursday":
		date = "Donnerstag"
	case "Friday":
		date = "Freitag"
	}

	// download the correct website
	resp, err := http.Get("http://singh-catering.de/cafe/")

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

	doc.Find(".menu-list.menu-list__dotted > .menu-list__title").Each(func(i int, t *goquery.Selection) {
		if strings.Contains(t.Text(), date) {
			t.Next().Next().Find(".menu-list__item").Each(func(i int, s *goquery.Selection) {

				name := s.Find(".menu-list__item-desc").Text()

				price := s.Find(".menu-list__item-price").Text()
				price = strings.Replace(price, "\n", "", -1)
				price = strings.Replace(price, " ", "", -1)
				price = strings.Replace(price, "€", "", -1)
				price = strings.Replace(price, "&euro;", "", -1)
				price = strings.Replace(price, ",", "", -1)
				price = strings.Replace(price, ".", "", -1)

				endprice, err := strconv.Atoi(price)

				if err != nil {
					return
				}

				vegetarian := strings.Contains(s.Find(".menu-list__item-highlight-wrapper > .menu-list__item-highlight-title").Text(), "VEGETARISCH")
				vegan := strings.Contains(s.Find(".menu-list__item-highlight-wrapper > .menu-list__item-highlight-title").Text(), "VEGAN")

				foodstuff[name] = food.Food{
					Name:       name,
					StudPrice:  endprice,
					ProfPrice:  endprice,
					Vegan:      vegan,
					Vegetarian: vegetarian,
					Fish:       false,
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
