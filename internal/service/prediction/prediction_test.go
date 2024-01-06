package prediction_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pdstuber/isit-a-cat/internal/service/prediction"
	"github.com/pdstuber/isit-a-cat/internal/service/prediction/mocks"
	prediction1 "github.com/pdstuber/isit-a-cat/pkg/prediction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testImageColorChannels  = int64(3)
	testStorageUploadFolder = "test"
	testImageID             = "123"
	testReplyToSubject      = "reply"
	mockErrorText           = "everything went to hell"
)

var (
	testPredictionInput  = prediction1.Input{ID: testImageID}
	invalidTestMessage   = []byte{9, 8, 7, 6, 5}
	mockPredictionResult = prediction1.Result{
		Class:       "banana",
		Probability: 0.99,
	}
	expectedErrorResult  = prediction1.ErrorResult{Message: "could not make prediction on image"}
	mockImage            = []byte{1, 2, 3, 4, 5, 6, 7}
	metricsPublisherMock *mocks.MetricsPublisher
	mockError            = errors.New(mockErrorText)
)

func init() {
	metricsPublisherMock = new(mocks.MetricsPublisher)

	metricsPublisherMock.On("RegisterIncomingRequest", mock.Anything).Return()
	metricsPublisherMock.On("RegisterPrediction", mock.Anything).Return()
	metricsPublisherMock.On("RegisterStorage", mock.Anything, mock.Anything).Return()

}
func Test_CalculateAndPublishPrediction_good_case(t *testing.T) {
	imagePredictorMock := new(mocks.ImagePredictor)
	storageReaderMock := new(mocks.StorageReader)

	imagePredictorMock.On("PredictImage", mock.Anything, mock.Anything).Return(&mockPredictionResult, nil)
	storageReaderMock.On("ReadFromBucketObject", mock.Anything, mock.Anything).Return(mockImage, nil)

	predictionService := prediction.NewService(metricsPublisherMock, imagePredictorMock, storageReaderMock, testStorageUploadFolder)

	result, err := predictionService.CalculatePredictionForStorageObject(testImageID)

	imagePredictorMock.AssertCalled(t, "PredictImage", mockImage, testImageColorChannels)
	storageReaderMock.AssertCalled(t, "ReadFromBucketObject", testStorageUploadFolder, testImageID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func Test_calculatePrediction_error_download_image_from_storage(t *testing.T) {
	imagePredictorMock := new(mocks.ImagePredictor)
	storageReaderMock := new(mocks.StorageReader)

	imagePredictorMock.On("PredictImage", mock.Anything, mock.Anything).Return(&mockPredictionResult, nil)
	storageReaderMock.On("ReadFromBucketObject", mock.Anything, mock.Anything).Return(nil, mockError)

	predictionService := prediction.NewService(metricsPublisherMock, imagePredictorMock, storageReaderMock, testStorageUploadFolder)

	prediction, err := predictionService.CalculatePredictionForStorageObject(testImageID)

	imagePredictorMock.AssertNumberOfCalls(t, "PredictImage", 0)
	storageReaderMock.AssertCalled(t, "ReadFromBucketObject", testStorageUploadFolder, testImageID)

	assert.Error(t, err)
	assert.Nil(t, prediction)
	assert.Condition(t, func() (success bool) {
		return strings.Contains(err.Error(), "could not fetch image from object storage")
	})
}

func Test_calculatePredictionCalculateAndPublishPrediction_error_prediction(t *testing.T) {
	imagePredictorMock := new(mocks.ImagePredictor)
	storageReaderMock := new(mocks.StorageReader)

	imagePredictorMock.On("PredictImage", mock.Anything, mock.Anything).Return(nil, mockError)
	storageReaderMock.On("ReadFromBucketObject", mock.Anything, mock.Anything).Return(mockImage, nil)

	predictionService := prediction.NewService(metricsPublisherMock, imagePredictorMock, storageReaderMock, testStorageUploadFolder)

	prediction, err := predictionService.CalculatePredictionForStorageObject(testImageID)

	imagePredictorMock.AssertCalled(t, "PredictImage", mockImage, testImageColorChannels)
	storageReaderMock.AssertCalled(t, "ReadFromBucketObject", testStorageUploadFolder, testImageID)

	assert.Error(t, err)
	assert.Nil(t, prediction)
	assert.Condition(t, func() (success bool) {
		return strings.Contains(err.Error(), "could not make prediction on image")
	})
}
