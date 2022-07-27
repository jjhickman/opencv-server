package info

import (
	"net/http"

	"github.com/jjhickman/telescope/internal/log"
	"gocv.io/x/gocv"
)

type Info struct {
	logger     *log.Logger
	apiVersion string
	buildTime  string
}

func (i *Info) Version(w http.ResponseWriter, r *http.Request) {
	i.logger.Info("/version", log.String("buildTime", i.buildTime), log.String("opencv", gocv.OpenCVVersion()), log.String("telescope", i.apiVersion))
}

func New(logger *log.Logger, version string, buildTime string) *Info {
	return &Info{logger: logger, apiVersion: version, buildTime: buildTime}
}
