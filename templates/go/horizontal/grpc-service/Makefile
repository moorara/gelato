# Include macros, variables, and rules
include ./make/common.mk
include ./make/go.mk      # test, test-short, test-coverage, clean-test, run, build, build-all, clean-build
include ./make/grpc.mk    # check-tools, protoc, protoc-gen-go, protobuf
include ./make/docker.mk  # docker, docker-test, push, push-latest, save-docker, load-docker, clean-docker

# Variables required by inclusions
name := grpc-service
proto_path := idl
go_out_path := internal/idl
docker_image := dockerid/grpc-service
docker_tag ?= $(version)
