# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app


COPY *.go ./

COPY go.mod ./
RUN go get 
RUN go mod download

RUN go build -o /function

EXPOSE 80

CMD [ "/function" ]