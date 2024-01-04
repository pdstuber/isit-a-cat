package prediction

import (
	"time"

	"github.com/pdstuber/isit-a-cat/internal/pkg/predict"
	"github.com/pkg/errors"
)

const (
	predictionErrorMessage                    = "could not make prediction on image"
	errorTextCouldNotUnmarshalIncomingMessage = "could not unmarshal incoming message"
	errorTextCouldNotFetchImageFromStorage    = "could not fetch image from object storage"
	errorTextCouldNotMakePredictionOnImage    = "could not make prediction on image"
)

// A MetricsPublisher publishes prometheus metrics
type MetricsPublisher interface {
	RegisterIncomingRequest(duration float64)
	RegisterPrediction(duration float64)
	RegisterStorage(duration float64, objectId string)
}

// A ImagePredictor predicts the class of an image
type ImagePredictor interface {
	PredictImage(imageBytes []byte) (*predict.Result, error)
	Stop() error
}

// A StorageReader reads data from the storage object with the given ID from the given folder
type StorageReader interface {
	ReadFromBucketObject(objectFolder string, objectId string) ([]byte, error)
}

// A Service for incoming prediction requests
type Service struct {
	metricsPublisher            MetricsPublisher
	imagePredictor              ImagePredictor
	storageReader               StorageReader
	colorChannels               int64
	storageUploadedImagesFolder string
}

// NewHandler creates a new instance of the Handler
func NewService(metricsPublisher MetricsPublisher, imagePredictor ImagePredictor, storageReader StorageReader, storageUploadedImagesFolder string) *Service {
	return &Service{
		metricsPublisher:            metricsPublisher,
		imagePredictor:              imagePredictor,
		storageReader:               storageReader,
		storageUploadedImagesFolder: storageUploadedImagesFolder,
	}
}

func (h *Service) CalculatePredictionForStorageObject(id string) (*predict.Result, error) {
	start := time.Now()

	defer (func() {
		elapsed := time.Since(start)
		go h.metricsPublisher.RegisterIncomingRequest(elapsed.Seconds())
	})()

	return h.calculatePrediction(id)

}

func (h *Service) calculatePrediction(id string) (*predict.Result, error) {
	start := time.Now()
	image, err := h.storageReader.ReadFromBucketObject(h.storageUploadedImagesFolder, id)
	elapsed := time.Since(start)
	go h.metricsPublisher.RegisterStorage(elapsed.Seconds(), id)

	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotFetchImageFromStorage)
	}

	start = time.Now()
	result, err := h.imagePredictor.PredictImage(image)
	elapsed = time.Since(start)
	go h.metricsPublisher.RegisterPrediction(elapsed.Seconds())

	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotMakePredictionOnImage)
	}

	return result, nil
}
