FROM debian:latest


RUN mkdir /horse_network
WORKDIR /horse_network
COPY . .

RUN mkdir /var/log/horse_network
RUN apt update
RUN apt-get update -y && apt-get install ca-certificates -y
RUN apt install golang-go -y
RUN go build cmd/main.go
RUN rm -r api/  cmd/  config/  docker/  Dockerfile  go.mod  go.sum  horse_network.code-workspace  LICENSE  README.md  tests/ .test.env .env .git/ .gitignore .github/

ENTRYPOINT ["./main"]
