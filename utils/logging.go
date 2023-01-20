package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var httpTracer *FileTracer
var applicationTracer *FileTracer
var once sync.Once

const LogPath string = "logs"

type FileTracer struct {
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	warningLogger *log.Logger
	logFile       *os.File
}

type Tracer interface {
	LogInfo(message ...any)
	LogError(message ...any)
	LogWarning(message ...any)
	Close()
}

func CloseLoggers() {
	if httpTracer != nil {
		httpTracer.Close()
	}
	if applicationTracer != nil {
		httpTracer.Close()
	}
}

func initLoggers() {
	once.Do(func() {
		httpTracer = &FileTracer{}
		httpTracer.initLogger("htaccess.log")

		applicationTracer = &FileTracer{}
		applicationTracer.initLogger("store.log")
	})
}

func createDirIfNotExists() {
	if _, err := os.Stat(LogPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(LogPath, os.ModePerm)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to create log directory %s Error: %v", LogPath, err))
		}
	}
}
func HttpTracer() *FileTracer {
	if httpTracer == nil {
		initLoggers()
	}

	return httpTracer
}

func ApplicationTracer() *FileTracer {
	if applicationTracer == nil {
		initLoggers()
	}

	return applicationTracer
}

func (h *FileTracer) initLogger(logName string) {
	createDirIfNotExists()
	logFile, err := os.OpenFile(filepath.Join(LogPath, logName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(fmt.Printf("Error initialising logger %s", err))
	}

	h.infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	h.errorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	h.warningLogger = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (h *FileTracer) LogInfo(message ...any) {
	h.infoLogger.Println(message...)
	fmt.Println(message...)
}

func (h *FileTracer) LogError(message ...any) {
	h.errorLogger.Println(message...)
	fmt.Println(message...)
}

func (h *FileTracer) LogWarning(message ...any) {
	h.warningLogger.Println(message...)
	fmt.Println(message...)
}

func (h *FileTracer) Close() {
	h.logFile.Close()
}
