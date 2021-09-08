package communication

import (
	"os"

	"github.com/slack-go/slack"
)

type Slack struct {
	sc *slack.Client
}

func NewSlack() *Slack {
	sc := slack.New(os.Getenv("SLACK_AUTH_TOKEN"), slack.OptionDebug(true))

	return &Slack{
		sc: sc,
	}
}

func (sc *Slack) Publish(message string) error {
	channelID := os.Getenv("SLACK_CHANNEL_ID")

	attachment := slack.Attachment{
		Text: message,
	}
	_, _, err := sc.sc.PostMessage(
		channelID,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		return err
	}

	return nil

}
