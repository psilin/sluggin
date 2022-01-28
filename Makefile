downloader:
	go build -o ./cmd/downloader/downloader ./cmd/downloader/main.go

server:
	go build -o ./cmd/server/server ./cmd/server/main.go

clean:
	rm ./cmd/downloader/downloader ./cmd/server/server
