FROM golang:latest
RUN set -e; \
    git clone https://github.com/hybridgroup/gocv.git; \
    cd gocv; \
    cat Makefile | sed -e 's/\<sudo\> //g' > temp.txt; \
    mv temp.txt Makefile; \
    make install_static BUILD_SHARED_LIBS=OFF; \
    make clean
WORKDIR /go/src
COPY . /go/src/opencv-server
RUN cd /go/src/opencv-server; \
    go build -o build/opencv-server -tags static /go/src/opencv-server/cmd/opencv-server

FROM debian:bullseye-slim
WORKDIR /opencv-server
RUN set -e; \
    apt-get update; \
    apt-get install --no-install-recommends -y libjpeg-dev; \
    apt-get autoremove -y; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*
COPY --from=0 /go/src/opencv-server/build/opencv-server ./
EXPOSE 8080
CMD ["./opencv-server"]
