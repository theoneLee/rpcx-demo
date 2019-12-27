go build -o serviceBin service/service.go
./serviceBin $1 $2
#usage ：
# sh start_service.sh -addr localhost:8973