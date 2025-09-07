FROM debian:latest


RUN mkdir /horse_network
WORKDIR /horse_network
COPY cmd/ .
COPY api/ .
COPY config/ .
COPY go.mod .
COPY go.sum .

RUN mkdir /var/log/horse_network
RUN apt update
RUN apt-get update -y && apt-get install ca-certificates -y
RUN apt install golang-go -y
RUN go build cmd/main.go
RUN rm -rf api/  cmd/  config/ go.mod  go.sum

ENTRYPOINT ["./main"]
