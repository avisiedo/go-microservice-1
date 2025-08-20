##
# Set default variable values for the project
##
APP_NAME ?= todo
export APP_NAME

BIN ?= bin
PATH := $(CURDIR)/$(BIN):$(PATH)
export PATH

CONFIG_PATH ?= $(PROJECT_DIR)/configs
export CONFIG_PATH
CONFIG_YAML := $(CONFIG_PATH)/config.yaml

COMPOSE_FILE ?= $(PROJECT_DIR)/deploy/docker-compose.yaml

CONTAINER_IMAGE_BASE ?= quay.io/$(firstword $(subst +, ,$(QUAY_USER)))/$(APP_NAME)-$(APP_COMPONENT)

# Tools and their dependencies
# Build dependencies
TOOLS_BIN := tools/bin

GODA := $(TOOLS_BIN)/goda
# GOJSONSCHEMA not in TOOLS actually
GOJSONSCHEMA := $(TOOLS_BIN)/gojsonschema
GOLANGCI_LINT := $(TOOLS_BIN)/golangci-lint
MOCKERY := $(TOOLS_BIN)/mockery
OAPI_CODEGEN := $(TOOLS_BIN)/oapi-codegen
PLANTER := $(TOOLS_BIN)/planter
YQ := $(TOOLS_BIN)/yq


TOOLS := \
	$(GODA) \
	$(GOLANGCI_LINT) \
	$(MOCKERY) \
	$(OAPI_CODEGEN) \
	$(PLANTER) \
	$(YQ)

TOOLS_DEPS := tools/go.mod tools/go.sum tools/tools.go | $(TOOLS_BIN)

#
# Kafka configuration variables
#

LOAD_DB_CFG_WITH_YQ := n
ifneq (,$(shell "$(YQ)" --version 2>/dev/null))
ifneq (,$(shell ls -1 "$(CONFIG_YAML)" 2>/dev/null))
LOAD_DB_CFG_WITH_YQ := y
endif
endif

DATABASE_CONTAINER_NAME="database"
ifeq (y,$(LOAD_DB_CFG_WITH_YQ))
$(info info:Trying to load DATABASE configuration from '$(CONFIG_YAML)')
DATABASE_HOST ?= $(shell "$(YQ)" -r -M '.database.host' "$(CONFIG_YAML)")
DATABASE_EXTERNAL_PORT ?= $(shell "$(YQ)" -M '.database.port' "$(CONFIG_YAML)")
DATABASE_NAME ?= $(shell "$(YQ)" -r -M '.database.name' "$(CONFIG_YAML)")
DATABASE_USER ?= $(shell "$(YQ)" -r -M '.database.user' "$(CONFIG_YAML)")
DATABASE_PASSWORD ?= $(shell "$(YQ)" -r -M '.database.password' "$(CONFIG_YAML)")
else
$(info info:Using DATABASE_* defaults)
DATABASE_HOST ?= localhost
DATABASE_EXTERNAL_PORT ?= 5432
DATABASE_NAME ?= database-db
DATABASE_USER ?= database-user
DATABASE_PASSWORD ?= database-secret
endif

#
# Minio Service (S3 storage)
#
MINIO_USER ?= minio
MINIO_PASSWORD ?= minio

#
# Kafka configuration variables
#

# The directory where the kafka data will be stored
KAFKA_DATA_DIR ?= $(PROJECT_DIR)/build/kafka/data

# The directory where the kafka configuration will be
# bound to the containers
KAFKA_CONFIG_DIR ?= $(PROJECT_DIR)/build/kafka/config

# The topics used by the repository
# Updated to follow the pattern used at playbook-dispatcher
KAFKA_TOPICS ?= platform.todo.todo-created

# The group id for the consumers; every consumer subscribed to
# a topic with different group-id will receive a copy of the
# message. In our scenario, any replica of the consumer wants
# only one message to be processed, so we only use a unique
# group id at the moment.
KAFKA_GROUP_ID ?= todos

# Application specific parameters
APP_EXPIRATION_TIME ?= 15
export APP_EXPIRATION_TIME
APP_PAGINATION_DEFAULT_LIMIT ?= 10
export APP_PAGINATION_DEFAULT_LIMIT
APP_PAGINATION_MAX_LIMIT ?= 100
export APP_PAGINATION_MAX_LIMIT
# Enable IS_FAKE_ENABLED for the ephemeral deployment
APP_ACCEPT_X_RH_FAKE_IDENTITY ?= true
export APP_ACCEPT_X_RH_FAKE_IDENTITY
APP_VALIDATE_API ?= true
export APP_VALIDATE_API

# Set the default token expiration in seconds (2 hours)
APP_TOKEN_EXPIRATION_SECONDS ?= 7200
export APP_TOKEN_EXPIRATION_SECONDS

# main secret for various MAC and encryptions like
# domain registration token and encrypted private JWKs
APP_SECRET ?= sFamo2ER65JN7wxZ48UZb5GbtDc053ahIPJ0Qx47bzA
export APP_SECRET
