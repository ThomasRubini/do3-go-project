FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd cmd
COPY pkg pkg

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    CGO_ENABLED=1 go build -o /binary cmd/nutritionapp/main.go

CMD ["/binary"]
