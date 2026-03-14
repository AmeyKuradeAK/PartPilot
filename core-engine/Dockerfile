FROM golang:1.25-alpine AS build

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/core-engine .

FROM gcr.io/distroless/static-debian12

WORKDIR /
COPY --from=build /out/core-engine /core-engine

USER nonroot:nonroot

ENTRYPOINT ["/core-engine"]
