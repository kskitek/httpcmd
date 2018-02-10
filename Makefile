all: deps test compile dbuild

deps:
	dep ensure

test: deps
	go test ./...

verify:
	go fmt ./...
	go vet ./...
	go test ./... -race

compile: test
	go build

clean:
	rm httpcmd
	rm httpcmd-linux

compile-linux:
	env GOOS=linux go build -o httpcmd-linux

run: compile
	./fridge

dbuild: compile-linux
	docker build -t registry.gitlab.com/kskitek/httpcmd .

drun: dbuild
	docker run --rm -p 8080:8080 registry.gitlab.com/kskitek/httpcmd
	# docker run --rm -p 8080:8080 --memory=20m registry.gitlab.com/kskitek/httpcmd

dpush: dbuild
	docker push registry.gitlab.com/kskitek/httpcmd