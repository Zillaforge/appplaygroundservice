OWNER ?= Zillaforge
PROJECT ?= AppPlaygroundService
ABBR ?= aps
IMAGE_NAME ?= app-playground-service
WORK_DIR ?= /home/appplaygroundservice
GOVERSION ?= 1.22.4
OS ?= ubuntu
ARCH ?= $(shell uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')
PREVERSION ?= 0.0.7
VERSION ?= $(shell cat VERSION)
LOCAL_PWD := $(shell pwd)
HOST_PWD := $(or $(HOST_PROJECT_PATH),$(LOCAL_PWD))
# GO_PROXY ?= "https://proxy.golang.org,http://proxy.pegasus-cloud.com:8078"

# Release Mode could be dev or prod,
# dev: default, will add commit id to version
# prod: will use version only
RELEASE_MODE ?= dev
COMMIT_ID ?= $(shell git rev-parse --short=8 HEAD)

sed = sed
ifeq ("$(shell uname -s)", "Darwin")	# BSD sed, like MacOS
	sed += -i ''
else	# GNU sed, like LinuxOS
	sed += -i''
endif

ifeq ($(RELEASE_MODE),prod)
    RELEASE_VERSION := $(VERSION)
else
    RELEASE_VERSION := $(VERSION)-$(COMMIT_ID)
endif

.PHONY: go-build
go-build:
	@echo "Build Binary"
	@go build -ldflags="-s -w" -o tmp/$(PROJECT)_$(VERSION)

.PHONY: build
build: go-build
ifeq ($(OS), ubuntu)
	@sh build/build-debian.sh
else
	@sh build/build-rpm.sh
endif

.PHONY: set-version
set-version:
	@echo "Set Version: $(RELEASE_VERSION)"
	@$(sed) -e'/$(PREVERSION)/{s//$(RELEASE_VERSION)/;:b' -e'n;bb' -e\} $(LOCAL_PWD)/build/$(PROJECT).spec
	@$(sed) -e'/$(PREVERSION)/{s//$(RELEASE_VERSION)/;:b' -e'n;bb' -e\} $(LOCAL_PWD)/constants/common.go
	@$(sed) -e'/$(PREVERSION)/{s//$(RELEASE_VERSION)/;:b' -e'n;bb' -e\} $(LOCAL_PWD)/etc/app-playground-service.yaml
	@$(sed) -e'/$(PREVERSION)/{s//$(RELEASE_VERSION)/;:b' -e'n;bb' -e\} $(LOCAL_PWD)/etc/aps-scheduler.yaml
	@$(sed) -e'/$(PREVERSION)/{s//$(RELEASE_VERSION)/;:b' -e'n;bb' -e\} $(LOCAL_PWD)/Makefile

.PHONY: build-container
build-container:
	@echo "Build Container"
	@rm -rf build/scratch_image/tmp/*
	@go build -o build/scratch_image/tmp/$(PROJECT)
	@sh build/scratch_image/build_scratch_img_env.sh

.PHONY: release
release:
	@make set-version
	@mkdir -p tmp
	@rm -rf tmp/$(OS)
	@docker rm -f build-env
	@docker run --name build-env --rm -v $(HOST_PWD):$(WORK_DIR) -w $(WORK_DIR) $(OWNER)/golang:$(GOVERSION)-$(OS)-$(ARCH) make OS=$(OS) build
	@mkdir tmp/$(OS)
	@mv tmp/$(PROJECT)* tmp/$(OS)

.PHONY: release-image
release-image:
	@make set-version
	@mkdir -p build/scratch_image/tmp
	@docker rm -f build-env
	@docker run --name build-env --rm -v $(HOST_PWD):$(WORK_DIR) -w $(WORK_DIR) $(OWNER)/golang:$(GOVERSION)-$(OS)-$(ARCH) make build-container
	@docker rmi -f $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION)
	@docker build -t $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION) build/scratch_image/
	@docker run --name build-env --rm -v $(HOST_PWD):$(WORK_DIR) -w $(WORK_DIR) $(OWNER)/golang:$(GOVERSION)-$(OS)-$(ARCH) rm -rf build/scratch_image/tmp/*

.PHONY: release-image-file
release-image-file: release-image
	@rm -rf tmp/container
	@mkdir -p tmp/container
	@docker save $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION) > tmp/container/$(IMAGE_NAME)_$(RELEASE_VERSION).image.tar

.PHONY: push-image
push-image:
	@echo "Check Image $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION)"
	@docker image inspect $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION) --format="image existed"
	@echo "Push Image"
	@docker logout
	@echo "<DOCKER HUB KEY>" | docker login -u $(OWNER) --password-stdin
	@docker image push $(OWNER)/$(IMAGE_NAME):$(RELEASE_VERSION)
	@docker logout

.PHONY: start
start:
	@go run main.go -c etc/app-playground-service.yaml serve

.PHONY: start-scheduler
start-scheduler:
	@go run main.go -c etc/app-playground-service.yaml -s etc/aps-scheduler.yaml scheduler start

.PHONY: init
init:
	@go run main.go -c etc/app-playground-service.yaml database sync

.PHONY: start-dev-env
start-dev-env:
	@make start-dev-persistent
	@make start-dev-system
	@make start-dev-service
	
.PHONY: start-dev-service
start-dev-service: docker-compose/service/docker-compose.*.yaml
	@for f in $^; do ARCH=$(ARCH) COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-service" up -d --no-recreate || true ; done

.PHONY: start-dev-system
start-dev-system: docker-compose/system/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-system" up -d --no-recreate || true ; done

.PHONY: start-dev-persistent
start-dev-persistent: docker-compose/persistent/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-system" up -d --no-recreate --no-start || true ; done

.PHONY: stop-dev-env # Stop and Remove current service only
stop-dev-env:
	COMPOSE_IGNORE_ORPHANS=True docker-compose -f docker-compose/service/docker-compose.${ABBR}.yaml -p "pegasus-service" down
	
.PHONY: stop-dev-all # Stop and Remove all dependency
stop-dev-all:
	@make stop-dev-service
	@make stop-dev-system

.PHONY: purge-dev-all # Stop and Remove all dependency include persistent network and volume
purge-dev-all:
	@make stop-dev-all
	@make clean-dev-persistent

.PHONY: stop-dev-service
stop-dev-service: docker-compose/service/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-service" down -v; done

.PHONY: stop-dev-system
stop-dev-system: docker-compose/system/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-system" down -v; done

.PHONY: clean-dev-persistent
clean-dev-persistent: docker-compose/persistent/docker-compose.*.yaml
	@for f in $^; do COMPOSE_IGNORE_ORPHANS=True docker-compose -f $${f} -p "pegasus-system" down -v; done
