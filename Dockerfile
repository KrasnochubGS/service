FROM golang:1.20.3

WORKDIR /service

COPY go.mod go.sum ./ 

RUN go mod download

COPY . .

RUN go build -o /docker-gs-ping

EXPOSE 5300

CMD ["/docker-gs-ping"]