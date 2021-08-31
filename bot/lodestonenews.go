package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LodestoneNewsResponse struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Time        time.Time `json:"time"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Emergency   bool      `json:"emergency"`
	Current     bool      `json:"current"`
	Category    string    `json:"category"`
}

type MaintenanceInfo struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Time      time.Time `json:"time"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Emergency bool      `json:"emergency"`
	Current   bool      `json:"current"`
}

type CurrentMainteenance struct {
	Companion LodestoneNewsResponse `json:"companion"`
	Game      LodestoneNewsResponse `json:"game"`
	Lodestone LodestoneNewsResponse `json:"lodestone"`
	Mog       LodestoneNewsResponse `json:"mog"`
	PSN       LodestoneNewsResponse `json:"psn"`
}

type Region string

const (
	baseURI = "http://%s.lodestonenews.com/news/%s"

	NorthAmerica  Region = "na"
	EuropeanUnion Region = "eu"
	France        Region = "fr"
	Germany       Region = "de"
	Japan         Region = "jp"
)

func Topics(locale Region) (news []LodestoneNewsResponse, err error) {
	return getData(locale, "topics")
}

func Notices(locale Region) (notices []LodestoneNewsResponse, err error) {
	return getData(locale, "notices")
}

func Maintenance(locale Region) (maintenance []LodestoneNewsResponse, err error) {
	return getData(locale, "maintenance")
}

func CurrentMaintenance(locale Region) (currentMaintenance CurrentMainteenance, err error) {
	targetURL := fmt.Sprintf(baseURI, locale, "maintenance/current")
	resp, err := http.Get(targetURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &currentMaintenance)
	return
}

func Updates(locale Region) (updates []LodestoneNewsResponse, err error) {
	return getData(locale, "updates")
}

func Status(locale Region) (status []LodestoneNewsResponse, err error) {
	return getData(locale, "status")
}

func Developers(locale Region) (status []LodestoneNewsResponse, err error) {
	return getData(locale, "developers")
}

func Feed(locale Region) (status []LodestoneNewsResponse, err error) {
	return getData(locale, "feed")
}

func getData(locale Region, urlPart string) (responseData []LodestoneNewsResponse, err error) {
	targetURL := fmt.Sprintf(baseURI, locale, urlPart)
	resp, err := http.Get(targetURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &responseData)
	return
}
