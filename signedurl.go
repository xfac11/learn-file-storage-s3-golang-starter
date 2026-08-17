package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)
	objInput := s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	presignRequest, err := presignClient.PresignGetObject(context.TODO(), &objInput, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", fmt.Errorf("Could not generate presign request with the generated url : %s", err)
	}
	return presignRequest.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return database.Video{}, fmt.Errorf("The video url is nil")
	}
	bucketKey := strings.Split(*video.VideoURL, ",")
	if len(bucketKey) != 2 {
		return database.Video{}, fmt.Errorf("Could not retrieve bucket and key from video URL: %s", *video.VideoURL)
	}
	bucket := bucketKey[0]
	key := bucketKey[1]

	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, time.Minute*1.0)
	if err != nil {
		return database.Video{}, fmt.Errorf("Could not generate presigned url: %s", err)
	}
	video.VideoURL = &presignedURL
	return video, nil

}
