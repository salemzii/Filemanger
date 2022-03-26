package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var BASE_URL, PARAMS, FILE_LINK, AUTH, BUCKET, REGION, UPLOADED_FILE_BASE_URI string

var s3sess *s3.S3
var sesss *session.Session

func init() {

	rand.Seed(time.Now().UnixNano())

	REGION = os.Getenv("REGION")
	BUCKET = os.Getenv("BUCKET")
	fmt.Println(REGION, BUCKET)

	// start aws sessions
	sesss := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	}))

	// initiate s3 with created session object
	s3sess = s3.New(sesss)

	//base twitter api url v2
	BASE_URL = "https://api.twitter.com/2/tweets/"

	//fields we're trying to query
	PARAMS = "?expansions=attachments.media_keys&media.fields=url"

	// link to media obj we're trying to fetch
	FILE_LINK = "https://twitter.com/CynthiaNomso02/status/1507602381122457600?s=20&t=CkTvEpyJiyOhriXcKDrT7g"
	// url of our base s3 bucket
	UPLOADED_FILE_BASE_URI = "https://salems-test-bucket.s3.eu-west-3.amazonaws.com/"
	AUTH = os.Getenv("TWITTER_BEARER")
	fmt.Println(AUTH)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	lambda.Start(Fetch_tweet)
}
