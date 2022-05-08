build-all:
	go build $(FLAGS) -o /dev/null ./cmd

gen:
	go generate dependencies/db/repositories/achievement.go
	go generate dependencies/db/repositories/game.go

test:
	go test -failfast -timeout 10m ./... -v

race:
	go test -race -failfast -timeout 10m ./...

coverage:
	bash ./tools/coverage.sh

coverage-visual: coverage
	go tool cover -html=./tools/coverage.out

lintci:
	golangci-lint run --config ./.golangci.yml

lintci-deps:
	rm -f $(GOPATH)/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.32.2
	go get honnef.co/go/tools/cmd/staticcheck@2020.1.6

run:
	go run cmd/main.go

tidy-check:
	bash ./ci/tidy.sh

tidy:
	go mod tidy

up:
	docker network create --driver=bridge --subnet=172.25.0.0/24 chainnet || true
	docker-compose up

deploy:
	sh ./truffle/run.sh


clean-ansible:
	rm -rf ansible/tmp
	rm -rf ansible/callbacks
	rm -rf ansible/inventory
	rm -rf ansible/terraform/.terraform
	rm ansible/terraform/terraform.tfstate

up-remote:
	cd ansible && ansible-playbook site.yml -vv