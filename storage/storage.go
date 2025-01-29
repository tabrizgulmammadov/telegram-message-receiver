package storage

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type MessageStorage interface {
	SaveVoiceMessage(chatID int64, username string, reader io.Reader, timestamp time.Time) error
	SaveTextMessage(chatID int64, username string, text string, timestamp time.Time) error
}

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) SaveVoiceMessage(chatID int64, username string, reader io.Reader, timestamp time.Time) error {
	voiceFolder := filepath.Join(s.basePath, "voices", fmt.Sprintf("%d_%s", chatID, username))
	log.Printf("Saving voice message to %s", voiceFolder)
	if err := os.MkdirAll(voiceFolder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filePath := filepath.Join(voiceFolder, fmt.Sprintf("%d.ogg", timestamp.Unix()))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to save voice message: %v", err)
	}

	fmt.Printf("Voice message saved: %s\n", filePath)
	return nil
}

func (s *LocalStorage) SaveTextMessage(chatID int64, username string, text string, timestamp time.Time) error {
	textFolder := filepath.Join(s.basePath, "texts")
	log.Println(textFolder)
	if err := os.MkdirAll(textFolder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filePath := filepath.Join(textFolder, fmt.Sprintf("%d.txt", timestamp.Unix()))
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	messageEntry := fmt.Sprintf("[%s]: %s\n", timestamp.Format(time.RFC3339), text)

	if _, err := file.WriteString(messageEntry); err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	fmt.Printf("Text message saved: %s\n", filePath)
	return nil
}
