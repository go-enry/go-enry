RESOURCES_SRC=../.shared
RESOURCES_DIR=./shared

SHELL := /bin/bash

all: $(RESOURCES_DIR)

$(RESOURCES_DIR): os-shared-lib
	cp -R $(RESOURCES_SRC) $(RESOURCES_DIR)
	mvn package

os-shared-lib:
	@os_name="$(shell uname -s)" && \
	if [ "$$os_name" == "Linux" ]; then \
		$(MAKE) linux-shared; \
	elif [ "$$os_name" == "Darwin" ]; then \
		$(MAKE) darwin-shared; \
	else \
		echo "Unsupported operating system, can't build shared library"; \
		exit 1; \
	fi;

linux-shared:
	$(MAKE) -C .. linux-shared

darwin-shared:
	$(MAKE) -C .. darwin-shared

test:
	mvn clean test

package:
	mvn package

clean:
	rm -rf $(RESOURCES_DIR)
	mvn clean

.PHONY: test package clean linux-shared darwin-shared os-shared-lib
