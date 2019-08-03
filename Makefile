GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINNAME=coinboard
UPX=upx

all: clean build

clean:
	$(GOCLEAN)
	rm -f $(BINNAME)

build:
	$(GOBUILD)

run:
	./$(BINNAME)

prod:
	$(GOBUILD) -ldflags="-s -w"
	$(UPX) $(BINNAME)
