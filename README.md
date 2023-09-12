# telemetry

## consul-demo

```shell
docker pull prom/prometheus
docker run -d --restart=always -p 9090:9090 -v /Users/momo/Desktop/consul-demo/prometheus.yml:/etc/prometheus/prometheus.yml --name prometheus prom/prometheus:latest

docker pull hashicorp/consul 
docker run -d --restart=always --name consul -d -p 8500:8500 hashicorp/consul 

go mod tidy
go run main.go
```

