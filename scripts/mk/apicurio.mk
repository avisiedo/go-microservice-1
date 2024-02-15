
# https://hub.docker.com/search?q=apicurio%2Fapicurito&type=image
.PHONY: apicurio-start
apicurio-start: ## Start and open apicurio web app
	$(CONTAINER_ENGINE) container exists apicurio || \
	$(CONTAINER_ENGINE) run --rm -d --name apicurio --publish 8090:8080 docker.io/apicurio/apicurito-ui:latest
	$(OPEN) http://localhost:8090/html

.PHONY: apicurio-stop
apicurio-stop: ## Stop apicurio web app
	$(CONTAINER_ENGINE) stop apicurio
	$(CONTAINER_ENGINE) container wait apicurio
