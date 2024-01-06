package prediction

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/pkg/errors"
	tf "github.com/wamuir/graft/tensorflow"
	"github.com/wamuir/graft/tensorflow/op"
	"golang.org/x/image/draw"
)

const (
	errorTextCouldNotCreateTensorFromImage     = "could not create tensor from input image"
	errorTextCouldNotDecodeJPEG                = "could not decode jpeg channels of image"
	errorTextCouldNotRunPreprocessImageSession = "could not run tensorflow session to preprocess input image"
	int32zero                                  = int32(0)
	vgg16ImagenetMeanRed                       = float32(123.68)
	vgg16ImagenetMeanGreen                     = float32(116.779)
	vgg16ImagenetMeanBlue                      = float32(103.939)
)

// TODO Needs to be synced with tensorflow code
const (
	targetImageSize    = 256
	targetImageMime    = "image/jpeg"
	targetImageQuality = 99
)

// VGG16 mean RGB values for the imagenet dataset
var imagenetMeans = []float32{vgg16ImagenetMeanRed, vgg16ImagenetMeanGreen, vgg16ImagenetMeanBlue}

// Preprocessing in specific to VGG16
func (s *Service) makeTensorFromImage(imageBytes []byte) (*tf.Tensor, error) {

	// DecodeJpeg uses a scalar String-valued tensor as inputOperation.
	tensor, err := tf.NewTensor(string(imageBytes))
	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotCreateTensorFromImage)
	}

	normalized, err := s.normalizationSession.Run(
		map[tf.Output]*tf.Tensor{*s.normalizationInput: tensor},
		[]tf.Output{*s.normalizationOutput},
		nil)
	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotRunPreprocessImageSession)
	}
	return normalized[0], nil
}

func decodeJPEGGraph(colorChannels int64) (*tf.Graph, *tf.Output, *tf.Output, error) {
	s := op.NewScope()

	mean := op.Const(s, imagenetMeans)
	input := op.Placeholder(s, tf.String)
	output := op.DecodeJpeg(s, input, op.DecodeJpegChannels(colorChannels))
	output = op.Cast(s, output, tf.Float)
	output = op.Sub(s, output, mean)
	output = op.ExpandDims(s, output, op.Const(s.SubScope("batch"), int32zero))

	graph, err := s.Finalize()
	return graph, &input, &output, err
}

// TODO user tensorflow for resizing
func resizeImage(imageBytes []byte) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}
	dst := image.NewRGBA(image.Rect(0, 0, targetImageSize, targetImageSize))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	jpeg.Encode(&buf, dst, &jpeg.Options{Quality: targetImageQuality})

	return buf.Bytes(), nil
}
