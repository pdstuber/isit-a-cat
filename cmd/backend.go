package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pdstuber/isit-a-cat/internal/api"
	"github.com/pdstuber/isit-a-cat/internal/dep"
	"github.com/pdstuber/isit-a-cat/internal/service/idgenerator"
	"github.com/pdstuber/isit-a-cat/internal/service/storage"
	"github.com/pdstuber/isit-a-cat/pkg/prediction"
	"github.com/spf13/cobra"
)

// backendCmd represents the backend command
var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Start the backend to process API requests from the frontend",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := api.ConfigFromEnv()
		if err != nil {
			log.Fatalf("could not create config from environment: %v\n", err)
		}

		imagePredictor := prediction.NewService(config.Model, config.Labels, defaultColorChannels, config.TFInputOperationName, config.TFOutputOperationName, config.TargetImageDimensions)

		storageService, err := storage.New(config.ObjectStorageBucketName, config.ObjectStorageObjectFolder, config.ObjectStorageEndpoint, config.ObjectStorageAccessKeyID, config.ObjectStorageSecretAccessKey, config.ObjectStorageUseTLS)
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
