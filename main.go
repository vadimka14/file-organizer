package main

import (
	"errors"
	"fmt"
	"os"
)

type FileOrganizer struct {
	sourceDir     string
	rulesMap      map[string]string
	ProcessedFile int
	logFile       *os.File
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
		fmt.Println("file:", info.Name())
	}

	return &FileOrganizer{}, nil
}
