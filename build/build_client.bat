SET BUILD_DATE=%DATE% %TIME%
SET GOARCH=amd64

SET GOOS=linux
go build -o ../bin/client/gophkeeper-client-linux-amd64 -ldflags "-X 'main.buildVersion=v0.0.1' -X 'main.buildDate=%BUILD_DATE%'" ../cmd/client

SET GOOS=darwin
go build -o ../bin/client/gophkeeper-client-darwin-amd64 -ldflags "-X 'main.buildVersion=v0.0.1' -X 'main.buildDate=%BUILD_DATE%'" ../cmd/client

SET GOOS=windows
go build -o ../bin/client/gophkeeper-client-windows-amd64.exe -ldflags "-X 'main.buildVersion=v0.0.1' -X 'main.buildDate=%BUILD_DATE%'" ../cmd/client
