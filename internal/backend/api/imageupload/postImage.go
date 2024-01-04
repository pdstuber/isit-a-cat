package imageupload

import (
	"io"
	"log"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/pdstuber/isit-a-cat-bff/imageretrieval"
)

const (
	headerNameContentType      = "Content-Type"
	headerValueContentTypeJSON = "application/json"
	fileFormKey                = "file"
	staticPictureName          = "picture.jpg"
	errorTextInvalidForm       = "invalid or missing http form data"
	errorTextInvalidFormFile   = "invalid form key. Please provide an image file under they key 'file'"
)

// A StorageWriter writes data to the storage object with the given ID
type StorageWriter interface {
	WriteToBucketObject(objectID string, data []byte) error
}

// A IDGenerator generates unique ids as string
type IDGenerator interface {
	GenerateID() string
}

// Handler handles http requests for uploading images
type Handler struct {
	storageWriter StorageWriter
	idGenerator   IDGenerator
}

// NewHandler creates a new http handler for uploading images
func NewHandler(storageWriter StorageWriter, idGenerator IDGenerator) *Handler {
	return &Handler{storageWriter, idGenerator}
}

// ServeHTTP requests on the image upload endpoint
func (h *Handler) Handle(c *fiber.Ctx) error {
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

	id := h.idGenerator.GenerateID()

	data, err := io.ReadAll(file)

	if err != nil {
		log.Printf("Could not read image from HTTP form: %v\n", err)
		return fiber.ErrInternalServerError
	}

	err = h.storageWriter.WriteToBucketObject(id, data)

	if err != nil {
		log.Printf("Could not upload image to object storage: %v\n", err)
		return fiber.ErrInternalServerError
	}
	imgParams := imageretrieval.ImgParams{
		ID:           id,
		OriginalName: staticPictureName,
	}

	return c.JSON(imgParams, headerValueContentTypeJSON)
}
