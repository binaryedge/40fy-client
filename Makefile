
GOCMD = go
GOINSTALL = $(GOCMD) install
CLIPATH = github.com/binaryedge/be-cli

default: all

clean:
	rm -rf pkg/* bin/*

all:
	$(GOINSTALL) $(CLIPATH)
