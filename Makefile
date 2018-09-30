NAME = cmd

BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
COMMIT = $(shell git rev-parse --short HEAD)
BUILDTIME = $(shell date +%Y-%m-%dT%T%z)
GOPATH = $(shell echo "$$GOPATH")
LD_OPTS = -ldflags="-X main.branch=${BRANCH} -X main.commit=${COMMIT} -X main.lasttag=${LASTTAG} -X main.buildtime=${BUILDTIME} -w "

all:  build run

all-with-deps: 
	cd cmd && GOOS=linux GOARCH=amd64 go build $(LD_OPTS)  -o $(NAME) .

run:
	cd ./cmd && rm -rf $(NAME) && cd .. && make build  && cd cmd && ./$(NAME) && ../
	
dist:
	cd ./cmd/ && GOOS=linux GOARCH=amd64 go build $(LD_OPTS)  -o $(NAME) .

build:
	 cd ./cmd/ && go build $(LD_OPTS) -o $(NAME) . && cd -

# Show to-do items per file.
todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow|nolint:.*' .
.PHONY: todo

