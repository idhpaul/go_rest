package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Presigner encapsulates the Amazon Simple Storage Service (Amazon S3) presign actions
// used in the examples.
// It contains PresignClient, a client that is used to presign requests to Amazon S3.
// Presigned requests contain temporary credentials and can be made from any HTTP client.
type Presigner struct {
	PresignClient *s3.PresignClient
}

// GetObject makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
func (presigner Presigner) GetObject(
	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

// PutObject makes a presigned request that can be used to put an object in a bucket.
// The presigned request is valid for the specified number of seconds.
func (presigner Presigner) PutObject(
	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

// DeleteObject makes a presigned request that can be used to delete an object from a bucket.
func (presigner Presigner) DeleteObject(bucketName string, objectKey string) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignDeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to delete object %v. Here's why: %v\n", objectKey, err)
	}
	return request, err
}

func s3client_init() {

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

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	//presignClient := s3.NewPresignClient(client)
	//presigner := Presigner{PresignClient: presignClient}

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}

	// log.Printf("Let's presign a request to Get Presigned the object.")
	// presignedGetRequest, err := presigner.GetObject("dolbyiohanedutech","input/0001.wav",60*30)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("Got a presigned %v request to URL:\n\t%v\n", presignedGetRequest.Method, presignedGetRequest.URL)

	// log.Printf("Let's presign a request to Put Presigned the object.")
	// presignedPutRequest, err := presigner.PutObject("dolbyiohanedutech","input/0001_out.wav",60*30)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("Got a presigned %v request to URL:\n\t%v\n", presignedPutRequest.Method, presignedPutRequest.URL)

}

func create_presignURL(num int) PresignURLs {

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

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)
	presigner := Presigner{PresignClient: presignClient}

	var urls []UrlData

	for i := 0; i < num; i++ {

		log.Printf("Let's presign a request to Get Presigned the object.")
		presignedGetRequest, err := presigner.GetObject(os.Getenv("S3_BUCKET_NAME"), "original/"+strconv.Itoa(i+1)+".wav", 60*30)
		if err != nil {
			panic(err)
		}
		log.Printf("Got a presigned %v presignedGetRequest to URL:\n\t%v\n", presignedGetRequest.Method, presignedGetRequest.URL)
		

		log.Printf("Let's presign a request to Put Presigned the object.")
		presignedPutRequest, err := presigner.PutObject(os.Getenv("S3_BUCKET_NAME"), "enhance/"+strconv.Itoa(i+1)+".wav", 60*30)
		if err != nil {
			panic(err)
		}
		log.Printf("Got a presigned %v presignedPutRequest to URL:\n\t%v\n", presignedPutRequest.Method, presignedPutRequest.URL)
		
		url := UrlData{presignedGetRequest.URL,presignedPutRequest.URL}
		urls= append(urls,url)
	}

	presignurls := PresignURLs{Count:num,Urls:urls}

	return presignurls
}
