package kaiserstueck

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pfandzelter/go-eat/pkg/food"
	"net/http"
	"strings"
	"time"
)

const price = 580

type kaiserstk struct{}

// New creates a new service to pull the menu from Kaiserstück.
func New() *kaiserstk {
	return &kaiserstk{}
}

func (m *kaiserstk) GetFood(t time.Time) ([]food.Food, error) {
	// download the correct website
	resp, err := http.Get("https://kaiserstueck.de/")

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

	doc.Find(".entry-content").Each(func(i int, s *goquery.Selection) {
		name := s.Find("p").Text()
		name = strings.Replace(name, "\n", " ", -1)

		veg := strings.Contains(name, "veg.")

		foodstuff[name] = food.Food{
			Name:       name,
			StudPrice:  price,
			ProfPrice:  price,
			Vegan:      false,
			Vegetarian: veg,
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
