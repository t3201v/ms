# Directories
PROTO_DIR = .
DESCRIPTOR_DIR = descriptor
GOOGLEAPIS_DIR = ../googleapis

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
	protoc -I$(GOOGLEAPIS_DIR) -I. \
		--go_out=identity/gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=identity/gen \
		--go-grpc_opt=paths=source_relative \
		--descriptor_set_out=$(DESCRIPTOR_DIR)/identity.pb \
		--include_imports \
		--include_source_info \
		identity/proto/identity.proto

	@echo "Generating protos for resource/"
	protoc -I$(GOOGLEAPIS_DIR) -I. \
		--go_out=resource/gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=resource/gen \
		--go-grpc_opt=paths=source_relative \
		--descriptor_set_out=$(DESCRIPTOR_DIR)/resource.pb \
		--include_imports \
		--include_source_info \
		resource/proto/resource.proto

# Generate only descriptor files (for Envoy)
descriptors: $(DESCRIPTOR_DIR)
	@echo "Generating descriptor for identity service"
	protoc -I$(GOOGLEAPIS_DIR) -I. \
		--descriptor_set_out=$(DESCRIPTOR_DIR)/identity.pb \
		--include_imports \
		--include_source_info \
		identity/proto/identity.proto

	@echo "Generating descriptor for resource service"
	protoc -I$(GOOGLEAPIS_DIR) -I. \
		--descriptor_set_out=$(DESCRIPTOR_DIR)/resource.pb \
		--include_imports \
		--include_source_info \
		resource/proto/resource.proto

# Generate combined descriptor file
combined-descriptor: $(DESCRIPTOR_DIR)
	@echo "Generating combined descriptor"
	protoc -I$(GOOGLEAPIS_DIR) -I. \
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

.PHONY: proto-gen descriptors combined-descriptor clean install-tools