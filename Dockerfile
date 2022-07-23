FROM golang:latest
RUN set -e; \
    git clone https://github.com/hybridgroup/gocv.git; \
    cd gocv; \
    cat Makefile | sed -e 's/\<sudo\> //g' > temp.txt; \
    mv temp.txt Makefile; \
    make install_static BUILD_SHARED_LIBS=OFF; \
    make clean
WORKDIR /go/src
COPY . /go/src/telescope
RUN cd /go/src/telescope; \
    go build -o build/telescope -ldflags="-s -w" -tags static /go/src/telescope/cmd/telescope

FROM debian:bullseye-slim
WORKDIR /telescope
RUN set -e; \
    apt-get update; \
    apt-get install --no-install-recommends -y libjpeg-dev; \
    apt-get autoremove -y; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*
COPY --from=0 /go/src/telescope/build/telescope ./
EXPOSE 8080
CMD ["./telescope"]
