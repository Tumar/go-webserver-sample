package main

import (
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func newMinioClient() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewEnvMinio(),
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
