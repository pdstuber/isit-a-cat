package prediction

import (
	"github.com/pdstuber/isit-a-cat/internal/dep"
	"github.com/pdstuber/isit-a-cat/pkg/prediction"
	"github.com/pkg/errors"
)

const (
	predictionErrorMessage                    = "could not make prediction on image"
	errorTextCouldNotUnmarshalIncomingMessage = "could not unmarshal incoming message"
	errorTextCouldNotFetchImageFromStorage    = "could not fetch image from object storage"
	errorTextCouldNotMakePredictionOnImage    = "could not make prediction on image"
)

type serviceDependencies interface {
	dep.HasStorageReader
	dep.HasImagePredictor
}

func CalculatePrediction(deps serviceDependencies, id string) (*prediction.Result, error) {
	image, err := deps.StorageReader().ReadFromBucketObject(id)

	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotFetchImageFromStorage)
	}

	result, err := deps.ImagePredictor().PredictImage(image)
	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotMakePredictionOnImage)
	}

	return result, nil
}
