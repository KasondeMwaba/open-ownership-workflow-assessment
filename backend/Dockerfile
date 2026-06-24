FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /workflow-api ./cmd/api

FROM alpine:3.22
RUN adduser -D -H appuser
USER appuser
COPY --from=build /workflow-api /workflow-api
EXPOSE 8080
ENTRYPOINT ["/workflow-api"]
