package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AntJanus/ngp-bot/config"
)

var (
	searchURL   = "games"
	searchQuery = "search \"%s\"; fields name,first_release_date,url; where version_parent = null;"
)

var client = &http.Client{}

/*
GameStruct represents the structure of a game returned from IGDB
*/
type GameStruct struct {
	Name             string `json:"name"`
	URL              string `json:"url"`
	FirstReleaseDate int64  `json:"first_release_date,omitempty"`
}

/*
GameList is a list of GameStruct
*/
type GameList []GameStruct

/*
Search fires off a request to IGDB and returns a game
*/
func Search(gameName string) (GameStruct, error) {
	requestURL := fmt.Sprintf("%s/%s", config.IGDBUrl, searchURL)
	searchBody := []byte(fmt.Sprintf(searchQuery, gameName))
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(searchBody))

	req.Header.Add("user-key", config.IGDBKey)
	req.Header.Add("Accepts", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		var empty GameStruct

		return empty, err
	}

	defer resp.Body.Close()

	var games GameList

	err = json.NewDecoder(resp.Body).Decode(&games)

	if err != nil {
		fmt.Println(err.Error())
		var empty GameStruct

		return empty, err
	}

	if len(games) == 0 {
		var empty GameStruct
		return empty, nil
	}

	game := games[0]

	return game, err
}
