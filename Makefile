# Directories
PROTO_DIR = .
DESCRIPTOR_DIR = descriptor
GOOGLEAPIS_DIR = ../googleapis
OPENAPIV2_DIR = ../grpc-gateway

# Create output directories
$(DESCRIPTOR_DIR):
	mkdir -p $(DESCRIPTOR_DIR)

identity/gen:
	mkdir -p identity/gen

resource/gen:
	mkdir -p resource/gen

# Generate Go code and descriptor files
proto-gen: identity/gen resource/gen $(DESCRIPTOR_DIR)
	@echo "Generating protos for identity/"
	protoc -I$(GOOGLEAPIS_DIR) -I$(OPENAPIV2_DIR) -I. \
		--go_out=identity/gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=identity/gen \
		--go-grpc_opt=paths=source_relative \
		--openapiv2_out ./docs \
		--descriptor_set_out=$(DESCRIPTOR_DIR)/identity.pb \
		--include_imports \
		--include_source_info \
		identity/proto/identity.proto

	@echo "Generating protos for resource/"
	protoc -I$(GOOGLEAPIS_DIR) -I$(OPENAPIV2_DIR) -I. \
		--go_out=resource/gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=resource/gen \
		--go-grpc_opt=paths=source_relative \
		--openapiv2_out ./docs \
		--descriptor_set_out=$(DESCRIPTOR_DIR)/resource.pb \
		--include_imports \
		--include_source_info \
		resource/proto/resource.proto

# Generate only descriptor files (for Envoy)
descriptors: $(DESCRIPTOR_DIR)
	@echo "Generating descriptor for identity service"
	protoc -I$(GOOGLEAPIS_DIR) -I$(OPENAPIV2_DIR) -I. \
		--descriptor_set_out=$(DESCRIPTOR_DIR)/identity.pb \
		--include_imports \
		--include_source_info \
		identity/proto/identity.proto

	@echo "Generating descriptor for resource service"
	protoc -I$(GOOGLEAPIS_DIR) -I$(OPENAPIV2_DIR) -I. \
		--descriptor_set_out=$(DESCRIPTOR_DIR)/resource.pb \
		--include_imports \
		--include_source_info \
		resource/proto/resource.proto

# Generate combined descriptor file
combined-descriptor: $(DESCRIPTOR_DIR)
	@echo "Generating combined descriptor"
	protoc -I$(GOOGLEAPIS_DIR) -I$(OPENAPIV2_DIR) -I. \
		--descriptor_set_out=$(DESCRIPTOR_DIR)/services.pb \
		--include_imports \
		--include_source_info \
		identity/proto/identity.proto \
		resource/proto/resource.proto

# Clean generated files
clean:
	rm -rf identity/gen resource/gen $(DESCRIPTOR_DIR)

# Install required tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

keys:
	openssl genpkey -algorithm RSA -out private.key -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in private.key -pubout -out public.key
	cd identity && go run ./cmd/jwks/main.go

# for debugging, also change svc domain to host.docker.internal for envoy cfg
run-infra:
	docker compose up postgres envoy -d --no-deps

run-all:
	docker compose up -d --build

down:
	docker compose down

.PHONY: proto-gen descriptors combined-descriptor clean install-tools keys run-infra run-all