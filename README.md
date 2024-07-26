# msg-proc
Micro-service for process messages.

## Build
```
go build -ldflags    "-X 'main.BuildVersion=$(cat VERSION)'\
   -X 'main.BuildDate=$(date +'%Y-%m-%d %H:%M')'\
   -X 'main.BuildCommit=$(git log --oneline -1|awk '{print $1}')'"\
   -tags netgo,osusergo\
   -o app cmd/app/main.go
```

## Run in docker
```
cd infra
docker-compose up -d
```