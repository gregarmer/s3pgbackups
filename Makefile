build: clean
	go install "github.com/dustin/go-humanize"
	go install "github.com/goamz/goamz/aws"
	go install "github.com/goamz/goamz/s3"
	go install "github.com/lib/pq"
	go build -o s3pgbackups *.go

clean:
	rm -f s3pgbackups
