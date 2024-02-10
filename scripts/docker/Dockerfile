FROM golang:alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /app/bin/ ./cmd/main.go

FROM node:alpine AS tailwind

WORKDIR /app

COPY ./static/css/main.css ./static/css/
COPY ./tailwind.config.js ./
COPY ./views ./views
RUN npx tailwindcss -i ./static/css/main.css -o ./static/css/tailwind.css

FROM scratch AS run

WORKDIR /app

COPY --from=build /app/bin ./
COPY --from=build /app/static ./static
COPY --from=build /app/views ./views
COPY --from=tailwind /app/static/css/tailwind.css ./static/css

ENTRYPOINT ["./main"]