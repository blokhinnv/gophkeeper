SET GOARCH=amd64

SET GOOS=linux
go build -o ../bin/server/gophkeeper-server-linux-amd64 ../cmd/server

SET GOOS=darwin
go build -o ../bin/server/gophkeeper-server-darwin-amd64 ../cmd/server

SET GOOS=windows
go build -o ../bin/server/gophkeeper-server-windows-amd64.exe ../cmd/server
