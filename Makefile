build: clean
	go build -o s3pgbackups *.go

clean:
	rm -f s3pgbackups
