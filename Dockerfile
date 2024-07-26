FROM golang:1.22.5 AS builder
WORKDIR /app
COPY . /app
RUN go build -ldflags "-X 'main.BuildVersion=$(cat VERSION)'\
      -X 'main.BuildDate=$(date +'%Y-%m-%d %H:%M')'\
      -X 'main.BuildCommit=$(git log --oneline -1|awk '{print $1}')'"\
      -tags netgo,osusergo\
      -o app cmd/app/main.go

FROM alpine:3.20
COPY --from=builder /app/app /
EXPOSE 8080
ENTRYPOINT ["/app"]