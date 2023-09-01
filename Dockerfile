############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

WORKDIR /chep
COPY . .

RUN go mod download
RUN apk --update add tzdata

# Build the binary
RUN date +%Y%m%d%T > /build-timestamp
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -a -o brambles . 

############################
# STEP 2 build a small image
############################
FROM golang:alpine AS cryptobot

WORKDIR /chep
COPY . .

ENV USER=appuser
ENV UID=10001

RUN apk --update add ca-certificates
RUN mkdir -p /etc/ssl/certs/
RUN update-ca-certificates

RUN echo -e "http://nl.alpinelinux.org/alpine/v3.5/main\nhttp://nl.alpinelinux.org/alpine/v3.5/community" > /etc/apk/repositories

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
	
# Import from builder.
COPY --from=builder /build-timestamp /build-timestamp

# Copy our static executable
COPY --from=builder /chep/brambles .

# Use an unprivileged user.
USER appuser:appuser

# Run the hello binary.
EXPOSE 8081 8081
ENTRYPOINT ["./brambles"]