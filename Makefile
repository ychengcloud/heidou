apps = 'heidou'

all: lint cover

.PHONY: test
test: 
	for app in $(apps) ;\
	do \
		go test ./... -covermode=count -coverprofile=dist/cover-$$app.out ;\
	done

.PHONY: build
build: 
	for app in $(apps) ;\
	do \
		go build -o dist/$$app ./cmd/$$app/; \
	done

.PHONY: install
install:
	for app in $(apps) ;\
	do \
		go install ./cmd/$$app/; \
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
	go vet ./... | grep -v assets/ && exit 1 || exit 0
