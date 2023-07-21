package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	//"github.com/tidwall/gjson"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	tr "github.com/aws/aws-sdk-go-v2/service/translate"
)

func trnaslate_en_to_kr(data string) string {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("S3_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("S3_ACCESSKEY"), os.Getenv("S3_PRIVATEDID"), "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := tr.NewFromConfig(cfg)


	translated, err := client.TranslateText(context.TODO(),
		&tr.TranslateTextInput{
			SourceLanguageCode: aws.String("en"),
			TargetLanguageCode: aws.String("ko"),
			Text: aws.String(data),
		},
	)
	if err != nil {
		log.Printf("Failed translate. err: %v\n", err)
	} else {
		log.Printf("translate done")
	}

	return *translated.TranslatedText
}
