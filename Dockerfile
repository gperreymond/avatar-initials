# stage 0
FROM golang:latest as builder
RUN mkdir -p /go/src/github.com/gperreymond/avatar-initials/app
WORKDIR /go/src/github.com/gperreymond/avatar-initials/app
COPY . .
RUN go get github.com/manyminds/api2go && \
    go get github.com/fogleman/gg
RUN GOARCH=amd64 GOOS=linux go build -ldflags "-linkmode external -extldflags -static -w"

# stage 1
FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/gperreymond/avatar-initials/app .
ENTRYPOINT ["/app"]
