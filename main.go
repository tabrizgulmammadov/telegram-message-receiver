package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"telegram-message-receiver/config"
	"telegram-message-receiver/handler"
	"telegram-message-receiver/logger"
	"telegram-message-receiver/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	logger := logger.NewLogger(config.Debug)

	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		logger.Error("Error starting bot: %v", err)
		os.Exit(1)
	}

	bot.Debug = config.Debug

	storage := storage.NewLocalStorage(config.StoragePath)
	handler := handler.NewMessageHandler(bot, config, storage, logger)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	// Graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Shutting down gracefully...")
		os.Exit(0)
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if err := handler.HandleMessage(update.Message); err != nil {
			logger.Error("Error handling message: %v", err)
		}
	}
}
