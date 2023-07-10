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
	"strconv"

	"github.com/joho/godotenv"
	//"github.com/tidwall/gjson"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	tr "github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
)

func create_transcribe(idx int) STTStatus {

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

	var sttResult STTStatus

	job := tr.StartTranscriptionJobInput{
		TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(idx+1)),
		Media: &types.Media{
			MediaFileUri: aws.String("s3://" + os.Getenv("S3_BUCKET_NAME") + "/equalize/" + strconv.Itoa(idx+1) + ".wav"),
		},
		MediaFormat:  "wav",
		LanguageCode: "ko-KR",

		// OutputKey 의 subtitles 파일명 영향 받음
		// Subtitles: &types.Subtitles{
		// 	Formats: []types.SubtitleFormat{types.SubtitleFormatVtt},
		// 	OutputStartIndex: aws.Int32(0),
		// },

		// https://docs.aws.amazon.com/ko_kr/AmazonS3/latest/userguide/access-bucket-intro.html
		// https://s3.us-east-2.amazonaws.com/{OutputBucketName}/{OutputKey}
		// https://docs.aws.amazon.com/transcribe/latest/APIReference/API_StartTranscriptionJob.html#transcribe-StartTranscriptionJob-request-OutputBucketName
		OutputBucketName: aws.String(os.Getenv("S3_BUCKET_NAME")),
		OutputKey:        aws.String("stt/" + strconv.Itoa(idx+1) + ".json"),
	}

	// start the transcription job
	resp, err := client.StartTranscriptionJob(context.TODO(), &job)
	if err != nil {
		log.Printf("Failed StartTranscriptionJob. err: %v\n", err)
		sttResult.Result1 = err.Error()
	} else {
		log.Println("Transcription started,", resp)
		sttResult.Result1 = "Transcription started"
	}

	jobOriginal := tr.StartTranscriptionJobInput{
		TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(idx+1) + "_original"),
		Media: &types.Media{
			MediaFileUri: aws.String("s3://" + os.Getenv("S3_BUCKET_NAME") + "/original/" + strconv.Itoa(idx+1) + ".wav"),
		},
		MediaFormat:  "wav",
		LanguageCode: "ko-KR",

		// OutputKey 의 subtitles 파일명 영향 받음
		// Subtitles: &types.Subtitles{
		// 	Formats: []types.SubtitleFormat{types.SubtitleFormatVtt},
		// 	OutputStartIndex: aws.Int32(0),
		// },

		// https://docs.aws.amazon.com/ko_kr/AmazonS3/latest/userguide/access-bucket-intro.html
		// https://s3.us-east-2.amazonaws.com/{OutputBucketName}/{OutputKey}
		// https://docs.aws.amazon.com/transcribe/latest/APIReference/API_StartTranscriptionJob.html#transcribe-StartTranscriptionJob-request-OutputBucketName
		OutputBucketName: aws.String(os.Getenv("S3_BUCKET_NAME")),
		OutputKey:        aws.String("stt_original/" + strconv.Itoa(idx+1) + ".json"),
	}

	// start the transcription job
	respOriginal, err := client.StartTranscriptionJob(context.TODO(), &jobOriginal)
	if err != nil {
		log.Printf("Failed StartTranscriptionJob_original. err: %v\n", err)
		sttResult.Result2 = err.Error()
	} else {
		log.Println("Transcription_original started,", respOriginal)
		sttResult.Result2 = "Transcription started"
	}

	return sttResult
}

func get_transcribe(idx int) STTStatus {

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
	var sttResult STTStatus

	outputJob, err := client.GetTranscriptionJob(context.TODO(), &tr.GetTranscriptionJobInput{
		TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(idx+1)),
	})
	if err != nil {
		log.Printf("Failed GetTranscriptionJob. err: %v\n", err)
		sttResult.Result1 = err.Error()
	} else {
		sttResult.Result1 = string(outputJob.TranscriptionJob.TranscriptionJobStatus)
	}

	outputJobOriginal, err := client.GetTranscriptionJob(context.TODO(), &tr.GetTranscriptionJobInput{
		TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(idx+1) + "_original"),
	})
	if err != nil {
		log.Printf("Failed GetTranscriptionJob_original. err: %v\n", err)
		sttResult.Result2 = err.Error()
	} else {
		sttResult.Result2 = string(outputJobOriginal.TranscriptionJob.TranscriptionJobStatus)
	}

	return sttResult
}

func cleanup_transcribe(idx int) STTStatus {

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

	var sttResult STTStatus
	deloutput, err := client.DeleteTranscriptionJob(context.TODO(), &tr.DeleteTranscriptionJobInput{
		TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(idx+1)),
	})
	if err != nil {
		log.Printf("Failed DeleteTranscriptionJob. err: %v\n", err)
		sttResult.Result1 = err.Error()
	} else {
		log.Println("delete transcription,", deloutput)

		cleanup_TranscribeData(idx, "stt/"+strconv.Itoa(idx+1)+".json")

		sttResult.Result1 = "cleanup done"
	}

	deloutputOriginal, err := client.DeleteTranscriptionJob(context.TODO(), &tr.DeleteTranscriptionJobInput{
		TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(idx+1)+"_original"),
	})
	if err != nil {
		log.Printf("Failed DeleteTranscriptionJob_original. err: %v\n", err)
		sttResult.Result2 = err.Error()
	} else {
		log.Println("delete transcription_original,", deloutputOriginal)

		cleanup_TranscribeData(idx, "stt_original/"+strconv.Itoa(idx+1)+".json")

		sttResult.Result2 = "original cleanup done"
	}

	return sttResult
}

func delete_trnascribe(num int) string {

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

	for i := 0; i < num; i++ {
		deloutput, err := client.DeleteTranscriptionJob(context.TODO(), &tr.DeleteTranscriptionJobInput{
			TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(i+1)),
		})
		if err != nil {
			log.Printf("Failed DeleteTranscriptionJob(idx : %v). err: %v\n", i, err)
		}
		log.Println("delete transcription,", deloutput)

		deloutputOriginal, err := client.DeleteTranscriptionJob(context.TODO(), &tr.DeleteTranscriptionJobInput{
			TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(i+1)+"_original"),
		})
		if err != nil {
			log.Printf("Failed DeleteTranscriptionJob_original(idx : %v). err: %v\n", i, err)
		}
		log.Println("delete transcription original,", deloutputOriginal)
	}

	for i := 0; i < num; i++ {
		cleanup_TranscribeData(i, "stt/"+strconv.Itoa(i+1))
		cleanup_TranscribeData(i, "stt_original/"+strconv.Itoa(i+1))
	}

	return "delete ok"

}
