package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func (fo *FileOrganizer) initLog() error {
	file, err := os.OpenFile("organizer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	fo.logFile = file
	log.SetOutput(fo.logFile)

	return nil
}

func (fo *FileOrganizer) logSuccess(message string) {
	log.Printf("[SUCCESS] %s\n", message)
}

func (fo *FileOrganizer) logError(message string) {
	log.Printf("[ERROR] %s\n", message)
}

func (fo *FileOrganizer) Close() error {
	if fo.logFile != nil {
		return fo.logFile.Close()
	}
	return nil
}

func (fo *FileOrganizer) moveFile(sourcePath, targetDir string) error {
	fullPath := filepath.Join(fo.sourceDir, targetDir)
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		msg := fmt.Sprintf("не удалось создать директорию %s: %v", targetDir, err)
		fo.logError(msg)
		return fmt.Errorf("moveFile: %w", err)
	}

	nameFile := filepath.Base(sourcePath)

	targetFullPath := filepath.Join(fullPath, nameFile)

	_, err = os.Stat(targetFullPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			msg := fmt.Sprintf("не удалось проверить файл %q: %v", targetFullPath, err)
			fo.logError(msg)
			return fmt.Errorf("moveFile: %w", err)
		}
	} else {
		ext := filepath.Ext(nameFile)
		base := strings.TrimSuffix(nameFile, ext)
		newNameFile := base + "_" + time.Now().Format("2006-01-02_15-04-05") + ext
		targetFullPath = filepath.Join(fullPath, newNameFile)
	}
	err = os.Rename(sourcePath, targetFullPath)
	if err != nil {
		msg := fmt.Sprintf("не удалось переместить файл %q: %v", nameFile, err)
		fo.logError(msg)
		return fmt.Errorf("moveFile: %w", err)
	}
	msg := fmt.Sprintf("файл %q перемещён в директорию %q", nameFile, targetDir)
	fo.logSuccess(msg)
	return nil
}

func (fo *FileOrganizer) Organize() error {
	err := fo.initLog()
	if err != nil {
		return err
	}

	err = filepath.WalkDir(fo.sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Dir(path) != fo.sourceDir {
			return nil
		}

		currentDir, ok := fo.rulesMap[strings.ToLower(filepath.Ext(path))]
		if ok {
			err := fo.moveFile(path, currentDir)
			if err != nil {
				return nil
			}
			fo.processedFiles++
		}

		return nil
	})
	return err
}
