FROM golang:1.15.2-alpine3.12 as builder

# All these steps will be cached
WORKDIR /app
COPY go.mod  go.sum ./

# Get dependancies - will also be cached if mod/sum are not changed
RUN go mod download

# COPY the source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app

## SECOND STEP: Build minimal image by copying the executable binary file from builder
FROM alpine:3.12.0

COPY --from=builder /go/bin/app /go/bin/app

WORKDIR /go/bin/

EXPOSE 8080

ENTRYPOINT ["/go/bin/app"]