PROJECT_NAME=imega/commerceml2teleport
GO_PROJECT=github.com/$(PROJECT_NAME)
CWD=/go/src/$(GO_PROJECT)
TAG=latest
IMG=imega/commerceml2teleport

GO_IMG=golang:1.10-alpine
GOLANG_VERSION="1.10"
GRPCURL_COMMIT="f203c2cddfe24b21f8343d989c86db68bf3872aa"

build: unit
	docker build --build-arg GRPCURL_COMMIT=$(GRPCURL_COMMIT) --build-arg GOLANG_VERSION=$(GOLANG_VERSION) -t $(IMG):$(TAG) .

.PHONY: acceptance
acceptance:
	@TAG=$(TAG) docker-compose up -d
	@docker run --rm \
		--network commerceml2teleport_default \
		-v $(CURDIR):$(CWD) \
		$(GO_IMG) sh -c "go test -v $(GO_PROJECT)/acceptance"

clean:
	@TAG=$(TAG) docker-compose rm -sfv

error:
	@docker ps --filter 'status=exited' -q | xargs docker logs

unit:
	@docker run --rm -v $(CURDIR):$(CWD) -w $(CWD) $(GO_IMG) \
		sh -c "go list ./... | grep -v 'vendor\|acceptance' | xargs go test"

pretest:
	@docker run --rm -v $(CURDIR):$(CWD) -w $(CWD) imega/gometalinter \
		$(LINTER_FLAGS) --vendor --deadline=600s --disable=gotype --disable=gocyclo --disable=gas \
		--exclude=/usr --exclude='api' ./...

test:
	echo units

release: build
	@docker login --username $(DOCKER_USER) --password $(DOCKER_PASS)
	@docker push $(IMG):$(TAG)

deploy:
	@curl -s -X POST -H "TOKEN: $(DEPLOY_TOKEN)" https://d.imega.ru -d '{"namespace":"imega-teleport", "project_name":"commerceml2teleport", "tag":"$(TAG)"}'
