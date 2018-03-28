FROM golang:1.10 as build
COPY breakit.go /breakit.go
WORKDIR /
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go get github.com/labstack/echo
RUN go build -ldflags "-s -w" /breakit.go

FROM scratch
EXPOSE 1323
COPY --from=build /breakit /breakit
CMD ["/breakit"]
