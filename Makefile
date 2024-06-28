# Directory definitions
PROTO_ROOT := submodules/InventoryManagement.Protos
OUTPUT_ROOT := pkg

# Find all .proto files
PROTO_FILES := $(shell find $(PROTO_ROOT) -name '*.proto')

# Tool definitions
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc

# Protoc flags
PROTOC_FLAGS := -I$(PROTO_ROOT)
GO_OUT_FLAGS := --go_out=paths=source_relative:$(OUTPUT_ROOT)
GRPC_OUT_FLAGS := --go-grpc_out=paths=source_relative:$(OUTPUT_ROOT)

# Phony targets
.PHONY: proto install-proto-plugins clean-proto env-example fetch-latest-submodules docker-build docker-run

# Default target
all: proto

# Generate Go code from all .proto files
proto: $(PROTO_FILES)
	@for protofile in $^; do \
		output_dir=$(OUTPUT_ROOT)/$$(dirname $${protofile#$(PROTO_ROOT)/}); \
		mkdir -p $$output_dir; \
		$(PROTOC) $(PROTOC_FLAGS) \
			$(GO_OUT_FLAGS) \
			$(GRPC_OUT_FLAGS) \
			--proto_path=$(PROTO_ROOT) \
			$$protofile; \
		echo "Compiled $$protofile"; \
	done

# Install protoc plugins
install-proto-plugins:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Clean generated proto files
clean-proto:
	rm -rf $(OUTPUT_ROOT)

# Create .env.example file
env-example:
	@echo "Creating .env.example file"
	@sed 's/=.*/=/' .env > .env.example
	@echo ".env.example file created."

# Fetch latest submodules
fetch-latest-submodules:
	@git submodule update --remote --merge
	@git add submodules
	@git commit -m "Updated submodules to latest version."