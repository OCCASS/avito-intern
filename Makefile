APP_OUTPUT=main
APP_CMD=./cmd/app
LOAD_TEST_CMD=./cmd/load_test
CONFIG_FILE=./config/local.yml

build:
	go build -o $(APP_OUTPUT) $(APP_CMD)

run:
	go run $(APP_CMD) -c $(CONFIG_FILE)

prod_run:
	$(APP_OUTPUT) -c $(CONFIG_FILE)

lint:
	golangci-lint run

load_test:
	go run $(LOAD_TEST_CMD) -c $(CONFIG_FILE)

