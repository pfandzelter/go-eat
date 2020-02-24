package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pfandzelter/go-eat/pkg/dynamo"
	"github.com/pfandzelter/go-eat/pkg/food"
	"github.com/pfandzelter/go-eat/pkg/kaiserstueck"
	"github.com/pfandzelter/go-eat/pkg/personalkantine"
	"github.com/pfandzelter/go-eat/pkg/stw"
	"log"
	"os"
	"time"
)

type mensa interface {
	GetFood(date time.Time) ([]food.Food, error)
}

// HandleRequest handles one request to the lambda function.
func HandleRequest() {
	tablename := os.Getenv("DYNAMODB_TABLE")
	region := os.Getenv("DYNAMODB_REGION")

	db, err := dynamo.New(region, tablename)

	if err != nil {
		log.Fatal(err)
	}

	type Canteen struct {
		Name     string
		SpecDiet bool
	}

	canteens := make(map[Canteen]mensa)

	canteens[Canteen{
		Name:     "Hauptmensa",
		SpecDiet: true,
	}] = stw.New(321)
	canteens[Canteen{
		Name:     "Veggie 2.0",
		SpecDiet: true,
	}] = stw.New(631)
	canteens[Canteen{
		Name:     "Kaiserstück",
		SpecDiet: false,
	}] = kaiserstueck.New()
	canteens[Canteen{
		Name:     "Personalkantine",
		SpecDiet: true,
	}] = personalkantine.New()

	t := time.Now()

	foodlists := make(map[Canteen][]food.Food)

	for c, m := range canteens {
		fl, err := m.GetFood(t)
		if err != nil {
			log.Print(err)
			continue
		}
		foodlists[c] = fl
	}

	for c, f := range foodlists {
		err := db.PutFood(c.Name, c.SpecDiet, f, t)
		if err != nil {
			log.Print(err)
		}
	}

}

func main() {
	lambda.Start(HandleRequest)
}
