
GOCMD = go
GOINSTALL = $(GOCMD) install
CLIPATH = github.com/binaryedge/40fy-client

default: all

clean:
	rm -rf pkg/* bin/*

all:
	$(GOINSTALL) $(CLIPATH)
