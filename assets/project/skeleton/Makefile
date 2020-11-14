apps = 'server'

.PHONY: run
run: 
	for app in $(apps) ;\
	do \
		 go run ./cmd/$$app -c configs/$$app ;\
	done

.PHONY: watch
watch:
	gowatch

.PHONY: wire
wire:
	wire ./cmd/server

.PHONY: test
test: 
	for app in $(apps) ;\
	do \
		go test -v ./app/$$app/... -f `pwd`/configs/$$app.yml -covermode=count -coverprofile=dist/cover-$$app.out ;\
	done

.PHONY: assets_gen
assets_gen:
	for app in $(apps) ;\
	do \
		go generate -tags=dev ./graphql... ;\
	done

.PHONY: build
build: assets_gen wire
	for app in $(apps) ;\
	do \
		GOOS=linux GOARCH="amd64" go build -o dist/$$app-linux-amd64 ./cmd/$$app/; \
		GOOS=darwin GOARCH="amd64" go build -o dist/$$app-darwin-amd64 ./cmd/$$app/; \
	done

.PHONY: cover
cover: test
	for app in $(apps) ;\
	do \
		go tool cover -html=dist/cover-$$app.out; \
	done

.PHONY: mock
mock:
	mockery --all

.PHONY: lint
lint:
	golint ./...

.PHONY: proto
proto:
	protoc -I api/proto ./api/proto/* --go_out=plugins=grpc:api/proto

.PHONY: dash
dash: # create grafana dashboard
	 for app in $(apps) ;\
	 do \
	 	jsonnet -J ./grafana/grafonnet-lib   -o ./configs/grafana/dashboards/$$app.json  --ext-str app=$$app ./scripts/grafana/dashboard.jsonnet ;\
	 done

.PHONY: pubdash
pubdash:
	 for app in $(apps) ;\
	 do \
	 	jsonnet -J ./grafana/grafonnet-lib  -o ./configs/grafana/dashboards-api/$$app-api.json  --ext-str app=$$app  ./scripts/grafana/dashboard-api.jsonnet ; \
	 	curl -X DELETE --user admin:admin  -H "Content-Type: application/json" 'http://localhost:3000/api/dashboards/db/$$app'; \
	 	curl -x POST --user admin:admin  -H "Content-Type: application/json" --data-binary "@./configs/grafana/dashboards-api/$$app-api.json" http://localhost:3000/api/dashboards/db ; \
	 done

.PHONY: rules
rules:
	for app in $(apps) ;\
	do \
	 	jsonnet  -o ./configs/prometheus/rules/$$app.yml --ext-str app=$$app  ./scripts/prometheus/rules.jsonnet ; \
	done

.PHONY: docker
docker-compose: build dash rules
	docker-compose -f deployments/docker-compose.yml up --build -d

all: lint cover docker
