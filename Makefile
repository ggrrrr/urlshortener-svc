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

go_run:
	docker compose up -d postgres
	DB_HOST=localhost \
		DB_USERNAME=test \
		DB_PASSWORD=test \
		DB_DATABASE=test \
		CORS_HOSTS="*" \
		go run be/shorturl/cmd/main.go

login:
	curl -q  -XPOST -d'{"username":"admin","password":"mypass"}' http://localhost:8080/login/v1

url_add:
	curl -q -H'Authorization: Bearer admin@secret' -XPOST -d'{"long_url":"http://yahoo.com"}' http://localhost:8080/admin/v1

url_list:
	curl -s -H'Authorization: Bearer admin@secret' -XGET http://localhost:8080/admin/v1 | jq '.'

url_delete:
	curl -q -H'Authorization: Bearer admin@secret' -XDELETE -d'{"key":"YtBX97ZQEVG"}' http://localhost:8080/admin/v1

url_forward:
	curl -q  -XGET http://localhost:8080/AOkGvLRAujA

