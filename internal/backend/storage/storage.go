package storage

import (
	"bytes"
	"io"
	"log"

	"github.com/minio/minio-go/v6"
	"github.com/pkg/errors"
)

const (
	errorTextCouldNotCreateClient = "could not create storage client"
	errorTextBucketWrite          = "could not write to bucket"
	errorTextBucketRead           = "could not read from bucket"
)

// Service handles writes and reads from object storage buckets
type Service struct {
	client              StorageObjectReaderWriter
	storageObjectFolder string
	storageBucketName   string
}

// A StorageObjectReaderWriter reads and writes storage bucket objects
type StorageObjectReaderWriter interface {
	PutObject(bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (n int64, err error)
	GetObject(bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
}

// New creates an instance of the storage service
func New(storageBucketName string, storageObjectFolder string, endpoint string, accessKeyID string, secretAccessKey string, secure bool) (*Service, error) {
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)

	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotCreateClient)
	}

	return &Service{
		storageBucketName:   storageBucketName,
		storageObjectFolder: storageObjectFolder,
		client:              minioClient,
	}, nil
}

// WriteToBucketObject writes data to the bucket object with the given ID
func (service *Service) WriteToBucketObject(objectID string, data []byte) error {
	storageObjectPath := service.storageObjectFolder + objectID
	dataLen := int64(len(data))

	if _, err := service.client.PutObject(service.storageBucketName, storageObjectPath, bytes.NewReader(data), dataLen, minio.PutObjectOptions{}); err != nil {

		return errors.Wrap(err, errorTextBucketWrite)
	}

	log.Printf("Successfully wrote %d bytes to %s/%s\n", len(data), service.storageBucketName, storageObjectPath)

	return nil
}

// ReadFromBucketObject reads data from the bucket object with the given ID
func (service *Service) ReadFromBucketObject(objectId string) ([]byte, error) {
	storageObjectPath := service.storageObjectFolder + objectId

	log.Printf("Trying to read from %v/%v\n", service.storageBucketName, storageObjectPath)

	object, err := service.client.GetObject(service.storageBucketName, storageObjectPath, minio.GetObjectOptions{})

	if err != nil {
		return nil, errors.Wrap(err, "could not get bucket object")
	}

	defer func() {
		err := object.Close()
		if err != nil {
			log.Printf("Error in closing object: %v\n", err)
		}
	}()

	data, err := io.ReadAll(object)

	if err != nil {
		return nil, errors.Wrap(err, errorTextBucketRead)
	}

	log.Printf("Successfully read %v bytes from %v/%v\n", len(data), service.storageBucketName, storageObjectPath)

	return data, err
}
