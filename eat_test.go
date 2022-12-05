package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pfandzelter/go-eat/pkg/kaiserstueck"
	"github.com/pfandzelter/go-eat/pkg/personalkantine"
	"github.com/pfandzelter/go-eat/pkg/singh"
	"github.com/pfandzelter/go-eat/pkg/stw"
)

func testCanteen(name string, specDiet bool, c mensa) error {

	fl, err := c.GetFood(time.Now())

	if err != nil {
		log.Print(err)
		return err
	}

	if len(fl) == 0 {
		log.Printf("No food found for %s", name)
		return fmt.Errorf("No food found for %s", name)
	}

	log.Printf("%s: %+v\n", name, fl)
	return nil
}

func TestHauptmensa(t *testing.T) {
	err := testCanteen("Hauptmensa", true, stw.New(321))

	if err != nil {
		t.Error(err)
	}
}

func TestVeggie(t *testing.T) {
	err := testCanteen("Pasteria Veggie 2.0", true, stw.New(631))

	if err != nil {
		t.Error(err)
	}
}

func TestKaiserstueck(t *testing.T) {
	err := testCanteen("Kaiserstück", false, kaiserstueck.New())

	if err != nil {
		t.Error(err)
	}
}

func TestPersonalkantine(t *testing.T) {
	err := testCanteen("Personalkantine", true, personalkantine.New())

	if err != nil {
		t.Error(err)
	}
}

func TestSingh(t *testing.T) {
	err := testCanteen("Mathe Café", true, singh.New())

	if err != nil {
		t.Error(err)
	}
}

func TestAll(t *testing.T) {
	TestHauptmensa(t)
	TestVeggie(t)
	TestKaiserstueck(t)
	TestPersonalkantine(t)
	TestSingh(t)
}
