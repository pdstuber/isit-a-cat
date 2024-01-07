package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/pdstuber/isit-a-cat/internal/api"
	"github.com/pdstuber/isit-a-cat/internal/dep"
	"github.com/pdstuber/isit-a-cat/internal/service/idgenerator"
	"github.com/pdstuber/isit-a-cat/internal/service/storage"
	"github.com/pdstuber/isit-a-cat/pkg/prediction"
	"github.com/spf13/cobra"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// backendCmd represents the backend command
var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Start the backend to process API requests from the frontend",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO parse flags

		var labels []prediction.Label

		modelPath := getEnv("MODEL_PATH", "/model")
		model, err := os.ReadFile(fmt.Sprintf("%s/model.pb", modelPath))
		if err != nil {
			log.Panic(err)
		}
		labelBytes, err := os.ReadFile(fmt.Sprintf("%s/labels.csv", modelPath))
		if err != nil {
			log.Panic(err)
		}

		if err := gocsv.UnmarshalBytes(labelBytes, &labels); err != nil {
			log.Fatalf("could not unmarshal labels csv: %v\n", err)
		}

		// TODO move to environment vars
		targetImageDimensions := 256
		inputOperationName := "input_1"
		outputOperationName := "dense_3/Softmax"

		imagePredictor := prediction.NewService(model, labels, defaultColorChannels, inputOperationName, outputOperationName, targetImageDimensions)

		objectStorageEndpoint := getEnv("OBJECT_STORAGE_ENDPOINT", "minio:9000")
		objectStorageAccessKeyID := getEnv("MINIO_ACCESS_KEY", "")
		objectStorageSecretAccessKey := getEnv("MINIO_SECRET_KEY", "")
		objectStorageUseTLS, _ := strconv.ParseBool(getEnv("OBJECT_STORAGE_USE_TLS", "false"))

		storageBucketName := getEnv("STORAGE_BUCKET_NAME", "isit-a-cat")
		storageObjectFolder := getEnv("STORAGE_OBJECT_FOLDER", "uploaded-images/")

		storageService, err := storage.New(storageBucketName, storageObjectFolder, objectStorageEndpoint, objectStorageAccessKeyID, objectStorageSecretAccessKey, objectStorageUseTLS)
		if err != nil {
			log.Fatalf("could not create storage service: %v\n", err)
		}

		deps := dep.NewAppDependencies().
			WithStorageService(storageService).
			WithIDGenerator(&idgenerator.Service{}).
			WithImagePredictor(imagePredictor)

		router := api.NewRouter(deps.Forward(), ":8080")

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		router.Start(ctx)

		<-ctx.Done()
		router.Stop(5 * time.Second)
	},
}

func init() {
	runCmd.AddCommand(backendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
