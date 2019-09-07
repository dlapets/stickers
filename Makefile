build:
	go build -o stickers cmd/stickers/main.go

install:
	go install ./cmd/stickers/

clean:
	rm -f *.log
	rm -f stickers

test:
	go test ./... -v
	golint ./... | grep -v ^vendor && exit 1 || exit 0

.PHONY:
	test clean
