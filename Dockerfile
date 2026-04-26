FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/vship-api ./cmd/api/main.go

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=build /out/vship-api ./vship-api
COPY web ./web
COPY migrations ./migrations

EXPOSE 3002

CMD ["./vship-api"]
