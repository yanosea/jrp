# shows help message defaultly
.DEFAULT_GOAL := help

#
# update
#
.PHONY: update.credits update.mocks update.swagger

# update `./CREDITS`
update.credits:
	gocredits -skip-missing . > ./CREDITS

# update mocks
update.mocks:
	# ./app/application
	mockgen -source=./app/application/wnjpn/word_query_service.go -destination=./app/application/wnjpn/word_query_service_mock.go -package=wnjpn
	# ./app/infrastructure
	mockgen -source=./app/infrastructure/database/connection.go -destination=./app/infrastructure/database/connection_mock.go -package=database
	mockgen -source=./app/infrastructure/database/connection_manager.go -destination=./app/infrastructure/database/connection_manager_mock.go -package=database
	# ./app/domain
	mockgen -source=./app/domain/jrp/history/history_repository.go -destination=./app/domain/jrp/history/history_repository_mock.go -package=history
	# ./app/presentation/api/jrp-server/server
	mockgen -source=./app/presentation/api/jrp-server/server/server.go -destination=./app/presentation/api/jrp-server/server/server_mock.go -package=server
	# ./app/presentation/cli/jrp/command
	mockgen -source=./app/presentation/cli/jrp/command/command.go -destination=./app/presentation/cli/jrp/command/command_mock.go -package=command
	# ./app/presentation/cli/jrp/formatter
	mockgen -source=./app/presentation/cli/jrp/formatter/formatter.go -destination=./app/presentation/cli/jrp/formatter/formatter_mock.go -package=formatter
	# ./pkg/proxy
	mockgen -source=./pkg/proxy/buffer.go -destination=./pkg/proxy/buffer_mock.go -package=proxy
	mockgen -source=./pkg/proxy/cobra.go -destination=./pkg/proxy/cobra_mock.go -package=proxy
	mockgen -source=./pkg/proxy/debug.go -destination=./pkg/proxy/debug_mock.go -package=proxy
	mockgen -source=./pkg/proxy/echo.go -destination=./pkg/proxy/echo_mock.go -package=proxy
	mockgen -source=./pkg/proxy/envconfig.go -destination=./pkg/proxy/envconfig_mock.go -package=proxy
	mockgen -source=./pkg/proxy/gzip.go -destination=./pkg/proxy/gzip_mock.go -package=proxy
	mockgen -source=./pkg/proxy/http.go -destination=./pkg/proxy/http_mock.go -package=proxy
	mockgen -source=./pkg/proxy/io.go -destination=./pkg/proxy/io_mock.go -package=proxy
	mockgen -source=./pkg/proxy/json.go -destination=./pkg/proxy/json_mock.go -package=proxy
	mockgen -source=./pkg/proxy/keyboard.go -destination=./pkg/proxy/keyboard_mock.go -package=proxy
	mockgen -source=./pkg/proxy/os.go -destination=./pkg/proxy/os_mock.go -package=proxy
	mockgen -source=./pkg/proxy/pflag.go -destination=./pkg/proxy/pflag_mock.go -package=proxy
	mockgen -source=./pkg/proxy/promptui.go -destination=./pkg/proxy/promptui_mock.go -package=proxy
	mockgen -source=./pkg/proxy/rand.go -destination=./pkg/proxy/rand_mock.go -package=proxy
	mockgen -source=./pkg/proxy/spinner.go -destination=./pkg/proxy/spinner_mock.go -package=proxy
	mockgen -source=./pkg/proxy/sql.go -destination=./pkg/proxy/sql_mock.go -package=proxy
	mockgen -source=./pkg/proxy/tablewriter.go -destination=./pkg/proxy/tablewriter_mock.go -package=proxy
	# ./pkg/utility
	mockgen -source=./pkg/utility/capture.go -destination=./pkg/utility/capture_mock.go -package=utility
	mockgen -source=./pkg/utility/download_util.go -destination=./pkg/utility/download_util_mock.go -package=utility
	mockgen -source=./pkg/utility/file_util.go -destination=./pkg/utility/file_util_mock.go -package=utility
	mockgen -source=./pkg/utility/json_util.go -destination=./pkg/utility/json_util_mock.go -package=utility
	mockgen -source=./pkg/utility/keyboard_util.go -destination=./pkg/utility/keyboard_util_mock.go -package=utility
	mockgen -source=./pkg/utility/prompt_util.go -destination=./pkg/utility/prompt_util_mock.go -package=utility
	mockgen -source=./pkg/utility/rand_util.go -destination=./pkg/utility/rand_util_mock.go -package=utility
	mockgen -source=./pkg/utility/spinner_util.go -destination=./pkg/utility/spinner_util_mock.go -package=utility
	mockgen -source=./pkg/utility/strings_util.go -destination=./pkg/utility/strings_util_mock.go -package=utility
	mockgen -source=./pkg/utility/tablewriter_util.go -destination=./pkg/utility/tablewriter_util_mock.go -package=utility
	mockgen -source=./pkg/utility/version_util.go -destination=./pkg/utility/version_util_mock.go -package=utility

