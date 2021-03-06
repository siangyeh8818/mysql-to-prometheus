FROM golang:1.14.4-stretch as builder

COPY go.mod /go/src/github.com/siangyeh8818/mysql-to-prometheus/go.mod
COPY go.sum /go/src/github.com/siangyeh8818/mysql-to-prometheus/go.sum

# Run golang at any directory, not neccessary $GOROOT, $GOPATH
ENV GO111MODULE=on
WORKDIR /go/src/github.com/siangyeh8818/mysql-to-prometheus

# RUN go mod init github.com/pnetwork/sre.monitor.metrics
RUN go mod download
COPY cmd /go/src/github.com/siangyeh8818/mysql-to-prometheus/cmd
COPY internal /go/src/github.com/siangyeh8818/mysql-to-prometheus/internal
#COPY pkg /go/src/github.com/pnetwork/sre.monitor.metrics/pkg

# Build the Go app
RUN env GOOS=linux GOARCH=amd64 go build -o mysql-to-prometheus -v -ldflags "-s" github.com/siangyeh8818/mysql-to-prometheus/cmd/

##### To reduce the final image size, start a new stage with alpine from scratch #####

FROM alpine:3.9
RUN apk --no-cache add ca-certificates libc6-compat busybox-extras

# Run as root
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/siangyeh8818/mysql-to-prometheus/mysql-to-prometheus /usr/local/bin/mysql-to-prometheus

# EXPOSE 8081

ENTRYPOINT [ "mysql-to-prometheus" ] 
