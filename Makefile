PROTO_ROOT = submodules/InventoryManagement.Protos
OUTPUT_ROOT = pkg

PROTO_FILES = $(shell find $(PROTO_ROOT) -name '*.proto')

PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc

# Define the protoc flags
PROTOC_FLAGS := -I$(PROTO_ROOT)
GO_OUT_FLAGS := --go_out=paths=source_relative:$(OUTPUT_ROOT)
GRPC_OUT_FLAGS := --go-grpc_out=paths=source_relative:$(OUTPUT_ROOT)

# Target to generate Go code from all .proto files
.PHONY: proto
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

# Ensure the protoc plugins are installed
.PHONY: install-proto-plugins
install-proto-plugins:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Clean generated proto files
.PHONY: clean-proto
clean-proto:
	rm -rf $(OUTPUT_ROOT)

env-example:
	@echo "Creating .env.example file"
	@sed 's/=.*/=/' .env > .env.example
	@echo ".env.example file created."

fetch-latest-submodules:
	@git submodule update --remote --merge
	@git add submodules
	@git commit -m "Updated submodules to latest version."