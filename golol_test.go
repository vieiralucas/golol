package golol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestConfig(t *testing.T) {
	if config.APIKey != os.Getenv("RIOT_API_KEY") {
		t.Error("expected " + config.APIKey + " to be " + os.Getenv("RIOT_API_KEY"))
	}

	if config.StaticDataUrl != "https://global.api.pvp.net/api/lol/static-data/" {
		t.Error("expected " + config.StaticDataUrl + " to be https://global.api.pvp.net/api/lol/static-data/")
	}
}

func TestConstants(t *testing.T) {
	if BR != "br" {
		t.Error("expected " + BR + " to be br")
	}
	if EU_NE != "eune" {
		t.Error("expected " + EU_NE + " to be eune")
	}
	if EU_W != "euw" {
		t.Error("expected " + EU_W + " to be euw")
	}
	if KR != "kr" {
		t.Error("expected " + KR + " to be kr")
	}
	if LAN != "lan" {
		t.Error("expected " + LAN + " to be lan")
	}
	if NA != "na" {
		t.Error("expected " + NA + " to be na")
	}
	if OCE != "oce" {
		t.Error("expected " + OCE + " to be oce")
	}
	if PBE != "pbe" {
		t.Error("expected " + PBE + " to be pbe")
	}
	if RU != "ru" {
		t.Error("expected " + RU + " to be ru")
	}
	if TR != "tr" {
		t.Error("expected " + TR + " to be tr")
	}
}

func TestSetAPIKey(t *testing.T) {
	SetAPIKey("whatever")
	if config.APIKey != "whatever" {
		t.Error("expected " + config.APIKey + " to be whatever")
	}
}

type fakeJsonStruct struct {
	status int
	some   string
}

func (f *fakeJsonStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(f.status)
	json.NewEncoder(w).Encode(f)
}

func TestRequest(t *testing.T) {
	body := fakeJsonStruct{status: http.StatusOK, some: "json"}
	res := &fakeJsonStruct{}
	s := httptest.NewServer(&body)
	defer s.Close()

	err := request(s.URL, &res)
	if err != nil {
		t.Error(fmt.Sprintf("unexpected error: %v", err))
	}

	body.status = http.StatusBadRequest
	err = request(s.URL, &res)

	if err.Error() != "Riot api returned "+strconv.Itoa(body.status) {
		t.Error("expected to return following err:", "Riot api returned "+strconv.Itoa(body.status))
	}
}
