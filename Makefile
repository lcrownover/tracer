.PHONY: all install clean

all:
	@go build -o bin/tracer cmd/tracer/main.go

install:
	@cp bin/tracer /usr/local/bin/tracer

clean:
	@rm -f bin/tracer /usr/local/bin/tracer

deploy: 
	@docker build --platform=linux/amd64 -t lcrowncrpublic.azurecr.io/tracer .
	@docker push lcrowncrpublic.azurecr.io/tracer
	@az container create --resource-group rg-prometheus-mock --name tracer --image lcrowncrpublic.azurecr.io/tracer --ports 80 --dns-name-label tracer
	@az container restart --resource-group rg-prometheus-mock --name tracer
