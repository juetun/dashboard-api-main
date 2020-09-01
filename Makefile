run:
	export GOPROXY=https://mirrors.aliyun.com/goproxy/
	go get ./...
	export GOPROXY=https://proxy.golang.org,direct
	go build -o  juetun_dashboard_api_main main.go
	./juetun_dashboard_api_main
dev:
	go run main.go