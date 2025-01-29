// handler/handler.go
package handler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-message-receiver/config"
	"telegram-message-receiver/logger"
	"telegram-message-receiver/storage"
)

type MessageHandler struct {
	bot     *tgbotapi.BotAPI
	config  *config.Config
	storage storage.MessageStorage
	logger  *logger.Logger
}

func NewMessageHandler(bot *tgbotapi.BotAPI, config *config.Config, storage storage.MessageStorage, logger *logger.Logger) *MessageHandler {
	return &MessageHandler{
		bot:     bot,
		config:  config,
		storage: storage,
		logger:  logger,
	}
}

func (h *MessageHandler) HandleMessage(message *tgbotapi.Message) error {
	if message == nil {
		return fmt.Errorf("received nil message")
	}

	username := h.sanitizeUsername(message.From.UserName)
	timestamp := time.Now()

	// Log message receipt
	h.logger.Debug("Received message from %s (Chat ID: %d)", username, message.Chat.ID)

	switch {
	case message.Voice != nil:
		return h.handleVoiceMessage(message.Chat.ID, username, message.Voice, timestamp)
	case message.Text != "":
		return h.handleTextMessage(message.Chat.ID, username, message.Text, timestamp)
	default:
		h.logger.Info("Unsupported message type received from %s", username)
		return nil
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, h.config.AcknowledgmentMessage)
	_, err := h.bot.Send(msg)
	if err != nil {
		h.logger.Error("Failed to send confirmation message: %v", err)
		return err
	}

	return nil
}

func (h *MessageHandler) handleVoiceMessage(chatID int64, username string, voice *tgbotapi.Voice, timestamp time.Time) error {
	h.logger.Debug("Processing voice message from %s (Duration: %d seconds)", username, voice.Duration)

	file, err := h.downloadFile(voice.FileID)
	if err != nil {
		return fmt.Errorf("failed to download voice message: %w", err)
	}
	defer file.Close()

	return h.storage.SaveVoiceMessage(chatID, username, file, timestamp)
}

func (h *MessageHandler) handleTextMessage(chatID int64, username, text string, timestamp time.Time) error {
	h.logger.Debug("Processing text message from %s (Length: %d)", username, len(text))

	// Filter out any potentially harmful characters from text
	sanitizedText := h.sanitizeText(text)

	if err := h.storage.SaveTextMessage(chatID, username, sanitizedText, timestamp); err != nil {
		return fmt.Errorf("failed to save text message: %w", err)
	}

	// If configured, send acknowledgment
	if h.config.SendAcknowledgment {
		if err := h.sendAcknowledgment(chatID); err != nil {
			h.logger.Error("Failed to send acknowledgment: %v", err)
			// Don't return error as the message was saved successfully
		}
	}

	return nil
}

func (h *MessageHandler) downloadFile(fileID string) (io.ReadCloser, error) {
	file, err := h.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	fileURL := fmt.Sprintf(h.config.BaseFileURL, h.bot.Token, file.FilePath)

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (h *MessageHandler) sanitizeUsername(username string) string {
	if username == "" {
		return "anonymous"
	}
	// Remove any potentially harmful characters
	username = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, username)
	return username
}

func (h *MessageHandler) sanitizeText(text string) string {
	// Remove any null bytes and other potentially harmful characters
	return strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, text)
}

func (h *MessageHandler) sendAcknowledgment(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, h.config.AcknowledgmentMessage)
	_, err := h.bot.Send(msg)
	return err
}

// ValidateFileType checks if the file extension is allowed
func (h *MessageHandler) validateFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedTypes := map[string]bool{
		".ogg": true,
		".mp3": true,
		".wav": true,
		".m4a": true,
	}
	return allowedTypes[ext]
}
