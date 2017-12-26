.PHONY: all clean

all: d c w

d:
	go install github.com/u35s/shell/shelld
c:
	go install github.com/u35s/shell/shellc
w:
	go install github.com/u35s/shell/webc

test:

bench:

clean:
	rm -f $(GOPATH)/bin/shelld $(GOPATH)/bin/shellc $(GOPATH)/bin/webc 
