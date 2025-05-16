compile:
	echo "Compiling for Linux OS"
	GOOS=linux GOARCH=386 go build
	
test:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./...

#docker relate commands
#push:
   