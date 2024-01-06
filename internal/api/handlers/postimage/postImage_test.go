package postimage_test

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/pdstuber/isit-a-cat-bff/imageretrieval"
	"gitlab.com/pdstuber/isit-a-cat-bff/imageupload/mocks"
)

const (
	postImageURL       = "/images"
	mockID             = "12346"
	mockErrorText      = "everything went to hell"
	invalidFormFileKey = "image"
)

var (
	mockImage     = []byte{9, 10, 17, 12, 16}
	errMock       = errors.New(mockErrorText)
	mockImgParams = imageretrieval.ImgParams{
		ID:           mockID,
		OriginalName: staticPictureName,
	}
)

func Test_ServeHTTP_good_case(t *testing.T) {
	storageServiceMock := new(mocks.StorageServiceMock)
	storageServiceMock.On("WriteToBucketObject", mock.Anything, mock.Anything).Return(nil)

	idGenerator := new(mocks.IDGenerator)
	idGenerator.On("GenerateID", mock.Anything).Return(mockID)

	postImageHandler := NewHandler(storageServiceMock, idGenerator)

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	go func() {
		defer writer.Close()
		part, _ := writer.CreateFormFile("file", "image.jpg")
		_, _ = part.Write(mockImage)
	}()

	req, err := http.NewRequest("POST", postImageURL, pr)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(postImageURL, postImageHandler.ServeHTTP)
	router.ServeHTTP(rr, req)

	storageServiceMock.AssertNumberOfCalls(t, "WriteToBucketObject", 1)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, headerValueContentTypeJSON, rr.Header().Get(headerNameContentType))

	expectedImgParams, _ := json.Marshal(mockImgParams)
	assert.Equal(t, expectedImgParams, rr.Body.Bytes())
}

func Test_ServeHTTP_Form_File_Error(t *testing.T) {
	storageServiceMock := new(mocks.StorageServiceMock)
	storageServiceMock.On("WriteToBucketObject", mock.Anything, mock.Anything).Return(nil)

	idGenerator := new(mocks.IDGenerator)
	idGenerator.On("GenerateID", mock.Anything).Return(mockID)

	postImageHandler := NewHandler(storageServiceMock, idGenerator)

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	go func() {
		defer writer.Close()
		part, _ := writer.CreateFormFile(invalidFormFileKey, "image.jpg")
		_, _ = part.Write(mockImage)
	}()

	req, err := http.NewRequest("POST", postImageURL, pr)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(postImageURL, postImageHandler.ServeHTTP)
	router.ServeHTTP(rr, req)

	storageServiceMock.AssertNumberOfCalls(t, "WriteToBucketObject", 0)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, errorTextInvalidFormFile+"\n", rr.Body.String())
}

func Test_ServeHTTP_multi_part_form_parsing_error(t *testing.T) {
	storageServiceMock := new(mocks.StorageServiceMock)
	storageServiceMock.On("WriteToBucketObject", mock.Anything, mock.Anything).Return(nil)

	idGenerator := new(mocks.IDGenerator)
	idGenerator.On("GenerateID", mock.Anything).Return(mockID)

	postImageHandler := NewHandler(storageServiceMock, idGenerator)

	req, err := http.NewRequest("POST", postImageURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(postImageURL, postImageHandler.ServeHTTP)
	router.ServeHTTP(rr, req)

	storageServiceMock.AssertNumberOfCalls(t, "WriteToBucketObject", 0)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, errorTextInvalidForm+"\n", rr.Body.String())
}

func Test_ServeHTTP_error_storage_service(t *testing.T) {
	storageServiceMock := new(mocks.StorageServiceMock)
	storageServiceMock.On("WriteToBucketObject", mock.Anything, mock.Anything).Return(errMock)

	idGenerator := new(mocks.IDGenerator)
	idGenerator.On("GenerateID", mock.Anything).Return(mockID)

	postImageHandler := NewHandler(storageServiceMock, idGenerator)

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	go func() {
		defer writer.Close()
		part, _ := writer.CreateFormFile("file", "image.jpg")
		_, _ = part.Write(mockImage)
	}()

	req, err := http.NewRequest("POST", postImageURL, pr)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(postImageURL, postImageHandler.ServeHTTP)
	router.ServeHTTP(rr, req)

	storageServiceMock.AssertNumberOfCalls(t, "WriteToBucketObject", 1)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError)+"\n", rr.Body.String())
}
