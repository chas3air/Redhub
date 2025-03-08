ARG GO_VERSION=latest
FROM golang:${GO_VERSION} AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=${TARGETARCH} go build -o /src/cli ./cmd/app 

FROM alpine:latest AS final

RUN apk --no-cache add ca-certificates tzdata

RUN mkdir /app
COPY --from=build /src /app

RUN mv /app/cli /cli

RUN adduser -D appuser
USER appuser

WORKDIR /

EXPOSE 50051

# ENTRYPOINT ["tail", "-f", "/dev/null"]
CMD ["./cli", "--config=/app/config/local.yaml"]