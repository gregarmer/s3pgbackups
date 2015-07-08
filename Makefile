DIST_DIR = dist

build: clean
	go get "github.com/dustin/go-humanize"
	go get "github.com/goamz/goamz/aws"
	go get "github.com/goamz/goamz/s3"
	go get "github.com/lib/pq"
	test -d $(DIST_DIR) || mkdir $(DIST_DIR)
	go build -o $(DIST_DIR)/s3pgbackups main.go

clean:
	rm -rf $(DIST_DIR)
