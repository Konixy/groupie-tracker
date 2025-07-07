package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Artist struct {
	Id           int
	Name         string
	Image        string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

const (
	apiURL = "https://groupietrackers.herokuapp.com/api"
)

func main() {
	response, err := http.Get(apiURL + "/artists")
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Fatal(err)
	}

	for _, artist := range artists {
		fmt.Println(artist.Name)
	}
}
