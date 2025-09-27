FROM debian:latest


RUN mkdir /horse_network
RUN mkdir /var/log/horse_network

WORKDIR /horse_network
COPY main .

RUN apt update
RUN apt-get update -y && apt-get install ca-certificates -y
RUN chmod +x ./main

ENTRYPOINT ["./main"]
