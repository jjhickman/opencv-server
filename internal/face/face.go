package face

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"
)

var upgrader = websocket.Upgrader{} // use default options

func Detect(w http.ResponseWriter, r *http.Request) {
	log.Print(gocv.OpenCVVersion())

}

func Blur(w http.ResponseWriter, r *http.Request) {
	log.Print(gocv.OpenCVVersion())

}

func BlurVideo(w http.ResponseWriter, r *http.Request) {
	log.Print(gocv.OpenCVVersion())
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func DetectVideo(w http.ResponseWriter, r *http.Request) {
	log.Print(gocv.OpenCVVersion())
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
