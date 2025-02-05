DOCKER_REPO ?= "local"
GIT_HASH ?= $(shell git log --format="%h" -n 1)

go_test:
	docker compose up -d postgres
	go test -cover -race ./be/...

go_lint:
	golangci-lint run -v ./be/...

go_clean:
	go clean -testcache
	# go clean -fuzzcache
	# go clean -modcache

go_build:
	docker build \
		-f ./be.Dockerfile \
		--tag ${DOCKER_REPO}/be/shorturl:${GIT_HASH} \
		./
	docker tag ${DOCKER_REPO}/be/shorturl:${GIT_HASH} ${DOCKER_REPO}/be/shorturl:latest

build_svc:
	docker build \
		-f ./.docker/be.Dockerfile \
		--tag "${DOCKER_REPO}/svc/rssaggregator:${GIT_HASH}" \
		./
	
	docker tag ${DOCKER_REPO}/svc/rssaggregator:${GIT_HASH} ${DOCKER_REPO}/svc/rssaggregator:latest

url_add:
	curl -q -H'Authorization: Some admin@secret' -XPOST -d'{"long_url":"http://yahoo.com"}' http://localhost:8080/admin/v1

url_list:
	curl -s -H'Authorization: Some admin@secret' -XGET http://localhost:8080/admin/v1 | jq '.'

url_delete:
	curl -q -H'Authorization: Some admin@secret' -XDELETE -d'{"key":"CjmlyvZwdaJ"}' http://localhost:8080/admin/v1
