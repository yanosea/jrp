# test all packages and output a report file
go test -v ../... -cover -coverprofile=./cover.out
go tool cover -html=./cover.out -o ./cover.html
