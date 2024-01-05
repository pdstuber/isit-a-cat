package getprediction

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/pdstuber/isit-a-cat-bff/prediction/mocks"
	"gitlab.com/pdstuber/isit-a-cat-bff/prediction/predictionservice"
)

const (
	predictionURL = "/predictions"
	testID        = "12345"
	mockErrorText = "everything went to hell"
)

var (
	mockPrediction = predictionservice.Prediction{
		Class:       "cat",
		Probability: 0.99,
	}
	errMock = errors.New(mockErrorText)
)

func Test_ServeHTTP_good_case(t *testing.T) {

	predictionServiceMock := new(mocks.PredictionService)

	predictionServiceMock.On("GetPredictionForImageID", mock.Anything).Return(&mockPrediction, nil)
	predictionHandler := NewHandler(predictionServiceMock, &websocket.Upgrader{})

	router := mux.NewRouter()
	router.HandleFunc(predictionURL+"/{id}", predictionHandler.ServeHTTP)

	s := httptest.NewServer(router)
	defer s.Close()

	url := fmt.Sprintf("%s%s%s", "ws", strings.TrimPrefix(s.URL, "http"), fmt.Sprintf(predictionURL+"/%s", testID))

	ws, _, _ := websocket.DefaultDialer.Dial(url, nil)

	defer ws.Close()

	prediction := predictionservice.Prediction{}

	_ = ws.ReadJSON(&prediction)

	assert.Equal(t, mockPrediction, prediction)
}

func Test_ServeHTTP_error_prediction_service(t *testing.T) {
	predictionServiceMock := new(mocks.PredictionService)

	predictionServiceMock.On("GetPredictionForImageID", mock.Anything).Return(nil, errMock)
	predictionHandler := NewHandler(predictionServiceMock, &websocket.Upgrader{})

	router := mux.NewRouter()
	router.HandleFunc(predictionURL+"/{id}", predictionHandler.ServeHTTP)

	s := httptest.NewServer(router)
	defer s.Close()

	url := fmt.Sprintf("%s%s%s", "ws", strings.TrimPrefix(s.URL, "http"), fmt.Sprintf(predictionURL+"/%s", testID))

	ws, _, _ := websocket.DefaultDialer.Dial(url, nil)

	defer ws.Close()

	errorResponse := ErrorResponse{}

	_ = ws.ReadJSON(&errorResponse)

	assert.Equal(t, errorTypeServerError, errorResponse.ErrorType)
}

func Test_ServeHTTP_missing_id(t *testing.T) {
	predictionServiceMock := new(mocks.PredictionService)

	predictionServiceMock.On("GetPredictionForImageID", mock.Anything).Return(&mockPrediction, nil)
	predictionHandler := NewHandler(predictionServiceMock, &websocket.Upgrader{})

	router := mux.NewRouter()
	router.HandleFunc(predictionURL, predictionHandler.ServeHTTP)

	s := httptest.NewServer(router)
	defer s.Close()

	url := fmt.Sprintf("%s%s%s", "ws", strings.TrimPrefix(s.URL, "http"), predictionURL)

	ws, _, _ := websocket.DefaultDialer.Dial(url, nil)

	defer ws.Close()

	errorResponse := ErrorResponse{}

	_ = ws.ReadJSON(&errorResponse)

	assert.Equal(t, errorTypeClientError, errorResponse.ErrorType)
	assert.Equal(t, errorTextMissingID, errorResponse.Message)
}
