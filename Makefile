# vim:set ts=8 ai:
BUILDDIR = ./build
TARGETS  = agent collector
PACKAGES = util agent/linux agent collector

GOTESTOPTS =

all: compile test

compile: $(TARGETS)

$(TARGETS):
	@mkdir -p $(BUILDDIR) || true
	go build -o $(BUILDDIR)/$@ ./$@

test:
	@for pkg in $(PACKAGES); do \
		go test $(GOTESTOPTS) ./$$pkg; \
	done

clean:
	rm -rf $(BUILDDIR)

over: clean all

.PHONY: all compile test clean over
.PHONY: $(TARGETS)

