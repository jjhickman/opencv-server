FROM gocv/opencv:latest
RUN mkdir -p /go/src/opencv-server
COPY . /go/src/opencv-server
