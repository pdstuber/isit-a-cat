/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocarina/gocsv"
	"github.com/pdstuber/isit-a-cat/internal/bot"
	"github.com/pdstuber/isit-a-cat/internal/pkg/predict"
	"github.com/spf13/cobra"
)

// botCmd represents the bot command
var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Start the telegram bot.",
	Run: func(cmd *cobra.Command, args []string) {
		botAPI, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
		if err != nil {
			log.Panic(err)
		}

		botAPI.Debug = true

		log.Printf("Authorized on account %s", botAPI.Self.UserName)
		var labels []predict.Label

		modelPath := os.Getenv("MODEL_PATH")
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

		bot := bot.New(botAPI, predict.New(model, labels, defaultColorChannels))

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
