# syntax=docker/dockerfile:1

FROM golang:1.17

WORKDIR /app

#COPY go.mod ./
#COPY go.sum ./

COPY . .

RUN go mod download

# RUN go build /app/cmd/main.go -o /default

CMD [ "go", "run", "/app/cmd/main.go", "-c"]
