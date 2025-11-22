APP_OUTPUT=main
APP_CMD=./cmd/app
CONFIG_FILE=./config/local.yml

build:
	go build -o $(APP_OUTPUT) $(APP_CMD)

run:
	go run $(APP_CMD) -c $(CONFIG_FILE)

prod_run:
	$(APP_OUTPUT) -c $(CONFIG_FILE)
