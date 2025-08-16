.PHONY: build run

build:
	go build -o bin/kustom-scheduler ./cmd

run: build
	./bin/kustom-scheduler

clean:
	rm -rf bin/kustom-scheduler
