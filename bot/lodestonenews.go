package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

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
	Companion MaintenanceInfo `json:"companion"`
	Game      MaintenanceInfo `json:"game"`
	Lodestone MaintenanceInfo `json:"lodestone"`
	Mog       MaintenanceInfo `json:"mog"`
	PSN       MaintenanceInfo `json:"psn"`
}

func getCurrentMaintenance() (currentMaintenance CurrentMainteenance, err error) {
	resp, err := http.Get("http://na.lodestonenews.com/news/maintenance/current")
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

func getMaintenance() (allMaintenance []MaintenanceInfo, err error) {
	resp, err := http.Get("http://na.lodestonenews.com/news/maintenance")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &allMaintenance)
	return
}
