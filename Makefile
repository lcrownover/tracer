.PHONY: all install clean

all:
	@go build -o bin/tracer cmd/tracer/main.go

install:
	@cp bin/tracer /usr/local/bin/tracer

container:
	@docker build -t tracer .

clean:
	@rm -f bin/tracer /usr/local/bin/tracer

