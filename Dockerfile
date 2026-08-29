# golden-field 应用容器（下游组合根；多阶段——依赖解析走公共 Go proxy）。
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/golden-field ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=build /app/golden-field /usr/local/bin/golden-field
ENTRYPOINT ["/usr/local/bin/golden-field"]