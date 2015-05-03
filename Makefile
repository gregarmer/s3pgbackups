build: clean
	go get .
	go build -o s3pgbackups *.go

clean:
	rm -f s3pgbackups
