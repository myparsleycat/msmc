package util

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateClient() (*s3.Client, error) {
	accessKey := os.Getenv("R2_ADMIN_ACCESS_KEY")
	secretKey := os.Getenv("R2_ADMIN_SECRET_KEY")
	accountId := os.Getenv("R2_ACCOUNT_ID")

	endpointUrl := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("apac"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("AWS 설정을 불러오는 중 오류 발생: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = &endpointUrl
	})

	return client, nil
}

func UploadFile(client *s3.Client, bucketName, fileName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("파일을 여는 중 오류 발생: %v", err)
	}
	defer file.Close()

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("S3 업로드 중 오류 발생: %v", err)
	}

	return nil
}

func ListFiles(client *s3.Client, bucketName string, prefix string) ([]string, error) {
	var files []string

	input := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix: aws.String(prefix),
	}

	result, err := client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("파일 목록 조회 중 오류 발생: %v", err)
	}

	for _, object := range result.Contents {
		files = append(files, *object.Key)
	}

	// 날짜순 정렬 (최신 순)
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	return files, nil
}

func DeleteFile(client *s3.Client, bucketName string, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    aws.String(key),
	}

	_, err := client.DeleteObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("파일 삭제 중 오류 발생: %v", err)
	}

	return nil
}

func DeleteFiles(client *s3.Client, bucketName string, keys []string) error {
	for _, key := range keys {
		err := DeleteFile(client, bucketName, key)
		if err != nil {
			return err
		}
	}
	return nil
}
