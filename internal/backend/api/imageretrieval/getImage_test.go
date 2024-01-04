// handlers_test.go
package imageretrieval

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/pdstuber/isit-a-cat-bff/imageretrieval/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	getImageURL = "/images"
	testID      = "12345"
)

var (
	storageServiceMockResponse []byte = []byte{1, 2, 3}
	errMock                           = errors.New("everything went to hell")
)

func Test_ServeHTTP_good_case(t *testing.T) {
	storageServiceMock := new(mocks.StorageService)

	storageServiceMock.On("ReadFromBucketObject", mock.Anything).Return(storageServiceMockResponse, nil)

	getImageHandler := NewHandler(storageServiceMock)

	req, err := http.NewRequest("GET", fmt.Sprintf(getImageURL+"/%s", testID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(getImageURL+"/{id}", getImageHandler.ServeHTTP)
	router.ServeHTTP(rr, req)

	storageServiceMock.AssertNumberOfCalls(t, "ReadFromBucketObject", 1)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, headerValueContentTypeJpeg, rr.Header().Get(headerNameContentType))
	assert.Equal(t, storageServiceMockResponse, rr.Body.Bytes())
}

func Test_ServeHTTP_id_missing(t *testing.T) {
	storageServiceMock := new(mocks.StorageService)

	storageServiceMock.On("ReadFromBucketObject", mock.Anything).Return(storageServiceMockResponse, nil)

	getImageHandler := NewHandler(storageServiceMock)

	req, err := http.NewRequest("GET", getImageURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(getImageURL, getImageHandler.ServeHTTP)
	router.ServeHTTP(rr, req)

	storageServiceMock.AssertNumberOfCalls(t, "ReadFromBucketObject", 0)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, errorTextMissingID+"\n", rr.Body.String())
}

func Test_ServeHTTP_error_storage_service(t *testing.T) {
	storageServiceMock := new(mocks.StorageService)

	storageServiceMock.On("ReadFromBucketObject", mock.Anything).Return(nil, errMock)

	getImageHandler := NewHandler(storageServiceMock)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", fmt.Sprintf(getImageURL+"/%s", testID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(getImageURL+"/{id}", getImageHandler.ServeHTTP)
	router.ServeHTTP(rr, req)

	storageServiceMock.AssertNumberOfCalls(t, "ReadFromBucketObject", 1)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError)+"\n", rr.Body.String())

}
