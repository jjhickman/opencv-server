package main

import (
	"flag"
	"net/http"
	"path/filepath"

	"github.com/jjhickman/telescope/internal/capture"
	"github.com/jjhickman/telescope/internal/face"
	"github.com/jjhickman/telescope/internal/info"
	"github.com/jjhickman/telescope/internal/log"
)

const version = "v0.1"

func home(w http.ResponseWriter, r *http.Request) {
}

func main() {
	var addr = flag.String("address", "localhost:8080", "http service address")
	var logPath = flag.String("logPath", "./data", "Directory to store logs")
	var logSize = flag.Int("logSize", 25, "Log file size in MB")
	var logAge = flag.Int("logAge", 1, "Log file age in days")
	var logBackups = flag.Int("logBackups", 5, "Number of log files backed up")
	var cascadeXmlFile = flag.String("cascadeXmlFile", "./data/haarcascade_frontalface_default.xml", "File path for Haar Cascade XML template")
	var enableCapture = flag.Bool("capture", false, "Enable capturing from backend camera if available")
	var deviceId = flag.Int("camera", 0, "Device ID of backend camera")
	flag.Parse()

	var tops = []log.TeeOption{
		{
			Filename: filepath.Join(*logPath, "telescope.log"),
			Ropt: log.RotateOptions{
				MaxSize:    *logSize,
				MaxAge:     *logAge,
				MaxBackups: *logBackups,
				Compress:   true,
			},
			Lef: func(lvl log.Level) bool {
				return lvl <= log.InfoLevel
			},
		},
	}

	logger := log.NewTeeWithRotate(tops)
	log.ResetDefault(logger)
	log.Info("telescope", log.String("status", "started"))

	i := info.New(logger, version)
	f := face.New(logger, *cascadeXmlFile)
	var c *capture.Capture
	if *enableCapture {
		c = capture.New(logger, *cascadeXmlFile, *deviceId)
		http.HandleFunc("/capture", c.Raw)
	}

	if f != nil {
		http.HandleFunc("/face/detect", f.Detect)
		http.HandleFunc("/face/detect/video", f.DetectVideo)
		http.HandleFunc("/face/blur", f.Blur)
		http.HandleFunc("/face/blur/video", f.BlurVideo)
	}

	http.HandleFunc("/version", i.Version)
	http.HandleFunc("/", home)
	err := http.ListenAndServe(*addr, nil).Error()
	if f != nil {
		f.Destroy()
	}
	if c != nil {
		c.Destroy()
	}
	log.Fatal("telescope", log.String("error", err))
}
