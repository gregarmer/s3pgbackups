build: clean
	go get "github.com/dustin/go-humanize"
	go get "github.com/goamz/goamz/aws"
	go get "github.com/goamz/goamz/s3"
	go get "github.com/lib/pq"
	[[ ! -d dist ]] && mkdir dist
	go build -o dist/s3pgbackups src/*.go

clean:
	rm -rf dist
