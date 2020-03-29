GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=bin/main
DATA_FOLDER=data/
OUTPUT_FOLDER=output
OUTPUT_FILE=$(OUTPUT_FOLDER)/final.csv
INDEX_FLAG=--index $(OUTPUT_FILE)

all: build test

build:
	$(GOBUILD) -o $(BINARY_NAME) main.go

test:
	mkdir -p $(OUTPUT_FOLDER)
	$(BINARY_NAME) build --sources $(DATA_FOLDER) $(INDEX_FLAG)
	$(BINARY_NAME) search $(INDEX_FLAG) --search-word hello,world
