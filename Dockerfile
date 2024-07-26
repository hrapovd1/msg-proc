ARG COMMIT
ARG DATE
FROM golang:1.22.5 AS builder
WORKDIR /app
COPY . /app
RUN cd /app && go build -ldflags "-X 'main.BuildVersion=$(cat VERSION)'\
   -X 'main.BuildDate=${DATE}'\
   -X 'main.BuildCommit=${COMMIT}'"\
   -tags netgo,osusergo\
   -o app cmd/app/main.go

FROM alpine:3.20
COPY --from=builder /app/app /
EXPOSE 8080
ENTRYPOINT ["/app"]