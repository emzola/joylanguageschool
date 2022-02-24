# ==================================================================================== #
# HELPERS
# ==================================================================================== #

# Include variables from the .envrc file
include .envrc

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/web: run the cmd/web application
.PHONY: run/web
run/web:
	@echo 'establishing a database connection...'
	@go run ./cmd/web -secret=${SECRET} -db-dsn=${JOY_DB_DSN} 

## db/mysql: connect to the databse using mysql
.PHONY: db/mysql
db/mysql:
	mysql ${JOY_DB_DSN}


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...


# ==================================================================================== #
# BUILD
# ==================================================================================== #
## build/api: build the cmd/api application
.PHONY: build/web
build/web:
	@echo 'Building cmd/web...'
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/web ./cmd/web

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #
production_host_ip = '137.184.228.246'
## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh joy@${production_host_ip}

