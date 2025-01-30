package storage

import (
	"encoding/json"
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
	SaveContactInfo(chatID int64, username string, phoneNumber string, timestamp time.Time) error
	HasContactInfo(chatID int64) (bool, error)
}

type ContactInfo struct {
	Username    string    `json:"username"`
	PhoneNumber string    `json:"phone_number"`
	Timestamp   time.Time `json:"timestamp"`
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

func (s *LocalStorage) HasContactInfo(chatID int64) (bool, error) {
	filePath := filepath.Join(s.basePath, "contacts", fmt.Sprintf("%d.json", chatID))
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking contact info: %v", err)
	}
	return true, nil
}

func (s *LocalStorage) SaveContactInfo(chatID int64, username string, phoneNumber string, timestamp time.Time) error {
	contactsFolder := filepath.Join(s.basePath, "contacts")
	if err := os.MkdirAll(contactsFolder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	contactInfo := ContactInfo{
		Username:    username,
		PhoneNumber: phoneNumber,
		Timestamp:   timestamp,
	}

	filePath := filepath.Join(contactsFolder, fmt.Sprintf("%d.json", chatID))

	data, err := json.MarshalIndent(contactInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal contact info: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save contact info: %v", err)
	}

	return nil
}
