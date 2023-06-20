package main

// 공식 문서
// https://docs.aws.amazon.com/transcribe/?icmpid=docs_homepage_ml
// https://docs.aws.amazon.com/transcribe/latest/APIReference/API_Operations_Amazon_Transcribe_Service.html

// github
// https://docs.aws.amazon.com/transcribe/latest/APIReference/API_Operations_Amazon_Transcribe_Service.html

// 예제
// https://github.com/aws-samples/cross-aws-sdk-workshop/blob/31eed8668ae573862fe9422184501170ad25d16d/lambda/go/start-transcription/main.go

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	tr "github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
)

func create_trnascribe() string {

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

	job := tr.StartTranscriptionJobInput{
		TranscriptionJobName: aws.String("dolbywav2text"),
		Media: &types.Media{
			MediaFileUri: aws.String("s3://dolbyio/equalize/10.wav"),
		},
		MediaFormat: "wav",
		LanguageCode: "ko-KR",

		// https://docs.aws.amazon.com/ko_kr/AmazonS3/latest/userguide/access-bucket-intro.html
		// https://s3.us-east-2.amazonaws.com/{OutputBucketName}/{OutputKey}
		// https://docs.aws.amazon.com/transcribe/latest/APIReference/API_StartTranscriptionJob.html#transcribe-StartTranscriptionJob-request-OutputBucketName
		OutputBucketName: aws.String("dolbyio"),
		OutputKey: aws.String("stt/10.json"),
	}

	// start the transcription job
	resp, err := client.StartTranscriptionJob(context.TODO(), &job)
	if err != nil {
		log.Printf("Failed StartTranscriptionJob. err: %v\n", err)
	}
	log.Println("transcription started,", resp)

	time.Sleep(30 * time.Second)

	outputJob, err := client.GetTranscriptionJob(context.TODO(),&tr.GetTranscriptionJobInput{
		TranscriptionJobName: aws.String("dolbywav2text"),
	})
	if err != nil {
		log.Printf("Failed StartTranscriptionJob. err: %v\n", err)
	}
	log.Println("transcription started,", outputJob.TranscriptionJob.Transcript.TranscriptFileUri)

	return get_TranscribeOutputJson("dolbyio","stt/10.json")
}