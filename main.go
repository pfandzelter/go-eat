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

	canteens := make(map [string]mensa)

	canteens["Hauptmensa"] = stw.New(321)
	canteens["Veggie 2.0"] = stw.New(321)
	canteens["Kaiserstück"] = kaiserstueck.New()
	canteens["Personalkantine"] = personalkantine.New()

	t := time.Now()

	foodlists := make(map [string][]food.Food)

	for c, m := range canteens {
		fl, err := m.GetFood(t)
		if err != nil {
			log.Print(err)
			continue
		}
		foodlists[c] = fl
	}

	for c, f := range foodlists {
		err := db.PutFood(c, f)
		if err != nil {
			log.Print(err)
		}
	}

}

func main() {
	lambda.Start(HandleRequest)
}
