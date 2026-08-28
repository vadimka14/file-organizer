package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math"
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
	statistics     map[string]*FileStats
	totalSize      int64
}

type FileStats struct {
	countFiles int
	sumSize    int64
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

	fmt.Println("=== Файловый органайзер ===")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите путь к директории для организации (Enter для текущей директории): ")
	input, _ := reader.ReadString('\n')
	sourcePath := strings.TrimSpace(input)
	if len(sourcePath) == 0 {
		sourcePath, _ = os.Getwd()
	}

	organizer, err := NewFileOrganizer(sourcePath, DefaultRules)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
	defer organizer.Close()

	fmt.Println("Начинаем организацию файлов...")
	err = organizer.Organize()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(organizer.generateReport())

	fmt.Println("Организация завершена! Подробности в файле organizer.log")
}

func NewFileOrganizer(sourceDir string, rulesMap map[string]string) (*FileOrganizer, error) {
	statistics := make(map[string]*FileStats)
	if len(sourceDir) == 0 {
		return nil, fmt.Errorf("директория не указана")
	}
	info, err := os.Stat(sourceDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("директория не найдена: %w", err)
		}
		return nil, fmt.Errorf("не удалось проверить директорию: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("директория не найдена: %s является файлом", info.Name())
	}

	return &FileOrganizer{
		sourceDir:      sourceDir,
		rulesMap:       rulesMap,
		processedFiles: 0,
		logFile:        nil,
		statistics:     statistics,
		totalSize:      0,
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

func (fo *FileOrganizer) moveFile(sourcePath, targetDir string, info os.FileInfo) error {
	sizeOfFile := info.Size()

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
	_, ok := fo.statistics[targetDir]
	if !ok {
		fo.statistics[targetDir] = &FileStats{countFiles: 0, sumSize: 0}
	}
	fo.statistics[targetDir].countFiles++
	fo.statistics[targetDir].sumSize += sizeOfFile
	fo.totalSize += sizeOfFile
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
			info, err := d.Info()
			if err != nil {
				return nil
			}
			err = fo.moveFile(path, currentDir, info)
			if err != nil {
				return nil
			}
			fo.processedFiles++
		}

		return nil
	})
	return err
}

func (fs *FileStats) String() string {
	return fmt.Sprintf("Файлов: %d, Размер: %.2f KB", fs.countFiles, math.Floor((float64(fs.sumSize)/1024)*100)/100)
}

func (fo *FileOrganizer) generateReport() string {
	var sb strings.Builder
	sb.WriteString("=== Отчёт о перемещении файлов ===\n\n")
	sb.WriteString(fmt.Sprintf("Всего обработано файлов: %d\n", fo.processedFiles))
	sb.WriteString(fmt.Sprintf("Общий размер: %.2f KB\n", math.Floor((float64(fo.totalSize)/1024)*100)/100))
	for category, stats := range fo.statistics {
		sb.WriteString(fmt.Sprintf("%s:\n  %s\n\n", category, stats.String()))
	}
	result := sb.String()
	return result
}
