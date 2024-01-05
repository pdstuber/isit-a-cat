package getprediction

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/pdstuber/isit-a-cat/internal/service/prediction"
	"github.com/pdstuber/isit-a-cat/pkg/predict"
	"github.com/pkg/errors"
)

const (
	errorTypeServerError = "SERVER_ERROR"
	errorTypeClientError = "CLIENT_ERROR"
	errorTextMissingID   = "request is missing mandatory path parameter 'id'"
)

var serverErrorResponse = ErrorResponse{
	ErrorType: errorTypeServerError,
}

// GetPredictionHandlerImpl handles http requests for predicting images
type Handler struct {
	predictionService *prediction.Service
}

// Result contains the result of a prediction for an image
type Result struct {
	*predict.Result
	error
}

// An ErrorResponse is sent back to the client in case an error occurred
type ErrorResponse struct {
	ErrorType string
	Message   string
}

// NewHandler creates an instance of the prediction handler
func NewHandler(predictionService *prediction.Service) *Handler {
	return &Handler{predictionService: predictionService}
}

// ServeHTTP requests on the image prediction endpoint
func (h *Handler) Handle(ctx *fiber.Ctx) error {
	websocketHandler := websocket.New(func(c *websocket.Conn) {
		id := c.Params("id")

		if id == "" {
			log.Println(errorTextMissingID)
			clientErrorResponse := ErrorResponse{
				ErrorType: errorTypeClientError,
				Message:   errorTextMissingID,
			}
			err := c.WriteJSON(&clientErrorResponse)
			if err != nil {
				log.Printf("error in writing websocket error response: %v\n", err)
			}
			return
		}

		h.triggerPrediction(id, c)
	})
	return websocketHandler(ctx)
}

func (h *Handler) triggerPrediction(id string, ws *websocket.Conn) {
	defer func() {
		err := ws.Close()
		if err != nil {
			log.Printf("Error closing websocket: %v\n", err)
		}
	}()

	predictionResultChannel := make(chan Result)

	go h.getPredictionFromService(id, predictionResultChannel)

	prediction := <-predictionResultChannel

	if prediction.error != nil {
		log.Printf("Error getting predictions: %v\n", prediction.error)
		err := ws.WriteJSON(&serverErrorResponse)
		if err != nil {
			log.Printf("error in writing websocket error response: %v\n", err)
		}
		return
	}
	if err := ws.WriteJSON(prediction.Result); err != nil {
		log.Printf("error writing json response to websocket : %v\n", err)

		err := ws.WriteJSON(&serverErrorResponse)
		if err != nil {
			log.Printf("error in writing websocket error response: %v\n", err)
		}
		return
	}
}

func (h Handler) getPredictionFromService(id string, predictionResultChannel chan Result) {
	imagePrediction, err := h.predictionService.CalculatePredictionForStorageObject(id)

	if err != nil {
		predictionResultChannel <- Result{nil, errors.Wrap(err, "Error getting prediction from prediction service")}
	}

	predictionResultChannel <- Result{imagePrediction, nil}
}
