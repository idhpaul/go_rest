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
	"sync"
	"time"

	"github.com/joho/godotenv"
	//"github.com/tidwall/gjson"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	tr "github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
)

func create_transcribe(idx int, isOriginal bool) STTStatus {

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

	if isOriginal {
		jobOriginal := tr.StartTranscriptionJobInput{
			TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(idx) + "_original"),
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
		_, err := client.StartTranscriptionJob(context.TODO(), &jobOriginal)
		if err != nil {
			log.Printf("Failed StartTranscriptionJob_original(%v). err: %v\n",idx, err)
			sttResult.Result = err.Error()
		} else {
			log.Printf("Transcription_original(%v) started", idx)
			sttResult.Result = "Transcription_original started"
		}
	} else {
		job := tr.StartTranscriptionJobInput{
			TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(idx)),
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
		_, err := client.StartTranscriptionJob(context.TODO(), &job)
		if err != nil {
			log.Printf("Failed StartTranscriptionJob(%v). err: %v\n",idx, err)
			sttResult.Result = err.Error()
		} else {
			log.Printf("Transcription(%v) started", idx)
			sttResult.Result = "Transcription started"
		}
	}

	return sttResult
}

func get_transcribe(idx int, isOriginal bool) STTStatus {

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

	if isOriginal {
		outputJobOriginal, err := client.GetTranscriptionJob(context.TODO(), &tr.GetTranscriptionJobInput{
			TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(idx) + "_original"),
		})
		if err != nil {
			log.Printf("Failed GetTranscriptionJob_original(%v). err: %v\n", idx, err)
			sttResult.Result = err.Error()
		} else {
			sttResult.Result = string(outputJobOriginal.TranscriptionJob.TranscriptionJobStatus)
		}
	} else {
		outputJob, err := client.GetTranscriptionJob(context.TODO(), &tr.GetTranscriptionJobInput{
			TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(idx)),
		})
		if err != nil {
			log.Printf("Failed GetTranscriptionJob(%v). err: %v\n", idx+1,err)
			sttResult.Result = err.Error()
		} else {
			sttResult.Result = string(outputJob.TranscriptionJob.TranscriptionJobStatus)
		}
	}

	return sttResult
}

func delete_transcribe(num int, isOriginal bool) string {

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

	if isOriginal {

		var isDone = false
		var waitDeleteTranscriptionJob sync.WaitGroup
		waitDeleteTranscriptionJob.Add(num)

		for i := 0; i < num; i++ {
			go func(i int, done bool) {
				defer waitDeleteTranscriptionJob.Done() //끝나면 .Done() 호출

				for {
					_, err := client.DeleteTranscriptionJob(context.TODO(), &tr.DeleteTranscriptionJobInput{
						TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(i) + "_original"),
					})
					if err != nil {
						log.Printf("Failed DeleteTranscriptionJob_original(idx : %v). err: %v\n", i, err)
					} else {
						log.Printf("delete transcription original(%v)", i)
						done = true
					}

					if done {
						break
					}

					time.Sleep(3 * time.Second)

				}
			}(i, isDone)
		}

		waitDeleteTranscriptionJob.Wait()

		var waitCleanUpS3 sync.WaitGroup

		waitCleanUpS3.Add(num)

		for i := 0; i < num; i++ {

			go func(i int) {
				defer waitCleanUpS3.Done() //끝나면 .Done() 호출
				cleanup_TranscribeData(i, "stt_original/"+strconv.Itoa(i+1))
			}(i)

		}

		waitCleanUpS3.Wait()

	} else {

		var isDone = false
		var waitDeleteTranscriptionJob sync.WaitGroup
		waitDeleteTranscriptionJob.Add(num)

		for i := 0; i < num; i++ {

			go func(i int, done bool) {

				defer waitDeleteTranscriptionJob.Done() //끝나면 .Done() 호출
				for {
					_, err := client.DeleteTranscriptionJob(context.TODO(), &tr.DeleteTranscriptionJobInput{
						TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(i)),
					})
					if err != nil {
						log.Printf("Failed DeleteTranscriptionJob(idx : %v). err: %v\n", i, err)
					} else {
						log.Printf("delete transcription(%v)", i)
						done = true
					}

					if done {
						break
					}

					time.Sleep(3 * time.Second)

				}
			}(i, isDone)
		}

		isDone = false
		waitDeleteTranscriptionJob.Wait()

		var waitCleanUpS3 sync.WaitGroup
		waitCleanUpS3.Add(num)

		for i := 0; i < num; i++ {

			go func(i int) {
				defer waitCleanUpS3.Done() //끝나면 .Done() 호출

				cleanup_TranscribeData(i, "stt/"+strconv.Itoa(i+1))

			}(i)
		}

		isDone = false
		waitCleanUpS3.Wait()
	}

	return "delete ok"

}

func test_delete_all(num int, isOriginal bool) string {

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

	if isOriginal {

		var isDone = false
		var waitDeleteTranscriptionJob sync.WaitGroup
		waitDeleteTranscriptionJob.Add(num)

		for i := 0; i < num; i++ {
			go func(i int, done bool) {
				defer waitDeleteTranscriptionJob.Done() //끝나면 .Done() 호출

				for {
					deloutputOriginal, err := client.DeleteTranscriptionJob(context.TODO(), &tr.DeleteTranscriptionJobInput{
						TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(i) + "_original"),
					})
					if err != nil {
						log.Printf("Failed DeleteTranscriptionJob_original(idx : %v). err: %v\n", i, err)
					} else {
						log.Println("delete transcription original,", deloutputOriginal)
						done = true
					}

					if done {
						break
					}

					time.Sleep(3 * time.Second)

				}
			}(i, isDone)
		}

		waitDeleteTranscriptionJob.Wait()


	} else {

		var isDone = false
		var waitDeleteTranscriptionJob sync.WaitGroup
		waitDeleteTranscriptionJob.Add(num)

		for i := 0; i < num; i++ {

			go func(i int, done bool) {

				defer waitDeleteTranscriptionJob.Done() //끝나면 .Done() 호출
				for {
					deloutput, err := client.DeleteTranscriptionJob(context.TODO(), &tr.DeleteTranscriptionJobInput{
						TranscriptionJobName: aws.String("dolbyEqualizeStt_" + strconv.Itoa(i)),
					})
					if err != nil {
						log.Printf("Failed DeleteTranscriptionJob(idx : %v). err: %v\n", i, err)
					} else {
						log.Println("delete transcription,", deloutput)
						done = true
					}

					if done {
						break
					}

					time.Sleep(3 * time.Second)

				}
			}(i, isDone)
		}

		waitDeleteTranscriptionJob.Wait()
	}

	return "delete ok"
}
