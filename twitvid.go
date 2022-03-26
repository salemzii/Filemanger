package main

/*


// function to get each media url found within data.Includes.Media[index].Url
func GetImage(url string) (bytes []byte, err error) {

	resp, err := http.Get(url)

	if err != nil {
		return bytes, err
	}
	defer resp.Body.Close()

	// create empty file to hold the incoming media body content.
	file, err := os.Create(RandStringRunes(10) + ".png")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// copy downloaded media content into our earlier created file ^
	_, err = io.Copy(file, resp.Body)
	fmt.Println(reflect.TypeOf(resp.Body))

	// upload file to s3
	err = UploadFile(sesss, file.Name())
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.ReadAll(resp.Body)
}

// function to upload file to s3
func UploadFile(session *session.Session, uploadFileDir string) error {

	curr_dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
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
*/
