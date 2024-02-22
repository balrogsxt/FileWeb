NAME=FileWeb
VERSION=0.0.1
build-all:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/${NAME}_${VERSION}_windows_amd64.exe
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ./bin/${NAME}_${VERSION}_linux_arm
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${NAME}_${VERSION}_linux_amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/${NAME}_${VERSION}_darwin_amd64
	go build -o ./bin/${NAME}_${VERSION}_darwin_arm64
