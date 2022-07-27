package main

import (
	"context"
	"flag"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/jjhickman/telescope/internal/capture"
	"github.com/jjhickman/telescope/internal/face"
	"github.com/jjhickman/telescope/internal/info"
	"github.com/jjhickman/telescope/internal/log"
	"github.com/mattn/go-mjpeg"
)

var (
	version   string
	buildTime string
)

func home(w http.ResponseWriter, r *http.Request) {

}

func main() {
	var address = flag.String("address", "0.0.0.0:8080", "http service address")
	var logPath = flag.String("logPath", "./data", "Directory to store logs")
	var logSize = flag.Int("logSize", 25, "Log file size in MB")
	var logAge = flag.Int("logAge", 1, "Log file age in days")
	var logBackups = flag.Int("logBackups", 5, "Number of log files backed up")
	var cascadeXmlFile = flag.String("cascadeXmlFile", "./data/haarcascade_frontalface_default.xml", "File path for Haar Cascade XML template")
	var enableCapture = flag.Bool("capture", false, "Enable capturing from backend camera if available")
	var deviceId = flag.Int("camera", 0, "Device ID of backend camera")
	var interval = flag.Duration("interval", 50*time.Millisecond, "interval")
	var videoHeight = flag.Int("videoHeight", 720, "Height of video capture in pixels")
	var videoWidth = flag.Int("videoWidth", 720, "Width of video capture in pixels")
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
	log.Info("telescope", log.String("version", version), log.String("buildTime", buildTime), log.String("logPath", *logPath), log.Bool("capture", *enableCapture))
	i := info.New(logger, version, buildTime)
	f := face.New(logger, *cascadeXmlFile)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	stream := mjpeg.NewStreamWithInterval(*interval)
	var c *capture.Capture
	if *enableCapture {
		c = capture.New(logger, *deviceId, stream, ctx, &wg, *videoHeight, *videoWidth)
		go c.Stream()
		http.HandleFunc("/capture", stream.ServeHTTP)
	}

	if f != nil {
		http.HandleFunc("/face/detect", f.Detect)
		http.HandleFunc("/face/detect/video", f.DetectVideo)
		http.HandleFunc("/face/blur", f.Blur)
		http.HandleFunc("/face/blur/video", f.BlurVideo)
	}

	http.HandleFunc("/version", i.Version)
	http.HandleFunc("/", home)
	err := http.ListenAndServe(*address, nil).Error()
	if f != nil {
		f.Destroy()
	}
	if stream != nil {
		stream.Close()
	}
	cancel()
	wg.Wait()
	log.Fatal("telescope", log.String("error", err))
}
