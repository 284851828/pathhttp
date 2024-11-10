
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64

go build -o ph.exe   main.go 

SET GOOS=linux
SET GOARCH=amd64

go build -o ph  main.go 


