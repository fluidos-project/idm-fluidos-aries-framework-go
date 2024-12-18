# Makefile for managing the model treatment service

# Define the base directory
BASE_DIR := $(shell realpath $(dir $(lastword $(MAKEFILE_LIST))))

# Define the paths to the scripts and docker-compose file
UPDATE_CERTS_SCRIPT := $(BASE_DIR)/updateCerts.sh
DOCKER_COMPOSE_FILE := $(BASE_DIR)/docker-compose.yml

# Run the model treatment service
.PHONY: run-model-treatment
run-model-treatment:
	@echo "Running updateCerts.sh..."
	@bash $(UPDATE_CERTS_SCRIPT)
	@echo "Building and starting Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Stop the model treatment service
.PHONY: stop-model-treatment
stop-model-treatment:
	@echo "Stopping Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Clean up the environment
.PHONY: clean
clean: stop-model-treatment
	@echo "Cleaning up Docker resources..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans
