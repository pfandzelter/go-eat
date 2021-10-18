package main

import (
	"log"
	"testing"
	"time"

	"github.com/pfandzelter/go-eat/pkg/kaiserstueck"
	"github.com/pfandzelter/go-eat/pkg/personalkantine"
	"github.com/pfandzelter/go-eat/pkg/singh"
	"github.com/pfandzelter/go-eat/pkg/stw"
)

func testCanteen(name string, specDiet bool, c mensa) {

	fl, err := c.GetFood(time.Now())

	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("%s: %+v\n", name, fl)
}

func TestHauptmensa(_ *testing.T) {
	testCanteen("Hauptmensa", true, stw.New(321))
}

func TestVeggie(_ *testing.T) {
	testCanteen("Veggie 2.0", true, stw.New(631))
}

func TestKaiserstueck(_ *testing.T) {
	testCanteen("Kaiserstück", false, kaiserstueck.New())
}

func TestPersonalkantine(_ *testing.T) {
	testCanteen("Personalkantine", true, personalkantine.New())
}

func TestSingh(_ *testing.T) {
	testCanteen("Mathe Café", true, singh.New())
}