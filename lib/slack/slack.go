package slack

import (
	"fmt"

	"github.com/slack-go/slack"

	"swarm_deploy/lib/config"
)

var cnf config.Config = config.LoadConfig()

// Sends a slack plain text message
func SendSimpleMessage(msg string) {
	api := slack.New(cnf.Slack.Token)

	channelID, timestamp, err := api.PostMessage(
		cnf.Slack.Channel,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
