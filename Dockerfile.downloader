FROM golang:1.24

WORKDIR /downloader

COPY go.mod .
COPY go.sum .

RUN go mod download

# Add git and ssh client for CronJob pull/push operations
RUN apt-get update && apt-get install -y --no-install-recommends \
    git openssh-client ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY . .

RUN go build -o downloader .

CMD ./downloader
