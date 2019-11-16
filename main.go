package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/higuching/slack_railway_bot/app"
)

func railways() (string, error) {
	text, err := app.Run()
	return text, err
}

func main() {
	lambda.Start(railways)
}
