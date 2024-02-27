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
RUN npm install -D tailwindcss @tailwindcss/typography
RUN npx tailwindcss -i ./static/css/main.css -o ./static/css/tailwind.css

FROM centos:latest AS run

WORKDIR /app

COPY --from=build /app/bin ./
COPY --from=build /app/static/assets ./static/assets
COPY --from=build /app/static/js ./static/js
COPY --from=build /app/static/css ./static/css
COPY --from=build /app/views ./views
COPY --from=tailwind /app/static/css/tailwind.css ./static/css

EXPOSE 3000

ENTRYPOINT ["./main"]
