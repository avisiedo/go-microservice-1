##
# This build the container image for the service
##


# https://catalog.redhat.com/software/containers/ubi9/go-toolset/61e5c00b4ec9945c18787690
FROM registry.access.redhat.com/ubi9/go-toolset:1.22.5-1731639025@sha256:c7bfd2501cb1be171366434a368db669b32f08a0198c1473b9bff0a379613fc3 as builder
LABEL todo-backend=builder
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
WORKDIR /go/src/app
COPY . .
USER 0
RUN make get-deps build


# https://catalog.redhat.com/software/containers/ubi9-minimal/61832888c0d15aff4912fe0d
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4-949.1714662671@sha256:2636170dc55a0931d013014a72ae26c0c2521d4b61a28354b3e2e5369fa335a3
LABEL todo-backend=backend
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
RUN mkdir -p /opt/bin /opt/bin/scripts/db /opt/bin/configs
WORKDIR /opt/bin
COPY --from=builder /go/src/app/bin/* ./
COPY scripts/db/migrations /opt/bin/scripts/db/migrations
USER 1001
# TODO Accomodate to the service arguments
CMD ["/opt/bin/service"]
