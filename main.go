package main

import (
	"errors"
	"fmt"
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

	for key, value := range DefaultRules {
		fmt.Printf("Расширение: %s -> Папка: %s\n", key, value)
	}

}

func NewFileOrganizer(sourceDir string) (*FileOrganizer, error) {
	if len(sourceDir) == 0 {
		return nil, fmt.Errorf("")
	}
	info, err := os.Stat(sourceDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Файла нет")
		}
		return nil, err
	}

	if info.IsDir() {
		fmt.Println("dir:", info.Name())
	} else {
		return nil, fmt.Errorf("Директорію не знайдено: %s є файлом", info.Name())
	}

	return &FileOrganizer{
		sourceDir:      sourceDir,
		rulesMap:       map[string]string{},
		processedFiles: 0,
		logFile:        nil,
	}, nil
}
