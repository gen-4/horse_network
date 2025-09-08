FROM debian:latest


RUN mkdir /horse_network
RUN mkdir /var/log/horse_network

WORKDIR /horse_network
COPY cmd api config go.mod go.sum ./

RUN ls -la /horse_network
RUN ls -la /horse_network/cmd || echo "cmd directory missing"
RUN find /horse_network -name "main.go" || echo "main.go not found"

RUN apt update
RUN apt-get update -y && apt-get install ca-certificates -y
RUN apt install golang-go -y
RUN go build ./cmd/main.go
RUN rm -rf api  cmd  config go.mod  go.sum

ENTRYPOINT ["./main"]
