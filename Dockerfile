##
# BUILD CONTAINER
##

FROM golang:1.10 as builder

WORKDIR /go/src/github.com/mvisonneau/evm

COPY Makefile .
RUN \
make setup

COPY . .
RUN \
make deps ;\
make build-docker

##
# RELEASE CONTAINER
##

FROM scratch

WORKDIR /

COPY --from=builder /go/src/github.com/mvisonneau/evm/evm /

ENTRYPOINT ["/evm"]
CMD [""]
