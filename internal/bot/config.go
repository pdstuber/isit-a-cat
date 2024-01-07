package bot

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/pdstuber/isit-a-cat/pkg/prediction"
	"github.com/pkg/errors"
)

type Config struct {
	TelegramBotToken      string
	Labels                []prediction.Label
	Model                 []byte
	TargetImageDimensions int
	TFInputOperationName  string
	TFOutputOperationName string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ConfigFromEnv() (*Config, error) {
	telegramBotToken := getEnv("TELEGRAM_BOT_TOKEN", "")
	if telegramBotToken == "" {
		return nil, errors.New("telegram bot token is mandatory")
	}

	var labels []prediction.Label

	modelPath := getEnv("MODEL_PATH", "/model")
	model, err := os.ReadFile(fmt.Sprintf("%s/model.pb", modelPath))
	if err != nil {
		return nil, errors.Wrap(err, "could not read model")
	}
	labelBytes, err := os.ReadFile(fmt.Sprintf("%s/labels.csv", modelPath))
	if err != nil {
		return nil, errors.Wrap(err, "could not read labels")
	}

	if err := gocsv.UnmarshalBytes(labelBytes, &labels); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal labels csv")
	}

	targetImageDimensions, err := strconv.Atoi(getEnv("TARGET_IMAGE_DIMENSIONS", "256"))
	if err != nil {
		return nil, errors.Wrap(err, "could not convert to integer, please use correct format")
	}

	inputOperationName := getEnv("TF_INPUT_OPERATION_NAME", "input_1")
	outputOperationName := getEnv("TF_OUTPUT_OPERATION_NAME", "dense_3/Softmax")

	return &Config{
		TelegramBotToken:      telegramBotToken,
		Labels:                labels,
		Model:                 model,
		TargetImageDimensions: targetImageDimensions,
		TFInputOperationName:  inputOperationName,
		TFOutputOperationName: outputOperationName,
	}, nil
}
