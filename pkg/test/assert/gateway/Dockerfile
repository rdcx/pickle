# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY *.go ./

COPY go.mod ./
RUN go get 
RUN go mod download

RUN go build -o /function

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o function

FROM scratch

COPY --from=0 /app/function .
CMD ["./function"]