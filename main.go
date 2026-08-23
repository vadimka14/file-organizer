package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type FileOrganizer struct {
	sourceDir      string
	rulesMap       map[string]string
	processedFiles int
	logFile        *os.File
}

func main() {
	var DefaultRules = map[string]string{
		".jpg":  "Images",
		".jpeg": "Images",
		".png":  "Images",
		".pdf":  "Documents",
		".doc":  "Documents",
		".docx": "Documents",
		".txt":  "Documents",
		".mp3":  "Music",
		".wav":  "Music",
		".mp4":  "Video",
		".avi":  "Video",
		".zip":  "Archives",
		".rar":  "Archives",
	}

	var dir string = "./banks"

	for key, value := range DefaultRules {
		fmt.Printf("Расширение: %s -> Папка: %s\n", key, value)
	}

	_, err := NewFileOrganizer(dir)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
	fmt.Printf("FileOrganizer создан для директории: %s\n", dir)

}

func NewFileOrganizer(sourceDir string) (*FileOrganizer, error) {
	if len(sourceDir) == 0 {
		return nil, fmt.Errorf("директория не указана")
	}
	info, err := os.Stat(sourceDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("директория не найдена: %w", err)
		}
		return nil, fmt.Errorf("директория не найдена: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("директория не найдена: %s является файлом", info.Name())
	}

	return &FileOrganizer{
		sourceDir:      sourceDir,
		rulesMap:       map[string]string{},
		processedFiles: 0,
		logFile:        nil,
	}, nil
}
