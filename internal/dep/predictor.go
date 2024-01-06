package dep

import "github.com/pdstuber/isit-a-cat/pkg/prediction"

type ImagePredictor interface {
	PredictImage(imageBytes []byte) (*prediction.Result, error)
	Stop() error
}

type HasImagePredictor interface {
	ImagePredictor() ImagePredictor
}

func (d AppDependencies) ImagePredictor() ImagePredictor {
	return d.imagePredictor
}
