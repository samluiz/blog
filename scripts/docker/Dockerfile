FROM golang:alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /app/bin/ ./cmd/main.go

FROM node:alpine AS tailwind

WORKDIR /app

COPY ./static/css/main.css ./static/css/
RUN npx tailwindcss -i ./static/css/main.css -o ./static/css/tailwind.css

FROM scratch AS run

WORKDIR /app

COPY --from=build /app/bin ./
COPY --from=tailwind /app/static/css/tailwind.css ./static/css

CMD ["./main"]