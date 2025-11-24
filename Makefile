APP_OUTPUT=main
APP_CMD=./cmd/app
TESTS_DIR=./tests
LOAD_TEST_CMD=./cmd/load_test
CONFIG_FILE=./config/local.yml
TEST_CONFIG_FILE=../config/test.yml
DOCKER_ENV_FILE=.env.docker

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

test:
	docker-compose -f docker-compose.test.yml --env-file $(DOCKER_ENV_FILE) up -d test-db

	@until docker-compose -f docker-compose.test.yml exec test-db pg_isready -U testuser; do \
		sleep 1; \
	done

	go test $(TESTS_DIR) -config $(TEST_CONFIG_FILE)

	docker-compose -f docker-compose.test.yml stop test-db
	docker-compose -f docker-compose.test.yml rm -f test-db
