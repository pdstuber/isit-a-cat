package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pdstuber/isit-a-cat/internal/bot"
	"github.com/pdstuber/isit-a-cat/pkg/prediction"
	"github.com/spf13/cobra"
)

// botCmd represents the bot command
var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Start the telegram bot.",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := bot.ConfigFromEnv()
		if err != nil {
			log.Fatalf("could not create config from environment: %v\n", err)
		}

		botAPI, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
		if err != nil {
			log.Panic(err)
		}

		botAPI.Debug = true

		imagePredictor := prediction.NewService(config.Model, config.Labels, defaultColorChannels, config.TFInputOperationName, config.TFOutputOperationName, config.TargetImageDimensions)

		bot := bot.New(botAPI, imagePredictor)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		bot.Start(ctx)

		log.Println("after start")

		<-ctx.Done()
		bot.Stop()
	},
}

func init() {
	runCmd.AddCommand(botCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// botCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// botCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
