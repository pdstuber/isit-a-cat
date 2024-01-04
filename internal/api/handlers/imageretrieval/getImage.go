package imageretrieval

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	headerNameContentType      = "Content-Type"
	headerValueContentTypeJpeg = "image/jpeg"
	errorTextMissingID         = "request is missing mandatory path parameter 'id'"
)

// A StorageReader reads data from the storage object with the given ID
type StorageReader interface {
	ReadFromBucketObject(objectID string) ([]byte, error)
}

// ImgParams are parameters that identify an image
type ImgParams struct {
	ID           string `json:"id"`
	OriginalName string `json:"originalName"`
}

// GetImageHandlerImpl handles http requests for retrieving images
type GetImageHandlerImpl struct {
	storageReader StorageReader
}

// NewHandler creates an instance of the get image handler
func NewHandler(storageReader StorageReader) http.Handler {
	return &GetImageHandlerImpl{storageReader}
}

// ServeHTTP requests for getting images
func (h GetImageHandlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	if id == "" {
		log.Println(errorTextMissingID)
		http.Error(w, errorTextMissingID, http.StatusBadRequest)
		return
	}

	image, err := h.storageReader.ReadFromBucketObject(id)

	if err != nil {
		log.Printf("Error retrieving image from object storage: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set(headerNameContentType, headerValueContentTypeJpeg)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(image)

	if err != nil {
		log.Printf("could not write response: %v\n", err)
	}
}
