set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
set OLDGOPATH=%GOPATH%
set GOPATH=%CD%;%GOPATH%
go install vc-export
set GOPATH=%OLDGOPATH%
