package capture

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jjhickman/telescope/internal/log"
	"gocv.io/x/gocv"
)

type Capture struct {
	logger     *log.Logger
	deviceId   int
	classifier gocv.CascadeClassifier
}

var upgrader = websocket.Upgrader{} // use default options

func (c *Capture) Raw(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("/capture", log.Int("deviceId", c.deviceId))
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

func (c *Capture) Destroy() {
	defer c.classifier.Close()
}

func New(logger *log.Logger, xmlFile string, deviceId int) *Capture {
	classifier := gocv.NewCascadeClassifier()
	if !classifier.Load(xmlFile) {
		return nil
	}
	return &Capture{logger: logger, deviceId: deviceId, classifier: classifier}
}
