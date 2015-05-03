build: clean
	go get "github.com/dustin/go-humanize"
	go get "github.com/goamz/goamz/aws"
	go get "github.com/goamz/goamz/s3"
	go get "github.com/lib/pq"
	go build -o s3pgbackups *.go

clean:
	rm -f s3pgbackups
