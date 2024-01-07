package bot

import (
	"context"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pdstuber/isit-a-cat/pkg/prediction"
)

const telegramBotErrorMessage = "there was a problem in processing your request at this time"

// A ImagePredictor predicts the class of an image
type ImagePredictor interface {
	PredictImage(imageBytes []byte) (*prediction.Result, error)
	Stop() error
}

type Bot struct {
	botAPI          *tgbotapi.BotAPI
	wg              *sync.WaitGroup
	workers         int
	fetchBuffer     int
	shutdownChannel chan interface{}
	imagePredictor  ImagePredictor
	httpClient      *http.Client
}

func New(botAPI *tgbotapi.BotAPI, imagePredictor ImagePredictor) *Bot {
	return &Bot{
		botAPI:          botAPI,
		wg:              &sync.WaitGroup{},
		workers:         4,
		fetchBuffer:     100,
		shutdownChannel: make(chan interface{}),
		httpClient:      http.DefaultClient,
		imagePredictor:  imagePredictor,
	}
}

func (b *Bot) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	ch := make(chan *tgbotapi.Update, b.fetchBuffer)
	go b.FetchAsync(ctx, u, ch)
	b.StartWorkers(ctx, ch)
}

func (b *Bot) Stop() {
	close(b.shutdownChannel)
	b.wg.Wait()
	b.imagePredictor.Stop()
}

func (bot *Bot) FetchAsync(ctx context.Context, config tgbotapi.UpdateConfig, ch chan *tgbotapi.Update) {
	for {
		select {
		case <-bot.shutdownChannel:
			close(ch)
			return
		case <-ctx.Done():
			close(ch)
			return
		default:
			updates, err := bot.botAPI.GetUpdates(config)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					ch <- &update
				}
			}
		}
	}
}

func (b *Bot) StartWorkers(ctx context.Context, ch chan *tgbotapi.Update) {
	for i := 0; i < b.workers; i++ {
		b.wg.Add(1)

		go func() {
			for update := range ch {
				if update.Message == nil {
					continue
				}

				if len(update.Message.Photo) == 0 {
					continue
				}

				msg := b.handlePhoto(ctx, update.Message)

				if _, err := b.botAPI.Send(msg); err != nil {
					log.Println(err)
				}
			}
			b.wg.Done()
		}()
	}
}

// TODO improve error messages
func (b *Bot) handlePhoto(ctx context.Context, message *tgbotapi.Message) tgbotapi.MessageConfig {
	fileConfig := tgbotapi.FileConfig{
		FileID: message.Photo[2].FileID,
	}

	file, err := b.botAPI.GetFile(fileConfig)
	if err != nil {
		log.Printf("could not retrieve information about your uploaded photo from the server: %w\n", err)
		return tgbotapi.NewMessage(message.Chat.ID, telegramBotErrorMessage)
	}

	link := file.Link(b.botAPI.Token)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		log.Printf("could not create http request: %w\n", err)
		return tgbotapi.NewMessage(message.Chat.ID, telegramBotErrorMessage)
	}

	response, err := b.httpClient.Do(request)
	if err != nil {
		log.Printf("could not perform http request: %w\n", err)
		return tgbotapi.NewMessage(message.Chat.ID, telegramBotErrorMessage)
	}

	defer response.Body.Close()

	photoBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("could read http response body: %w\n", err)
		return tgbotapi.NewMessage(message.Chat.ID, telegramBotErrorMessage)
	}

	result, err := b.imagePredictor.PredictImage(photoBytes)
	if err != nil {
		log.Printf("could not retrieve information about your uploaded photo from the server: %w\n", err)
		return tgbotapi.NewMessage(message.Chat.ID, telegramBotErrorMessage)
	}

	return tgbotapi.NewMessage(message.Chat.ID, result.String())
}
