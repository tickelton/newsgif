
SRCDIR = .
BUILDDIR = build

BIN = $(BUILDDIR)/newsgif

all: $(BIN)

$(BUILDDIR):
	@[ -d $(BUILDDIR) ] || mkdir $(BUILDDIR)

$(BIN): $(BUILDDIR) FORCE
	go build -o $(BIN)

run:
	./$(BIN)
clean:
	@rm -rf $(BUILDDIR)

FORCE:
