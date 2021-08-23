package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type MaintenanceInfo struct {
	ID    string    `json:"id"`
	URL   string    `json:"url"`
	Title string    `json:"title"`
	Time  time.Time `json:"time"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
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
