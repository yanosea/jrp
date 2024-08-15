# test packages and output a report file
go test -v -p 1 ../... -cover -coverprofile=./cover.out
awk '!/(mock)\//' ./cover.out > temp_cover.out && mv temp_cover.out cover.out
go tool cover -html=./cover.out -o ./cover.html
