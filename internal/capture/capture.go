package capture

import (
	"context"
	"fmt"
	"sync"

	"github.com/jjhickman/telescope/internal/log"
	"github.com/mattn/go-mjpeg"
	"gocv.io/x/gocv"
)

type Capture struct {
	logger      *log.Logger
	deviceId    int
	stream      *mjpeg.Stream
	ctx         context.Context
	wg          *sync.WaitGroup
	videoHeight int
	videoWidth  int
}

func (c *Capture) Stream() {
	webcam, err := gocv.OpenVideoCapture(c.deviceId)
	webcam.Set(gocv.VideoCaptureFrameWidth, c.videoWidth)
	webcam.Set(gocv.VideoCaptureFrameHeight, c.videoHeight)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to allocate camera %d", c.deviceId)
		c.logger.Error("telescope/capture", log.String("description", errorMessage))
		return
	}
	img := gocv.NewMat()
	for len(c.ctx.Done()) == 0 {
		var buf []byte
		if c.stream.NWatch() > 0 {
			if ok := webcam.Read(&img); !ok {
				continue
			}
			nbuf, err := gocv.IMEncode(".jpg", img)
			if err != nil {
				continue
			}
			buf = nbuf.GetBytes()
			defer nbuf.Close()
		}
		err = c.stream.Update(buf)
		if err != nil {
			break
		}
	}
	defer img.Close()
	defer webcam.Close()
}

func New(logger *log.Logger, deviceId int, stream *mjpeg.Stream, ctx context.Context, wg *sync.WaitGroup, height int, width int) *Capture {
	return &Capture{logger: logger, deviceId: deviceId, stream: stream, ctx: ctx, wg: wg, videoHeight: height, videoWidth: width}
}
