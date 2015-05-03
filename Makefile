build: clean
	go install
	go build -o s3pgbackups *.go

clean:
	rm -f s3pgbackups
