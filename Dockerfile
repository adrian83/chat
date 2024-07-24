FROM golang:1.22-bookworm

ADD . /chat
WORKDIR /chat

EXPOSE 7070

RUN go build -o chat cmd/chat/main.go

CMD ["./chat"]