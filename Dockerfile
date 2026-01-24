FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

# -ldflags="-s -w"
RUN go build -o /app/main  ./cmd

FROM alpine:latest

RUN apk add --no-cache gettext

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY config/*.yaml config/config.yaml
COPY entrypoint.sh entrypoint.sh

RUN chmod u+x entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]

CMD [ "/app/main" ]