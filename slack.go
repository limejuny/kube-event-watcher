package main

import (
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"os"
)

type SlackConf struct {
	Token   string
	Channel string
}

var slackConf = loadSlackEnvVal()

var slackColors = map[string]string{
	"Normal":  "good",
	"Warning": "warning",
	"Danger":  "danger",
}

func loadSlackEnvVal() SlackConf {
	var s SlackConf
	s.Token = os.Getenv("SLACK_TOKEN")
	s.Channel = os.Getenv("SLACK_CHANNEL")
	return s
}

//実際postする以外にprivateチャンネルの存在確認する方法は…
func validateSlack() error {
	//s := loadSlackEnvVal()
	if slackConf.Token == "" || slackConf.Channel == "" {
		return errors.New("slack error: token or channel is empty")
	}
	fmt.Printf("slack channel: %v\n", slackConf.Channel)
	api := slack.New(slackConf.Token)
	title := "kube-event-watcher (beta)"
	text := "application start"
	params := prepareParams(title, text, "good")
	_, _, err := api.PostMessage(slackConf.Channel, "", params)
	if err != nil {
		return err
	}
	return nil
}

func prepareParams(title string, text string, color string) slack.PostMessageParameters {
	params := slack.PostMessageParameters{
		AsUser: true,
	}
	attachment := slack.Attachment{
		Color: color,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: title,
				Value: text,
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	return params
}

func postEventToSlack(message string, action string, status string) error {
	api := slack.New(slackConf.Token)
	title := "kubernetes event : " + action
	color, ok := slackColors[status]
	if !ok {
		color = "danger"
	}
	params := prepareParams(title, message, color)
	_, _, err := api.PostMessage(slackConf.Channel, "", params)
	if err != nil {
		return err
	}
	return nil
}
