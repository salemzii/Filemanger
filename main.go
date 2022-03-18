package main

/*


import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var BUCKET, REGION string
var s3session *s3.S3

var imageLinks []string
var statuscode int
var ContentType string
var image = "image/"

type Output struct {
	StatusCode      int               `json:"statuscode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:isBase64encoded`
}

func init() {
	REGION = os.Getenv("REGION")
	BUCKET = os.Getenv("BUCKET")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	}))

	s3session = s3.New(sess)

	resp, err := s3session.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(BUCKET)})
	if err != nil {
		log.Println(err)
		return
	}

	for _, obj := range resp.Contents {
		imageLinks = append(imageLinks, formats3link(*obj.Key))
		fmt.Println(formats3link(*obj.Key))

		if strings.HasSuffix(*obj.Key, "jpeg") {
			ContentType = image + "jpeg"
		} else if strings.HasSuffix(*obj.Key, "png") {
			ContentType = image + "png"
		} else if strings.HasSuffix(*obj.Key, "jpg") {
			ContentType = image + "jpg"
		} else {
			log.Println("Unknown Image format")
		}
	}
}

func main() {
	lambda.Start(handler)
}

func handler() (o Output, err error) {
	link := getRandom()
	img, err := getImage(link)
	if err != nil {
		log.Fatal(err)
		statuscode = 400
	}

	switch statuscode {
	case 400:
		o = Output{
			StatusCode: statuscode,
		}

	default:
		o = Output{
			StatusCode:      200,
			Headers:         map[string]string{"Content-Type": ContentType},
			Body:            base64.StdEncoding.EncodeToString(img),
			IsBase64Encoded: true,
		}
	}
	return o, err
}

func getRandom() string {
	img := rand.Intn(len(imageLinks))
	return imageLinks[img]
}
func getImage(url string) (bytes []byte, err error) {

	resp, err := http.Get(url)

	if err != nil {
		return bytes, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func formats3link(k string) string {
	return fmt.Sprintf("https://%v.s3-%v.amazonaws.com/%v", BUCKET, REGION, k)
}

*/
