# Directory definitions
PROTO_ROOT := pkg/api
OUTPUT_ROOT := pkg/api

# Find all .proto files
PROTO_FILES := $(shell find $(PROTO_ROOT) -name '*.proto')

# Tool definitions
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc

# Protoc flags
PROTOC_FLAGS := -I$(PROTO_ROOT) -I/usr/include
GO_OUT_FLAGS := --go_out=paths=source_relative:$(OUTPUT_ROOT)
GRPC_OUT_FLAGS := --go-grpc_out=paths=source_relative:$(OUTPUT_ROOT)

# Phony targets
.PHONY: proto env-example

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

env-example:
	@echo "Creating .env.example file"
	@sed 's/=.*/=/' .env > .env.example
	@echo ".env.example file created."