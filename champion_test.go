package golol

import (
	"encoding/json"
	"errors"
	"testing"
)

var unMocked func(url string, w interface{}) error
var failRequest = func(url string, w interface{}) error {
	return errors.New("Error message here")
}

func mockRequest(v interface{}) func(url string, w interface{}) error {
	unMocked = request
	fakeResp, err := json.Marshal(v)
	if err != nil {
		return func(url string, w interface{}) error {
			return err
		}
	}

	return func(url string, w interface{}) error {
		err := json.Unmarshal(fakeResp, &w)
		return err
	}
}

func restoreRequest() {
	request = unMocked
}

func TestStatsAttackSpeed(t *testing.T) {
	s := &Stats{AttackSpeedDelay: 1.5}
	aSpeed := s.AttackSpeed()
	if aSpeed != 0.25 {
		t.Error("expected", aSpeed, "to be ", 0.25)
	}
}

func TestGetChampions(t *testing.T) {
	fakeChampions := Champions{
		Data: map[string]Champion{
			"whatever": Champion{Name: "whatever"},
		},
	}

	// ok case
	request = mockRequest(fakeChampions)
	defer restoreRequest()

	cs, err := GetChampions(BR)
	if err != nil {
		t.Error("unexpected err", err)
	}

	c, found := cs["whatever"]
	if !found {
		t.Error("expected whatever to be found in map[string]Champion")
	}

	if c.Name != "whatever" {
		t.Error("expected", c.Name, "to be", "whatever")
	}

	// bad case
	request = failRequest
	cs, err = GetChampions(BR)
	if err.Error() != "Error message here" {
		t.Error("expectec to return an error")
	}
}

func TestGetChapionById(t *testing.T) {
	fakeChampion := Champion{Name: "whatever"}

	// ok case
	request = mockRequest(fakeChampion)
	defer restoreRequest()

	c, err := GetChampionById(1, BR)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if c.Name != fakeChampion.Name {
		t.Error("expetected", c.Name, "to be", fakeChampion.Name)
	}

	// bad case
	request = failRequest
	c, err = GetChampionById(1, BR)
	if err.Error() != "Error message here" {
		t.Error("expectec to return an error")
	}
}

func TestGetChampionsByName(t *testing.T) {
	fakeChampions := Champions{
		Data: map[string]Champion{
			"whatever": Champion{Name: "whatever"},
		},
	}

	// ok case
	request = mockRequest(fakeChampions)
	defer restoreRequest()

	c, err := GetChampionByName("whatever", BR)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if c.Name != "whatever" {
		t.Error("expected", c.Name, "to be whatever")
	}

	// bad cases
	c, err = GetChampionByName("not found", BR)
	if err.Error() != "Couldn't find a champion named: not found" {
		t.Error("expectec to return an error")
	}

	request = failRequest
	c, err = GetChampionByName("whatever", BR)
	if err.Error() != "Error message here" {
		t.Error("expectec to return an error")
	}
}
