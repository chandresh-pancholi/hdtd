package sink

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWS struct {
	s3 *s3.S3
}

func NewAwsS3Client() *AWS {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		log.Fatalf("s3 connection failed. Error: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	return &AWS{
		s3: svc,
	}
}

func (s3 *AWS) Upload(source, destination string) (string, error) {
	f, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fileInfo, _ := f.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	f.Read(buffer)

	uploader := s3manager.NewUploaderWithClient(s3.s3)

	// Perform an upload.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("DUMP_BUCKET")),
		Key:         aws.String(fmt.Sprintf("%s/%s/%s", os.Getenv("SERVICE_NAME"), os.Getenv("POD_NAME"), destination)),
		Body:        bytes.NewReader(buffer),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(http.DetectContentType(buffer)),
	})
	if err != nil {
		return "", err
	}

	log.Printf("successfully uploaded file. location: %s", result.Location)
	return result.Location, nil
}