# update swagger files
update.swagger:
	swag init -g ./app/presentation/api/jrp-server/main.go --parseDependency --output ./docs

#
# container
#
.PHONY: container.build container.down

# build container
container.build:
	@set -e; \
	if [ -f "./container.exist" ]; then \
		echo "container already exist"; \
		exit 1; \
	fi; \
	docker-compose -f docker-compose.yml build --no-cache; \
	touch ./container.exist

# down container
container.down:
	@set -e; \
	docker-compose down; \
	docker image prune -af; \
	if [ -f "./container.exist" ]; then \
		rm ./container.exist; \
	fi

#
# test
#
.PHONY: test.local test.container test.container.once

# execute tests in local
test.local:
	@set -e; \
	if [ -f "./test.run" ]; then \
		echo "test already running"; \
		exit 1; \
	fi; \
	touch test.run; \
	go test -v -p 1 ./... -cover -coverprofile=./cover.out; \
	grep -v -E "(_mock\.go|/mock/|/proxy/|/docs/docs\.go)" ./cover.out > ./cover.out.tmp && mv ./cover.out.tmp ./cover.out; \
	go tool cover -html=./cover.out -o ./docs/coverage.html; \
	rm ./cover.out; \
	if [ -f "./test.run" ]; then \
		rm ./test.run; \
	fi

# execute tests in container
test.container:
	@set -e; \
	if ! [ -f "./container.exist" ]; then \
		echo "container not exist"; \
		exit 1; \
	fi; \
	if [ -f "./test.run" ]; then \
		echo "test already running"; \
		exit 1; \
	fi; \
	touch test.run; \
	docker-compose -f docker-compose.yml up --abort-on-container-exit jrp-test-container; \
	CONTAINER_ID=$$(docker ps -a -q --filter "name=jrp-test-container" --filter "status=exited"); \
	docker cp $${CONTAINER_ID}:/jrp/docs/coverage.html ./docs/coverage.html; \
	rm ./test.run

# execute tests in container (once)
test.container.once:
	@set -e; \
	if [ -f "./container.exist" ]; then \
		echo "container already exist"; \
		exit 1; \
	fi; \
	if [ -f "./test.run" ]; then \
		echo "test already running"; \
		exit 1; \
	fi; \
	touch ./container.exist; \
	touch test.run; \
	docker-compose -f docker-compose.yml build --no-cache; \
	docker-compose -f docker-compose.yml up --abort-on-container-exit jrp-test-container; \
	CONTAINER_ID=$$(docker ps -a -q --filter "name=jrp-test-container" --filter "status=exited"); \
	docker cp $${CONTAINER_ID}:/jrp/docs/coverage.html ./docs/coverage.html; \
	docker-compose down; \
	docker image prune -af; \
	rm ./test.run; \
	rm ./container.exist

#
# misc
#
.PHONY: misc.run-server

# run api server
misc.run-server:
	go run ./app/presentation/api/jrp-server/main.go


# required phony targets for standards
all: help
clean:
	@rm -f ./cover.out ./co
	@rm -f ./test.run ./con
	@docker-compose down
	@docker image prune -af
test: test.local

# help
.PHONY: help
help:
	@echo ""
	@echo "available targets:"
	@echo ""
	@echo "  [update]"
	@echo "    update.credits       - update ./CREDITS file"
	@echo "    update.mocks         - update all mocks"
	@echo "    update.swagger       - update docs/swagger files"
	@echo ""
	@echo "  [container]"
	@echo "    container.build      - build container for testing"
	@echo "    container.down       - down container and remove images"
	@echo ""
	@echo "  [test]"
	@echo "    test.local           - execute all tests in local"
	@echo "    test.container       - execute all tests in container"
	@echo "    test.container.once  - build container and execute all tests in container once, then remove container and images"
	@echo ""
	@echo "  [miscellaneous]"
	@echo "    misc.run-server      - run api server"
	@echo ""
