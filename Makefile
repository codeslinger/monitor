# vim:set ts=8 ai:

AGENT_SRCS 	= $(wildcard agent/*.go util/*.go)
COLLECTOR_SRCS  = $(wildcard collector/*.go util/*.go)

BIN 	 = agent collector
BUILDDIR = build

all: $(BIN)

$(BUILDDIR)/%:
	mkdir -p $(dir $@)
	cd $* && go build -o $(abspath $@)

$(BIN): %: $(BUILDDIR)/%

$(BUILDDIR)/agent: $(AGENT_SRCS)
$(BUILDDIR)/collector: $(COLLECTOR_SRCS)

clean:
	rm -rf $(BUILDDIR)

.PHONY: all clean $(BIN)

