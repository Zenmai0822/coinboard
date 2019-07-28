GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINNAME=cryptogo

all: clean build

clean:
	$(GOCLEAN)
	rm -f $(BINNAME)

build:
	$(GOBUILD)

run:
	./$(BINNAME)
