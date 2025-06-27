# go-microservice-1

<!--

TODO:
- Update the title
- Update the summary below about the API you are
  implementing.
-->

The goal of this repository is to inspire a getting
started golang microservice which use a design first
approach and provide a productive way to create
resources for the API.

## Getting started

- **Pre-requisites** on Fedora 41

```sh
# Install required packages
$ sudo dnf install git golang make podman podman-compose delve

# Create and edit the file confis/config.yaml
$ cp -vf configs/config.example.yaml configs/config.yaml

# Install local tools
$ make tidy
$ make install-tools

# Build all the requirements and start the local infra
$ make compose-build compose-clean clean build compose-up
```

Day to day:

- Start the service: `make run`
- Run some test request by: `./test/sripts/todos-list.sh`
- Run tests by: `make test` or `make test-unit` or `make test-integration`

- Open apicurio to load, edit and save the openapi specification
  by: `make apicurio-start`

If you want to use TDD, then start implementing some tests to fail, and
cover every small change. Run tests by `test`, `test-unit` or `test-integration`
make rules and add the necessary code.

## Repository layout

```raw
/api - Contain all the APIs
  /http - Contain the openapi specification.
  /event - Contain json schema to define the event formats.
/internal
  /data - the database model
  /handler - the http and event handlers
    /http
    /event
  /infrastructure
    All the code directly related with
    frameworks and libraries.
  /interface
    The interfaces for presenter, interactor and
    repositories are stored here.
    /presenter
      Adapter between the framework and the
      business logic independent code; every
      http handler and middleware has a
      presenter that wrap the interactor.
      /echo
      /gin
      ...
        **/<resource-name>
    /repository
      Data input/output of the system.
      /client
        /http (this can be self-generated from the openapi.specification).
        /event (producer could be possible to generate code based on json
        schema).
      /db
      /s3
      /cache (redis)
    /interactor
      Free form layout to represent the business
      logic; this is framework independent.
  /usecase
    Same layout that interface, with the implementation.
  /test
    Test helper, mocks, integration tests.
    /mock
    /integration - Hold all the integration tests.
    /builder
      /db
      /api
/test
  /data
    /http
    /event
  /scripts
  /http
```

TODO

- Add workload tests by using
  [locust](https://docs.locust.io/en/stable/what-is-locust.html).
- Add end2end tests by using
  [playwright](https://playwright.dev/).
- Refactor better separability.
- Add a redis sample of caching data from an external system.
- Add a S3 sample for uploading, downloading, and browsing the bucket
  primitives.

## Acknowledgements

- [Content Sources team](https://github.com/content-services/content-sources-backend).
- [insighst-idm team](https://github.com/podengo-project/idmsvc-backend).

