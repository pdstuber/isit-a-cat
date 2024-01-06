package postimage

import (
	"io"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pdstuber/isit-a-cat/internal/dep"
)

const (
	headerNameContentType      = "Content-Type"
	headerValueContentTypeJSON = "application/json"
	fileFormKey                = "file"
	staticPictureName          = "picture.jpg"
	errorTextInvalidForm       = "invalid or missing http form data"
	errorTextInvalidFormFile   = "invalid form key. Please provide an image file under they key 'file'"
)

// ImgParams are parameters that identify an image
type ImgParams struct {
	ID           string `json:"id"`
	OriginalName string `json:"originalName"`
}

type handerDependencies interface {
	dep.HasStorageWriter
	dep.HasIDGenerator
}

// Handler handles http requests for uploading images
type Handler struct {
	deps handerDependencies
}

// NewHandler creates a new http handler for uploading images
func NewHandler(deps handerDependencies) *Handler {
	return &Handler{deps}
}

// ServeHTTP requests on the image upload endpoint
func (h *Handler) Handle(c *fiber.Ctx) error {
	log.Println("inside post image handler")
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("invalid http form: %v\n", err)
		return fiber.NewError(fiber.StatusBadRequest, errorTextInvalidForm)
	}
	files := form.File[fileFormKey]

	if len(files) != 1 {
		log.Printf("invalid number of files in form: %d\n", len(files))
		return fiber.NewError(fiber.StatusBadRequest, errorTextInvalidFormFile)
	}

	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Could not extract file from HTTP form: %v\n", err)
		return fiber.NewError(fiber.StatusBadRequest, errorTextInvalidFormFile)
	}

	defer file.Close()
	id := h.deps.IDGenerator().GenerateID()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Could not read image from HTTP form: %v\n", err)
		return fiber.ErrInternalServerError
	}

	err = h.deps.StorageWriter().WriteToBucketObject(id, data)
	if err != nil {
		log.Printf("Could not upload image to object storage: %v\n", err)
		return fiber.ErrInternalServerError
	}
	imgParams := ImgParams{
		ID:           id,
		OriginalName: staticPictureName,
	}

	return c.JSON(imgParams, headerValueContentTypeJSON)
}
