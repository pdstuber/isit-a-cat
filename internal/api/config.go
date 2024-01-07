package api

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/pdstuber/isit-a-cat/pkg/prediction"
	"github.com/pkg/errors"
)

type Config struct {
	ListenPort                   string
	Labels                       []prediction.Label
	Model                        []byte
	TargetImageDimensions        int
	TFInputOperationName         string
	TFOutputOperationName        string
	ObjectStorageEndpoint        string
	ObjectStorageAccessKeyID     string
	ObjectStorageSecretAccessKey string
	ObjectStorageUseTLS          bool
	ObjectStorageBucketName      string
	ObjectStorageObjectFolder    string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ConfigFromEnv() (*Config, error) {
	listenPort := getEnv("LISTEN_PORT", ":8080")

	var labels []prediction.Label

	modelPath := getEnv("MODEL_PATH", "/model")
	model, err := os.ReadFile(fmt.Sprintf("%s/model.pb", modelPath))
	if err != nil {
		return nil, errors.Wrap(err, "could not read model")
	}
	labelBytes, err := os.ReadFile(fmt.Sprintf("%s/labels.csv", modelPath))
	if err != nil {
		return nil, errors.Wrap(err, "could not read labels")
	}

	if err := gocsv.UnmarshalBytes(labelBytes, &labels); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal labels csv")
	}

	targetImageDimensions, err := strconv.Atoi(getEnv("TARGET_IMAGE_DIMENSIONS", "256"))
	if err != nil {
		return nil, errors.Wrap(err, "could not convert to integer, please use correct format")
	}

	inputOperationName := getEnv("TF_INPUT_OPERATION_NAME", "input_1")
	outputOperationName := getEnv("TF_OUTPUT_OPERATION_NAME", "dense_3/Softmax")

	objectStorageEndpoint := getEnv("OBJECT_STORAGE_ENDPOINT", "minio:9000")
	objectStorageAccessKeyID := getEnv("MINIO_ACCESS_KEY", "")
	if objectStorageAccessKeyID == "" {
		return nil, errors.New("object storage access key id is mandatory")
	}
	objectStorageSecretKey := getEnv("MINIO_SECRET_KEY", "")
	if objectStorageSecretKey == "" {
		return nil, errors.New("object storage secret key is mandatory")
	}
	objectStorageUseTLS, err := strconv.ParseBool(getEnv("OBJECT_STORAGE_USE_TLS", "false"))
	if err != nil {
		return nil, errors.Wrap(err, "could not parse as boolean")
	}

	storageBucketName := getEnv("STORAGE_BUCKET_NAME", "isit-a-cat")
	storageObjectFolder := getEnv("STORAGE_OBJECT_FOLDER", "uploaded-images/")

	return &Config{
		ListenPort:                   listenPort,
		Labels:                       labels,
		Model:                        model,
		TargetImageDimensions:        targetImageDimensions,
		TFInputOperationName:         inputOperationName,
		TFOutputOperationName:        outputOperationName,
		ObjectStorageEndpoint:        objectStorageEndpoint,
		ObjectStorageAccessKeyID:     objectStorageAccessKeyID,
		ObjectStorageSecretAccessKey: objectStorageSecretKey,
		ObjectStorageUseTLS:          objectStorageUseTLS,
		ObjectStorageBucketName:      storageBucketName,
		ObjectStorageObjectFolder:    storageObjectFolder,
	}, nil
}
