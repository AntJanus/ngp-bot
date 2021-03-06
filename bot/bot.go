package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/AntJanus/ngp-bot/config"
	"github.com/AntJanus/ngp-bot/igdb"
	"github.com/AntJanus/ngp-bot/sheet"
	"github.com/bwmarrin/discordgo"
)

var BotID string
var goBot *discordgo.Session

var (
	yesEmoji   = "white_check_mark"
	noEmoji    = "no_entry"
	plusEmoji  = "white_check_mark"
	minusEmoji = "no_entry"
)

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")

	return
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		command := strings.TrimPrefix(m.Content, config.BotPrefix)

		fmt.Println("Runs")
		fmt.Println(command)

		if m.Author.ID == BotID {
			return
		}

		if strings.HasPrefix(command, "salty") {
			query := strings.TrimPrefix(command, "salty ")

			gameMatch, _ := sheet.ReadSheet(query)
			// ngpStatus := ""
			message := ""

			if gameMatch.ExactMatch == false {
				message += fmt.Sprintf("Cannot find an exact match for: %s", query)
				message += "Did you mean? \n"
				message += fmt.Sprintf("Game: %s\nDate: %s\n", gameMatch.Name, gameMatch.Date)

				_, _ = s.ChannelMessageSend(m.ChannelID, message)

				return
			}

			if gameMatch.NGP != 1 && gameMatch.NGP != -1 {
				message += fmt.Sprintf("Game has not been NGP rated yet: %s", query)

				_, _ = s.ChannelMessageSend(m.ChannelID, message)

				return
			}

			message += "Saltiness is still being configured..."
			// Store.saltySave(game.Name, user)

			// message += fmt.Sprintf("%s your saltiness have been registered", user)

			_, _ = s.ChannelMessageSend(m.ChannelID, message)

			return
		}

		if strings.HasPrefix(command, "check") {
			query := strings.TrimPrefix(command, "check ")
			game, err := igdb.Search(query)

			if err != nil {
				fmt.Println("Error is here")
				fmt.Println(err.Error())
				return
			}

			if game.Name == "" {
				message := fmt.Sprintf("Cannot find Game: %s", query)

				_, _ = s.ChannelMessageSend(m.ChannelID, message)

				return
			}

			releaseDate := game.FirstReleaseDate
			unixDate := time.Unix(releaseDate, 0)
			humanDate := unixDate.Format("01/02/2006")

			nowDate := time.Now()
			dateDifference := nowDate.Sub(unixDate)
			yearDiff := dateDifference.Hours() / 24 / 365

			statusEmoji := ""

			if yearDiff > 15 {
				statusEmoji = yesEmoji
			} else {
				statusEmoji = noEmoji
			}

			// cover
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", game.URL))

			// game information
			message := ""

			if releaseDate == 0 {
				message = fmt.Sprintf("Game: %s \nDate: [unknown] \n", game.Name)

				_, _ = s.ChannelMessageSend(m.ChannelID, message)
			} else {
				message = fmt.Sprintf("Game: %s \nDate: %s \nEligible: :%s:", game.Name, humanDate, statusEmoji)

				_, _ = s.ChannelMessageSend(m.ChannelID, message)

				if yearDiff < 15 {
					fmt.Println("Game ineligible for NGP")
					return
				}
			}

			query = game.Name
			gameMatch, _ := sheet.ReadSheet(query)
			ngpStatus := ""
			message = ""

			if len(gameMatch.Name) == 0 {
				message += "Could not find game in master list"
			} else {
				if gameMatch.ExactMatch == false {
					message += "The closest match I could find: \n"
					message += fmt.Sprintf("Game: %s\nDate: %s\n", gameMatch.Name, gameMatch.Date)
				} else {
					message += "Game is in the master list\n"
				}
			}

			if gameMatch.NGP == 1 {
				ngpStatus = plusEmoji
				message += fmt.Sprintf("NGP: :%s: %s\n", ngpStatus, gameMatch.Salty)
				message += fmt.Sprintf("Ep#: %s", gameMatch.EpisodeNum)
			} else if gameMatch.NGP == -1 {
				ngpStatus = minusEmoji
				message += fmt.Sprintf("NGP: :%s: %s", ngpStatus, gameMatch.Salty)
			}

			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		}
	}

}
