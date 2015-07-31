package golol

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Image struct {
	Width  int    `json:"w"`
	Heigth int    `json:"w"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Group  string `json:"group"`
	Sprite string `json:"sprite"`
	Full   string `json:"full"`
}

var config = struct {
	APIKey        string
	StaticDataUrl string
}{
	os.Getenv("RIOT_API_KEY"),
	"https://global.api.pvp.net/api/lol/static-data/",
}

const (
	BR    = "br"
	EU_NE = "eune"
	EU_W  = "euw"
	KR    = "kr"
	LAN   = "lan"
	NA    = "na"
	OCE   = "oce"
	PBE   = "pbe"
	RU    = "ru"
	TR    = "tr"
)

func SetAPIKey(key string) {
	config.APIKey = key
}

var request = func(url string, v interface{}) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("Riot api returned " + strconv.Itoa(res.StatusCode))
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return err
	}

	return nil
}
