include .envrc

MAKEFLAGS += --no-print-directory

LDFLAGS=-ldflags="-s -w"
BUILD_OPTS=-installsuffix cgo
CGO=CGO_ENABLED=0
BIN_DIR=./bin

LINUX_ARCHS=386 amd64 arm arm64
DARWIN_ARCHS=amd64 arm64
WINDOWS_ARCHS=386 amd64 arm

SWAGGER_CONF=./config/swagger.config.yml
SVCS=soccer-manager

APP_NAME=soccer-manager

# Build executable for your arch and os
build:
ifeq ($(OS),Windows_NT) 
	@$(MAKE) build/batch name=${APP_NAME}
else
	@$(MAKE) build/bash name=${APP_NAME}
endif

# Build executable for every arch and os
build/all:
ifeq ($(OS),Windows_NT) 
	@$(MAKE) build/all/batch name=${APP_NAME}
else
	@$(MAKE) build/all/bash name=${APP_NAME}
endif

# Generates swagger config file according to SVCS
swag-config:
ifeq ($(OS),Windows_NT) 
	@$(MAKE) swag-config/batch
else
	@$(MAKE) swag/bash
endif

# Generates swagger documentation file according to $(svc)
swag-gen:
ifeq ($(OS),Windows_NT) 
	@$(MAKE) swag-gen/batch svc=$(svc)
else
	@$(MAKE) swag-gen/bash svc=$(svc)
endif

# Default build scripts linked to your current os and arch
build/batch:
	set CGO_ENABLED=0&& go build -ldflags="-s -w" -installsuffix cgo -o ./bin/$(name).exe ./cmd/$(name)/main.go

build/bash:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o ./bin/$(name) ./cmd/$(name)/main.go

# Run all tests
test: 
	go test -cover ./internal/...

build/all/bash:
	@$(MAKE) build/linux/bash name=$(name)
	@$(MAKE) build/darwin/bash name=$(name)
	@$(MAKE) build/windows/bash name=$(name)

build/darwin/bash:
	@$(MAKE) build/unix/bash archs='$(DARWIN_ARCHS)' os=darwin name=$(name)

build/linux/bash:
	@$(MAKE) build/unix/bash archs='$(LINUX_ARCHS)' os=linux name=$(name)

build/unix/bash:
	$(foreach arch,$(archs),${CGO} GOOS=$(os) GOARCH=$(arch) go build ${LDFLAGS} ${BUILD_OPTS} -o ${BIN_DIR}/$(name)_$(os)_$(arch) ./cmd/$(name)/main.go;)

build/windows/bash:
	$(foreach arch,${WINDOWS_ARCHS},${CGO} GOOS=windows GOARCH=$(arch) go build ${LDFLAGS} ${BUILD_OPTS} -o ${BIN_DIR}/$(name)_windows_$(arch).exe ./cmd/$(name)/main.go;)

build/all/batch:
	@$(MAKE) build/linux/batch name=$(name)
	@$(MAKE) build/darwin/batch name=$(name)
	@$(MAKE) build/windows/batch name=$(name)

build/linux/batch:
	@$(MAKE) build/unix/batch archs='${LINUX_ARCHS}' os=linux name=$(name)
	
build/darwin/batch:
	@$(MAKE) build/unix/batch archs='${DARWIN_ARCHS}' os=darwin name=$(name)

build/windows/batch:
	for %%f in (${WINDOWS_ARCHS}) do set ${CGO}&& set GOOS=windows&& set GOARCH=%%f&& go build ${LDFLAGS} ${BUILD_OPTS} -o ${BIN_DIR}/$(name)_windows_%%f.exe ./cmd/$(name)/main.go

build/unix/batch:
	for %%f in ($(archs)) do set ${CGO}&& set GOOS=$(os)&& set GOARCH=%%f&& go build ${LDFLAGS} ${BUILD_OPTS} -o ${BIN_DIR}/$(name)_$(os)_%%f ./cmd/$(name)/main.go

# Genrates sqlc files according to $(db)
sqlc:
	@sqlc generate -f ./sql/$(db)/sqlc.yml

# Migrates database up according to $(db) and credentials from .envrc
migrate-up:
	@$(MAKE) migrate db=$(db) way=up

# Migrates database down according to $(db) and credentials from .envrc
migrate-down:
	@$(MAKE) migrate db=$(db) way=down

# Migrates database according to $(way){up/down} and $(db), gets credentials and values form .envrc
migrate:
	migrate -path ./sql/$(db)/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5433/$(db)?sslmode=disable" -verbose $(way)

# Initializes migrate up & down files for $(db)
migrate-init:
	migrate create -ext sql -dir ./sql/$(db)/migrations/ -seq init_schema

# PLATFORM SPECIFIC
swag-gen/bash:
	@swag init --dir cmd/$(svc)/,internal/$(svc)/delivery --parseDependency --output ./api/swagger --outputTypes yaml
	mv ./api/swagger/swagger.yaml ./api/swagger/$(svc).swagger.yaml

swag-gen/batch:
	@swag init --dir cmd/$(svc)/,internal/$(svc)/delivery --parseDependency --output ./api/swagger --outputTypes yaml
	@if exist api\swagger\$(svc).swagger.yaml del api\swagger\$(svc).swagger.yaml
	ren api\swagger\swagger.yaml $(svc).swagger.yaml

swag-config/bash:
	@echo "urls:" > $(SWAGGER_CONF)
	@for svc in $(SVCS); do \
		echo "  - url: "./$$svc.swagger.yaml"" >> $(SWAGGER_CONF); \
		echo "    name: "$$svc"" >> $(SWAGGER_CONF); \
	done

swag-config/batch:
	@echo urls: > $(SWAGGER_CONF)
	@for %%s in ($(SVCS)) do ( \
		echo   - url: "./%%s.swagger.yaml" >> $(SWAGGER_CONF) && \
		echo     name: "%%s" >> $(SWAGGER_CONF) \
	)
