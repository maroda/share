FROM ghcr.io/ghostgroup/docker-hub/golang:alpine3.20
LABEL app="verificat"
LABEL version="0.0.1"
LABEL org.opencontainers.image.source https://github.com/GhostGroup/verificat
EXPOSE 4330
WORKDIR /go/src/verificat/
COPY . .
RUN go get .
RUN go build -o /bin/verificat
ENTRYPOINT ["/bin/verificat"]
