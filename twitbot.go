package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Formats3link(k string) string {
	return fmt.Sprintf("https://%v.s3-%v.amazonaws.com/%v", BUCKET, REGION, k)
}

// Structfield to represent our uploaded file meta-data
type UploadInfo struct {
	Key     string
	Uri     string
	Created time.Time
	Size    int64
	Type    string
}

// Structfields representing twitter response body
type Media struct {
	Media_Key string
	Type      string
	Url       string
}
type Media_key struct {
	Key string
}
type Includes struct {
	Media []Media
}
type Attachments struct {
	Media_keys []Media_key
}
type Data struct {
	Attachments Attachments
	Id          string
	Text        string
	Includes    Includes
}

//Struct type for representing events
type TweetImage struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

// function to fetch, manipulate and upload media obj to our s3 bucket.
func Fetch_tweet(ctx context.Context, event TweetImage) (*Data, error) {
	/*
		most url to be inputted by users are most likely to be in the twitter v1 api
		format, so we parse the tweet id into the twitter v2 api url format.
	*/
	idls := strings.Split(event.Url, "/")
	baseid := idls[len(idls)-1]
	baseidls := strings.Split(baseid, "?")
	uri := BASE_URL + baseidls[0] + PARAMS
	fmt.Println(uri)
	client := http.Client{}

	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{AUTH},
	}
	data := new(Data)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// decode the response body to a  bytes
	respByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// decode bytes response to a readable string
	respStr := string(respByte)
	fmt.Println(respStr)

	// decode byte response to our data struct
	e := json.Unmarshal(respByte, &data)
	if e != nil {
		log.Fatal(e)
	}

	// find the media url in the data struct and send it to the getImage function for fetching img
	if len(data.Includes.Media) >= 1 {

		for i := 0; i < len(data.Includes.Media); i++ {
			err := FetchImage(data.Includes.Media[i].Url)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("File url -------")
			fmt.Println((data.Includes.Media[i].Url))
		}
	}

	return data, nil
}

func FetchImage(url string) (err error) {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// upload file to s3
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if err = UploadImage(sesss, resp.ContentLength, body); err != nil {
		log.Fatal(err)
	}

	return nil
}

// function to upload file to s3
func UploadImage(session *session.Session, contentLength int64, body []byte) error {

	imageName := RandStringRunes(10) + ".png"
	uploadInfo := &UploadInfo{
		Key:     imageName,
		Created: time.Now(),
		Uri:     UPLOADED_FILE_BASE_URI + imageName,
		Size:    contentLength,
	}
	fmt.Println(uploadInfo)
	_, err := s3sess.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(BUCKET),
		Key:                  aws.String(imageName),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(body),
		ContentLength:        aws.Int64(contentLength),
		ContentType:          aws.String(http.DetectContentType(body)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
