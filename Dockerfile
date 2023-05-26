FROM golang:1.20.4

WORKDIR /service

COPY go.mod go.sum ./ 

RUN go mod download

COPY . .

RUN go build -o /docker-gs-ping

# EXPOSE 50051

# CMD ["/docker-gs-ping"]
CMD ["/docker-gs-ping"]