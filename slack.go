package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
)

func Sendslack(c *gin.Context, message string) {
	err := errors.New("")

	channel := c.Query("channel")
	// Throw error if channel is empty
	if channel == "" {
		channel = "#junk"
	}
	// Get token from environment variable
	token := getEnv("SLACK_TOKEN", "")
	// Throw error if token is empty
	if token == "" {
		log.Fatal("Environment variable `SLACK_TOKEN` is empty")
		c.String(http.StatusBadRequest, "Environment variable `SLACK_TOKEN` is empty")
		os.Exit(3)
	}
	// Send message to Slack
	api := slack.New(token)
	_, _, err = api.PostMessage(channel, slack.MsgOptionText(message, false))
	// Throw error if message could not be sent
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}
