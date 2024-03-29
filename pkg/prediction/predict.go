package prediction

import (
	"log"

	"github.com/pkg/errors"
	tf "github.com/wamuir/graft/tensorflow"
)

const (
	errorTextTensorflowEmptyResponse          = "tensorflow session produced empty result"
	errorTextCouldNotExecuteTensorflowSession = "could not execute tensorflow session"
	errorTextCouldNotProcessInputImage        = "could not process input image"
)

// Service predicts images using an imported tensorflow model
type Service struct {
	inputOperation        *tf.Operation
	outputOperation       *tf.Operation
	session               *tf.Session
	normalizationSession  *tf.Session
	normalizationInput    *tf.Output
	normalizationOutput   *tf.Output
	labels                []Label
	targetImageDimensions int
}

// NewService creates a new service instance from the given model and labels
func NewService(model []byte, labels []Label, colorChannels int64, inputOperationName, outputOperationName string, targetImageDimensions int) *Service {
	graph, err := createTensorFlowGraphFromModel(model)

	if err != nil {
		log.Fatalf("could not import tensorflow graph: %v\n", err)
	}

	inputOperation := graph.Operation(inputOperationName)
	outputOperation := graph.Operation(outputOperationName)

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Fatalf("could not create tensorflow session: %v/n", err)
	}

	// Creates a tensorflow graph to decode the jpeg image
	normalizationGraph, normalizationInput, normalizationOutput, err := decodeJPEGGraph(colorChannels)
	if err != nil {
		log.Fatalf("could not create tensorflow graph: %v/n", err)
	}

	// Execute that graph to decode this one image
	normalizationSession, err := tf.NewSession(normalizationGraph, nil)
	if err != nil {
		log.Fatalf("could not create tensorflow session: %v/n", err)
	}

	return &Service{inputOperation, outputOperation, session, normalizationSession, normalizationInput, normalizationOutput, labels, targetImageDimensions}
}

// PredictImage with the imported tensorflow model and labels
func (s *Service) PredictImage(imageBytes []byte) (*Result, error) {
	resizedImageBytes, err := s.resizeImage(imageBytes)
	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotProcessInputImage)
	}

	inputTensor, err := s.makeTensorFromImage(resizedImageBytes)
	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotProcessInputImage)
	}

	results, err := s.session.Run(
		map[tf.Output]*tf.Tensor{
			s.inputOperation.Output(0): inputTensor,
		},
		[]tf.Output{
			s.outputOperation.Output(0),
		},
		nil)

	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotExecuteTensorflowSession)
	} else if len(results) == 0 {
		return nil, errors.New(errorTextTensorflowEmptyResponse)
	}

	predictions := results[0].Value().([][]float32)[0]
	className, probability := findClassWithMaxProbability(predictions, s.labels)

	log.Printf("Prediction finished. Predicted class=[%v] with probability=[%v]", className, probability)
	return &Result{Class: className, Probability: probability}, nil
}

func createTensorFlowGraphFromModel(model []byte) (*tf.Graph, error) {
	// Construct an in-memory graph from the serialized form.
	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		return nil, err
	}

	return graph, nil
}

func (s *Service) Stop() error {
	if err := s.session.Close(); err != nil {
		return err
	}
	if err := s.normalizationSession.Close(); err != nil {
		return err
	}

	return nil
}
