## We specify the base image we need for our
## go application
FROM golang:1.16.5-buster
## We create an /app directory within our
## image that will hold our application source
## files

# ENV GOPATH=''

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build

RUN /app/textlooker-backend migrate

EXPOSE 8080
ENTRYPOINT ["/app/textlooker-backend", "run"]