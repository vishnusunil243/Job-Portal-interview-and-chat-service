FROM golang:1.21.5-bullseye AS build

RUN apt-get update

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o interviewchat-service

FROM busybox:latest 

WORKDIR /interviewchat/cmd

COPY --from=build /app/cmd/interviewchat-service .

COPY --from=build /app/.env /interviewchat

EXPOSE 8000

CMD ["./interviewchat-service"]