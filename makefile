projdir = $(shell pwd)
projname = $(notdir $(projdir))
manfile = readme.md
settingsdir = ~/.local/share/$(projname)

installdir = $(GOPATH)/bin
builddir = $(projdir)/bin

buildpath = $(builddir)/$(projname)
installpath = $(installdir)/$(projname)

ifeq ($(shell uname),Linux)
	os=Linux
	mandir=/usr/local/man/man1
all: install update_manpages
endif
ifeq ($(shell uname),Darwin)
	os=Darwin
	mandir=/usr/local/share/man/man1
all: install 
endif
ifeq ($(OS),Windows_NT)
	os=Windows_NT
all: install_bin
endif

gruffarchive = $(mandir)/$(projname)
gruffpath = $(gruffarchive).1

PANDOC = $(shell which pandoc)

.PHONY: tidy update fmt run build_bin build_docs install_bin install_docs install docs test_doc remove

install: install_bin install_docs 
remove: remove_install remove_docs

tidy:
	@go mod tidy

update:
	git fetch
	git merge
	$(MAKE) -s tidy

upgrade: update
	$(MAKE) -s default

fmt:
	@gofumpt -w .

run: tidy fmt
	@go run .

build_bin: tidy fmt
	@rm -f $(buildpath)
	@mkdir -p $(builddir)
	@go build -o $(buildpath)
	@chmod +x $(buildpath)
	@echo "Successfully built \"$(projname)\" into $(buildpath)"

install_bin: build_bin
	@rm -f $(installpath)
	@mv $(buildpath) $(installpath)
	@echo "Successfully installed \"$(projname)\" at \"$(installpath)\""
	@rm -rf $(builddir)

	

build_docs: $(PANDOC)
	@pandoc $(projdir)/$(manfile) -s -t man -o $(projdir)/$(projname).1

install_docs: $(mandir) build_docs
	@mkdir -p $(mandir)
	sudo mv $(projname).1 $(gruffpath)
	@echo "Successfully installed man-page at $(gruffpath)"
	@echo ""

test_doc: $(PANDOC) build_docs
	man $(projname).1
	@rm $(projname).1
	@echo "Successfully tested docs!"


remove_install:
	sudo rm -rf $(builddir)
	@echo "Erased build directory from \"$(builddir)\""

	sudo rm -f $(installpath)
	@echo "Erased installation from \"$(installpath)\""

	sudo rm -rf $(settingsdir)
	@echo "Erased settings from \"$(settingsdir)\""

remove_docs: $(mandir)
	sudo rm -f $(gruffpath)
	@echo "Erased man-page gruff file from \"$(gruffpath)\""

	sudo rm -f $(gruffarchive)
	@echo "Erased man-page archive from \"$(gruffarchive)\""

update_manpages: $(mandir) $(LINUX)
	echo $(LINUX)
	echo $(mandir)
	sudo mandb
	@echo "Updated man-pages"
