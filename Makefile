APP_NAME = std-deviation-calculator
BUILD_DIR =.build

.PHONY: all verify test build run build-docker run-docker

all: fmt test build run

verify:
	go vet ./...

test: verify
	go test ./...

fmt: verify
	go fmt ./...

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) .

run:
	$(BUILD_DIR)/$(APP_NAME)

run-mock:
	cd mockRandomOrg && go build -o ../$(BUILD_DIR)/mockRandom .
	$(BUILD_DIR)/mockRandom

build-docker:
	docker build -t $(APP_NAME) .

run-docker:
	docker run --rm --net=host $(APP_NAME)
