# vim:set ts=8 ai:
AGENT_SRCS = $(wildcard agent/*.go util/*.go)
COLLECTOR_SRCS = $(wildcard collector/*.go util/*.go)

TARGETS  = agent collector
BUILDDIR = build

all: $(TARGETS)

$(BUILDDIR)/%:
	@mkdir -p $(dir $@)
	cd $* && go build -o $(abspath $@)

$(TARGETS): %: $(BUILDDIR)/%

$(BUILDDIR)/agent: $(AGENT_SRCS)
$(BUILDDIR)/collector: $(COLLECTOR_SRCS)

test:
	@for dir in $(TARGETS); do \
	  cd $$dir; \
	  go test; \
	  cd ..; \
	done

deps:
	@for dir in $(TARGETS); do \
	  cd $$dir; \
	  go get -v; \
	  cd ..; \
	done

clean:
	rm -rf $(BUILDDIR)

over: clean all

.PHONY: all test deps clean over
.PHONY: $(TARGETS)

