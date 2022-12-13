FROM golang:1.19-alpine AS development

WORKDIR /ryanpujo/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /build ./cmd

EXPOSE 8080

CMD [ "/build" ]




FROM alpine:latest AS production

WORKDIR /app

COPY userApp /

CMD ["/userApp"]