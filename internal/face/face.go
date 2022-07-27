package face

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jjhickman/telescope/internal/log"
	"gocv.io/x/gocv"
)

type Face struct {
	logger     *log.Logger
	classifier gocv.CascadeClassifier
}

var upgrader = websocket.Upgrader{} // use default options

func (f *Face) Detect(w http.ResponseWriter, r *http.Request) {
	f.logger.Info("telescope/face/detect", log.String("source", "..."))
}

func (f *Face) Blur(w http.ResponseWriter, r *http.Request) {
	f.logger.Info("/face/blur", log.String("source", "..."))

}

func (f *Face) BlurVideo(w http.ResponseWriter, r *http.Request) {
	f.logger.Info("telescope/face/blur/video", log.String("source", "..."))
	/*
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	*/
}

func (f *Face) DetectVideo(w http.ResponseWriter, r *http.Request) {
	f.logger.Info("telescope/face/detect/video", log.String("source", "..."))
	/*
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	*/
}

func (f *Face) Destroy() {
	defer f.classifier.Close()
}

func New(logger *log.Logger, xmlFile string) *Face {
	classifier := gocv.NewCascadeClassifier()
	if !classifier.Load(xmlFile) {
		return nil
	}
	return &Face{logger: logger, classifier: classifier}
}
