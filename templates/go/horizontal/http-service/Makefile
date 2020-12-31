# Include macros, variables, and rules
include ./make/common.mk
include ./make/go.mk      # test, test-short, test-coverage, clean-test, run, build, build-all, clean-build
include ./make/docker.mk  # docker, docker-test, push, push-latest, save-docker, load-docker, clean-docker

# Variables required by inclusions
name := http-service
docker_image := dockerid/http-service
docker_tag ?= $(version)
