package sheet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AntJanus/ngp-bot/config"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

var apiURL = "https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?key=%s"

type sheetStruct struct {
	Range          string     `json:"range"`
	MajorDimension string     `json:"majorDimension"`
	Values         [][]string `json:"values"`
}

/*
RowStruct describes structure of spreadsheet data
*/
type RowStruct struct {
	Name        string
	Date        string
	Platform    string
	NGP         int
	EpisodeNum  string
	EpisodeLink string
	Guest       string
	PatreonPool string
	ExactMatch  bool
	Salty       string
}

/*
ReadSheet fetches spreadsheet, reads it, and matches a game to input gameName
*/
func ReadSheet(gameName string) (RowStruct, error) {
	query := fmt.Sprintf(apiURL, config.WorkBookID, config.SheetName, config.GoogleAPIKey)

	resp, err := http.Get(query)

	if err != nil {
		fmt.Println("Error is here")
		fmt.Println(err.Error())
		return RowStruct{}, err
	}

	defer resp.Body.Close()

	var sheet sheetStruct

	err = json.NewDecoder(resp.Body).Decode(&sheet)

	if err != nil {
		fmt.Println("Error is here")
		fmt.Println(err.Error())
		return RowStruct{}, err
	}

	gameTitles := []string{}
	sheetMap := make(map[string]RowStruct)

	for key, val := range sheet.Values {
		// skip first 2 rows
		if key == 0 || key == 1 {
			continue
		}

		/*
			Structure:
				Title
				Date
				Original System
				Episode #,
				NG+
		*/

		gameTitles = append(gameTitles, val[0])

		rowMap := RowStruct{
			Name:       val[0],
			Date:       checkIfExists(1, val),
			Platform:   checkIfExists(2, val),
			EpisodeNum: checkIfExists(3, val),
			NGP:        0,
			ExactMatch: false,
			Salty:      "",
		}

		if checkIfExists(4, val) != "" {
			if val[4] == "-" || val[4] == "--" {
				rowMap.NGP = -1
			} else if val[4] == "NG++" || val[4] == "NG+" {
				rowMap.NGP = 1
			}
		}

		sheetMap[strings.ToLower(val[0])] = rowMap
	}

	var gameListing RowStruct

	if val, ok := sheetMap[strings.ToLower(gameName)]; ok {
		gameListing = val
		gameListing.ExactMatch = true
	} else {
		findResult := fuzzy.RankFindFold(gameName, gameTitles)

		fmt.Println("Results")
		fmt.Println(findResult)
		if len(findResult) == 0 {
			return RowStruct{}, nil
		}

		firstResult := findResult[0]
		gameMatch := firstResult.Target
		matchAccuracy := firstResult.Distance

		fmt.Println(matchAccuracy)

		gameListing = sheetMap[strings.ToLower(gameMatch)]
	}

	if val, ok := config.Salty[gameListing.Name]; ok {
		gameListing.Salty = "and " + val + " is salty about it"
	}

	return gameListing, nil
}

func checkIfExists(idx int, list []string) string {
	if len(list) > idx {
		return list[idx]
	}

	return ""
}
