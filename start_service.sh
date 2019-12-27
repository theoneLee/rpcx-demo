go build -o serviceBin service/service.go
./serviceBin $1 $2
#usage ：
# start the first node :  sh start_service.sh -addr localhost:8973
# start the second : sh start_service.sh -addr localhost:8974