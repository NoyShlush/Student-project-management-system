package aws

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Info struct {
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}

var (
	e Info
)

var AllowTypes = []string{
	"application/pdf",
	"application/vnd.ms-powerpoint",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/msword",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/zip",
	"application/x-7z-compressed",
	"application/x-rar-compressed",
	"image/jpeg",
}
var AllowExtensions = []string{
	".pdf",
	".ppt",
	".pptx",
	".doc",
	".docx",
	".zip",
	".7z",
	".rar",
	".jpeg",
	".jpg",
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return false
		}
	}
	return true
}

func Configure(c Info) {
	e = c
}

func FileToChat(file multipart.File, header *multipart.FileHeader) (string, error) {

	// Upload file to chat
	path, err := addFileToS3("chat", file, header)
	if err != nil {
		return "", err
	}

	return path, nil
}

func FileToProject(file multipart.File, header *multipart.FileHeader) (string, error) {

	// Upload file to chat
	path, err := addFileToS3("project", file, header)
	if err != nil {
		return "", err
	}

	return path, nil
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func addFileToS3(folder string, f multipart.File, h *multipart.FileHeader) (string, error) {

	// Create a single AWS session
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(e.Region),
		Credentials: credentials.NewStaticCredentials(
			e.AccessKey,
			e.SecretKey,
			""),
	})

	if err != nil {
		return "", err
	}

	// Get file size and read the file content into a buffer
	size := h.Size
	buffer := make([]byte, size)
	f.Read(buffer)

	fileName, _ := uuid.NewRandom()
	path := folder + "/" + fileName.String() + filepath.Ext(h.Filename)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(e.Bucket),
		Key:                  aws.String(path),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		return "", err
	}

	return "http://" + e.Bucket + "/" + path, nil
}

func SendSMS(phoneNumber string, message string) error {

	// Create a single AWS session
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			e.AccessKey,
			e.SecretKey,
			""),
	})

	// Set the SMS title
	attrs := map[string]*sns.MessageAttributeValue{}
	attrs["AWS.SNS.SMS.SenderID"] = &sns.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String("SPMS"),
	}

	// Sends a text message (SMS message) directly to a phone number.
	_, err = sns.New(s).Publish(&sns.PublishInput{
		PhoneNumber:       aws.String(phoneNumber),
		Message:           aws.String(message),
		MessageAttributes: attrs,
	})

	if err != nil {
		return err
	}

	return nil
}
