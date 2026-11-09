FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY src/ ./src/
RUN CGO_ENABLED=0 go build -o /orders-api ./src

FROM gcr.io/distroless/static:nonroot
COPY --from=build /orders-api /orders-api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/orders-api"]
