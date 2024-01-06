package imageretrieval

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pdstuber/isit-a-cat/internal/dep"
)

const (
	headerNameContentType      = "Content-Type"
	headerValueContentTypeJpeg = "image/jpeg"
	errorTextMissingID         = "request is missing mandatory path parameter 'id'"
)

type handlerDependencies interface {
	dep.HasStorageReader
}

// ImgParams are parameters that identify an image
type ImgParams struct {
	ID           string `json:"id"`
	OriginalName string `json:"originalName"`
}

// GetImageHandlerImpl handles http requests for retrieving images
type Handler struct {
	deps handlerDependencies
}

// NewHandler creates an instance of the get image handler
func NewHandler(deps handlerDependencies) *Handler {
	return &Handler{deps}
}

// ServeHTTP requests for getting images
func (h *Handler) Handle(c *fiber.Ctx) error {

	id := c.Params("id")

	if id == "" {
		log.Println(errorTextMissingID)
		return fiber.NewError(fiber.StatusBadRequest, errorTextMissingID)
	}

	image, err := h.deps.StorageReader().ReadFromBucketObject(id)

	if err != nil {
		log.Printf("Error retrieving image from object storage: %v\n", err)
		return fiber.ErrInternalServerError
	}

	c.Set(headerNameContentType, headerValueContentTypeJpeg)
	return c.Send(image)
}
