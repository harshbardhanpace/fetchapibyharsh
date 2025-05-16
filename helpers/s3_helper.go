package helpers

import (
	"bytes"
	"fmt"
	"io"
	"space/constants"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/tealeg/xlsx/v3"
)

var awsSession *session.Session

var awsConfig *AWSS3CredConfig

// AWSS3CredConfig forms the aws config
type AWSS3CredConfig struct {
	Access_Key_Id     string `mapstructure:"aws_access_key_id"`
	Secret_Access_key string `mapstructure:"secret_access_key"`
	My_Region         string `mapstructure:"my_region"`
	BucketName        string `mapstructure:"bucket_name"`
}

// StartAwsSession - makes a new aws session
func StartAwsSession() (err error) {
	env := constants.Env

	awsConfig = &AWSS3CredConfig{
		My_Region:         constants.AWSRegion,
		Access_Key_Id:     constants.AWSAccessKeyID,
		Secret_Access_key: constants.AWSSecretAccessKey,
	}

	if env == constants.LocalEnv {
		awsConfig = &AWSS3CredConfig{
			My_Region:         constants.AWSRegion,
			Access_Key_Id:     constants.AWSAccessKeyID,
			Secret_Access_key: constants.AWSSecretAccessKey,
		}
	}

	awsSession, err = session.NewSession(
		&aws.Config{
			Region: aws.String(awsConfig.My_Region),
			Credentials: credentials.NewStaticCredentials(
				awsConfig.Access_Key_Id,
				awsConfig.Secret_Access_key,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		return fmt.Errorf("failed to create aws session, err: %s", err.Error())
	}

	return nil
}

// UploadToS3 - uploads file to aws s3
func UploadToS3(key string, body io.Reader) (error, string) {

	if awsSession == nil {
		if err := StartAwsSession(); err != nil {
			return fmt.Errorf("upload reconnection, err: %s", err.Error()), ""
		}
	}

	uploader := s3manager.NewUploader(awsSession)
	uploadRes, err := uploader.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(constants.AWSBucketName),
		ACL:    aws.String("public-read"),
		Key:    aws.String(key),
		//Metadata:    metaData,
		//ContentType: aws.String("binary/octet-stream"),
	})

	fmt.Printf("upload res=%v\n", uploadRes)

	if err != nil {
		return fmt.Errorf("s3 upload error, err: %s", err.Error()), ""
	}

	return nil, uploadRes.Location
}

func UploadFileToS3AndGetPresignedURL(folderName string, fileName string, file *xlsx.File, expiryHours int64) (string, error) {

	// Create a new AWS session using the default credential chain
	if awsSession == nil {
		if err := StartAwsSession(); err != nil {
			return "", fmt.Errorf("upload reconnection, err: %s", err.Error())
		}
	}

	// Create an S3 client
	s3Client := s3.New(awsSession)

	// Specify the object key
	objectKey := fmt.Sprintf("%s/%s", folderName, fileName)

	// Convert *xlsx.File to bytes
	var buffer bytes.Buffer
	if err := file.Write(&buffer); err != nil {
		return "", fmt.Errorf("failed to write XLSX file to buffer: %v", err)
	}

	// Perform the actual upload to S3
	uploader := s3manager.NewUploaderWithClient(s3Client)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(constants.AWSBucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(buffer.Bytes()),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload XLSX file to S3: %v", err)
	}

	// Generate a pre-signed URL for the object
	expiryDuration := time.Duration(expiryHours) * time.Hour
	urlReq, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(constants.AWSBucketName),
		Key:    aws.String(objectKey),
	})
	presignUrl, err := urlReq.Presign(expiryDuration)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %v", err)
	}

	return presignUrl, nil
}
