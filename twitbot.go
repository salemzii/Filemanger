package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var BASE_URL, PARAMS, FILE_LINK, AUTH, BUCKET, REGION, UPLOADED_FILE_BASE_URI string

var s3sess *s3.S3
var sesss *session.Session

func init() {

	REGION = os.Getenv("REGION")
	BUCKET = os.Getenv("BUCKET")

	sesss := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	}))

	s3sess = s3.New(sesss)

	BASE_URL = "https://api.twitter.com/2/tweets/"

	PARAMS = "?expansions=attachments.media_keys&media.fields=url"

	FILE_LINK = "https://twitter.com/0x/status/1500560662933774342?s=20&t=yQ9z8sR7cvuLzmkBDhsB0g"
	UPLOADED_FILE_BASE_URI = "https://salems-test-bucket.s3.eu-west-3.amazonaws.com/"
	AUTH = os.Getenv("TWITTER_BEARER")
}

func Formats3link(k string) string {
	return fmt.Sprintf("https://%v.s3-%v.amazonaws.com/%v", BUCKET, REGION, k)
}

type UploadInfo struct {
	Key     string
	Uri     string
	Created time.Time
	Size    int64
	Type    string
}
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

func main() {
	fmt.Println(fetch_tweet(FILE_LINK))
}

func fetch_tweet(url string) (*Data, []byte, error) {
	idls := strings.Split(url, "/")
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
		"Authorization": []string{"Bearer AAAAAAAAAAAAAAAAAAAAAEMIaAEAAAAAMFa82YQWqtKV3qSetdowN%2FV9avA%3DniIQuHgi8GJWp0jrbvMmIWYkmT7vl2sP3lScHRzhvMKW4vfPzT"},
	}
	data := new(Data)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	respByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	respStr := string(respByte)
	fmt.Println(respStr)

	e := json.Unmarshal(respByte, &data)
	if e != nil {
		log.Fatal(e)
	}
	img, err := GetImage(data.Includes.Media[0].Url)
	if err != nil {
		log.Fatal(err)
	}
	return data, img, nil
}

func GetImage(url string) (bytes []byte, err error) {

	resp, err := http.Get(url)

	if err != nil {
		return bytes, err
	}
	defer resp.Body.Close()

	file, err := os.Create("twitimage.png")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	err = uploadFile(sesss, "twitimage.png")
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.ReadAll(resp.Body)
}

func uploadFile(session *session.Session, uploadFileDir string) error {

	curr_dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(curr_dir)
	upFle := curr_dir + "/" + uploadFileDir

	upFile, err := os.Open(upFle)
	if err != nil {
		return err
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()

	uploadInfo := &UploadInfo{
		Key:     uploadFileDir,
		Created: time.Now(),
		Uri:     UPLOADED_FILE_BASE_URI + uploadFileDir,
		Size:    fileSize,
	}
	fmt.Println(uploadInfo)
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	_, err = s3sess.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(BUCKET),
		Key:                  aws.String(uploadFileDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
