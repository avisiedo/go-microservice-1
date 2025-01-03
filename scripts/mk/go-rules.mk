##
# Golang rules to build the binaries, tidy dependencies,
# generate vendor directory, download dependencies and clean
# the generated binaries.
##

GOVERSION := 1.22
export GOVERSION

ifeq (,$(shell ls -1d vendor 2>/dev/null))
MOD_VENDOR :=
else
MOD_VENDOR ?= -mod vendor
endif

.PHONY: install-go-tools
install-go-tools: $(TOOLS) ## Install Go tools

# used by ipa-hcc backend test
.PHONY: install-xrhidgen
install-xrhidgen: $(XRHIDGEN) ## Install xrhidgen tool

.PHONY: install-tools
install-tools: install-go-tools install-python-tools ## Install tools used to build, test and lint

.PHONY: build-all
build-all: ## Generate code and build binaries
	$(MAKE) generate-api
	# $(MAKE) generate-event
	$(MAKE) generate-mock
	$(MAKE) generate-diagrams
	$(MAKE) build

# Meta rule to add dependency on the binaries generated
.PHONY: build
build: $(patsubst cmd/%,$(BIN)/%,$(wildcard cmd/*)) ## Build binaries

$(BIN) $(TOOLS_BIN):
	mkdir -p $@

# export CGO_ENABLED
# $(BIN)/%: CGO_ENABLED=0
# Build by path, not by referring to main.go. It's required to bake VCS
# information into binary, see golang/go#51279.
$(BIN)/%: cmd/%/main.go $(BIN)
	go build -C $(dir $<) $(MOD_VENDOR) -buildvcs=true -o "$(CURDIR)/$@" .

$(TOOLS_BIN)/%: $(TOOLS_DEPS)
	cd tools && GOBIN="$(PROJECT_DIR)/tools/bin" go install $(shell grep $(notdir $@) tools/tools.go | awk '{print $$2}')

.PHONY: clean
clean: ## Clean binaries and testbin generated
	@[ ! -e "$(BIN)" ] || for item in cmd/*; do rm -vf "$(BIN)/$${item##cmd/}"; done

.PHONY: cleanall
cleanall: ## Clean and remove all binaries
	rm -rf $(BIN) $(TOOLS_BIN)

.PHONY: run
run: $(BIN)/service .compose-wait-db ## Run the api & kafka consumer locally
	"$(BIN)/service"

# See: https://go.dev/doc/modules/managing-dependencies#synchronizing
.PHONY: tidy
tidy:  ## Synchronize your code's dependencies
	go mod tidy -go=$(GOVERSION)
	cd tools && go mod tidy -go=$(GOVERSION)

.PHONY: get-deps
get-deps: ## Download golang dependencies
	go mod download

.PHONY: update-deps
update-deps: ## Update all golang dependencies
	go get -u -t ./...
	$(MAKE) tidy

.PHONY: vet
vet:  ## Run go vet ignoring /vendor directory
	go vet $(shell go list ./... | grep -v /vendor/)

.PHONY: go-fmt
go-fmt:  ## Run go fmt ignoring /vendor directory
	go fmt $(shell go list ./... | grep -v /vendor/)

.PHONY: vendor
vendor: ## Generate vendor/ directory populated with the dependencies
	go mod vendor

# Exclude /internal/test/mock /internal/api directories because the content is
# generated.
# Exclude /vendor in case it exists
# Exclude /internal/interface directories because only contain interfaces
TGF_PREFIX := github.com/avisiedo/go-microservice-1
TEST_GREP_FILTER := -v \
  -e '^$(TGF_PREFIX)/vendor/' \
  -e '^$(TGF_PREFIX)/internal/test' \
  -e '^$(TGF_PREFIX)/internal/interface/' \
  -e '^$(TGF_PREFIX)/internal/api/metrics' \
  -e '^$(TGF_PREFIX)/internal/api/private' \
  -e '^$(TGF_PREFIX)/internal/usecase/repository/s3' \
  -e '^$(TGF_PREFIX)/internal/usecase/repository/event' \
  -e '^$(TGF_PREFIX)/internal/api/http/public' \
  -e '^$(TGF_PREFIX)/internal/api/http/private' \
  -e '^$(TGF_PREFIX)/internal/api/event' \


.PHONY: test
test: test-unit test-integration  ## Run unit tests and integration tests

.PHONY: test-unit
test-unit: ## Run unit tests
	go test -parallel 4 -coverprofile="coverage.out" -covermode count $(MOD_VENDOR) $(shell go list ./... | grep $(TEST_GREP_FILTER) )

.PHONY: test-ci
test-ci: ## Run tests for ci
	go test $(MOD_VENDOR) ./...

.PHONY: test-integration
test-integration:  ## Run integration tests
	CONFIG_PATH="$(PROJECT_DIR)/configs" go test -parallel 1 ./internal/test/integration/... -test.failfast -test.v

# Add dependencies from binaries to all the the sources
# so any change is detected for the build rule
$(patsubst cmd/%,$(BIN)/%,$(wildcard cmd/*)): $(shell find $(PROJECT_DIR)/cmd -type f -name '*.go') $(shell find $(PROJECT_DIR)/pkg -type f -name '*.go' 2>/dev/null) $(shell find $(PROJECT_DIR)/internal -type f -name '*.go' 2>/dev/null)

# # Regenerate code when message schema changes
# $(shell find "$(EVENT_MESSAGE_DIR)" -type f -name '*.go'): $(SCHEMA_YAML_FILES)
# 	$(MAKE) gen-event-messages


############### TOOLS

# https://github.com/RedHatInsights/playbook-dispatcher/blob/master/Makefile
API_LIST := api/public.openapi.yaml api/internal.openapi.yaml
.PHONY: generate-api
generate-api: $(OAPI_CODEGEN) $(API_LIST) ## Generate server stubs from openapi
	# Public API
	$(OAPI_CODEGEN) -generate spec -package public -o internal/api/http/public/spec.gen.go api/http/public.openapi.yaml
	$(OAPI_CODEGEN) -generate server -package public -o internal/api/http/public/server.gen.go api/http/public.openapi.yaml
	$(OAPI_CODEGEN) -generate types -package public -o internal/api/http/public/types.gen.go -alias-types api/http/public.openapi.yaml
	# Internal API # FIXME Update -import-mapping options
	$(OAPI_CODEGEN) -generate server -package private -o internal/api/http/private/server.gen.go api/http/internal.openapi.yaml
	$(OAPI_CODEGEN) -generate types -package private -o internal/api/http/private/types.gen.go api/http/internal.openapi.yaml

$(API_LIST):
	#git submodule update --init

.PHONY: update-api
update-api:
	#git submodule update --init --remote
	$(MAKE) generate-api

EVENTS := todo_created
# Generate event types
.PHONY: generate-event
generate-event: $(GOJSONSCHEMA) $(SCHEMA_JSON_FILES)  ## Generate event messages from schemas
	@[ -e "$(EVENT_MESSAGE_DIR)" ] || mkdir -p "$(EVENT_MESSAGE_DIR)"
	for event in $(EVENTS); do \
		$(GOJSONSCHEMA) -p event "$(EVENT_SCHEMA_DIR)/$${event}.event.json" -o "$(PROJECT_DIR)/internal/api/event/$${event}.event.types.gen.go"; \
	done

.PHONY: generate-event-debug
generate-event-debug:
	@echo SCHEMA_JSON_FILES=$(SCHEMA_JSON_FILES)
	@echo GOJSONSCHEMA=$(GOJSONSCHEMA)
	@echo EVENT_MESSAGE_DIR=$(EVENT_MESSAGE_DIR)
	@echo EVENT_SCHEMA_DIR=$(EVENT_SCHEMA_DIR)

# Generic rule to generate the JSON files
$(EVENT_SCHEMA_DIR)/%.event.json: $(EVENT_MESSAGE_DIR)/%.event.yaml
	@[ -e "$(EVENT_MESSAGE_DIR)" ] || mkdir -p "$(EVENT_MESSAGE_DIR)"
	yaml2json "$<" "$@"

# Mockery support
MOCK_DIRS := internal/api/http/private \
	internal/api/http/public \
	internal/api/http/openapi \
	internal/api/http/metrics \
	internal/api/http/healthcheck \
	internal/interface/presenter/echo \
	internal/interface/interactor \
	internal/interface/repository/event \
	internal/interface/repository/client \
	internal/interface/repository/db \
	internal/handler/http \
	internal/infrastructure/event \
	internal/infrastructure/service \
	internal/infrastructure/middleware \

.PHONY: generate-mock
generate-mock: $(MOCKERY)  ## Generate mock by using mockery tool
	for item in $(MOCK_DIRS); do \
	  PKG="$${item##*/}"; \
	  DEST_DIR="internal/test/mock/$${item#*/}"; \
	  [ -e "$${DEST_DIR}" ] || mkdir -p "$${DEST_DIR}"; \
	  $(MOCKERY) \
	    --all \
	    --outpkg "$${PKG}" \
	    --dir "$${item}" \
		--output "$${DEST_DIR}" \
		--case underscore || exit 1; \
	done

.PHONY: generate-deps
generate-deps: $(GODA)
	$(GODA) graph "github.com/avisiedo/go-microservice-1/..." | dot -Tsvg -o docs/service-dependencies.svg

coverage.tmp.out:
	go test -parallel 4 -coverprofile="coverage.tmp.out" -covermode count $(MOD_VENDOR) $(shell go list ./... | grep $(TEST_GREP_FILTER) )

coverage.out: coverage.tmp.out
	grep -v -e '.gen.go' coverage.tmp.out > coverage.out

.PHONY: coverage
coverage:  coverage.out  ## Printout coverage
	go tool cover -func ./coverage.out
