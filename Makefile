DIST_DIR = dist

deps:
	go get "github.com/dustin/go-humanize"
	go get "github.com/goamz/goamz/aws"
	go get "github.com/goamz/goamz/s3"
	go get "github.com/lib/pq"

build: clean deps
	test -d $(DIST_DIR) || mkdir $(DIST_DIR)
	go build -o $(DIST_DIR)/s3pgbackups main.go

clean:
	rm -rf $(DIST_DIR)

test:
	@go test -run=. -test.v ./config
	@go test -run=. -test.v ./database
	@go test -run=. -test.v ./dest
	@go test -run=. -test.v ./utils
