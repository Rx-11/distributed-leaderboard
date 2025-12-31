FROM golang:1.23-alpine AS builder

RUN apk add --no-cache protobuf git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY proto/ proto/
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/*.proto

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /dl ./main.go


FROM gcr.io/distroless/static-debian12

COPY --from=builder /dl /dl

ENV DATA_DIR=/data
ENV HTTP_PORT=8080
ENV GRPC_PORT=9090

EXPOSE 8080 9090

VOLUME ["/data"]

ENTRYPOINT ["/dl"]