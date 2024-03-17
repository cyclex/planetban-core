# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build

# Main target
all: build

# Build the executable
build:
	$(GOBUILD) -o engine cmd/*.go

# Run the application
run:
	$(GOBUILD) -o engine cmd/.go
	mv engine ../
	sudo systemctl restart planetban-server
	sudo systemctl restart planetban-webhook

# Default target to run the application
.DEFAULT_GOAL := run
