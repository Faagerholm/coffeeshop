CURRENT_DIR := $(dir $(MAKEFILE_LIST)) PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go

$(PROTOC_GEN_GO):
	go install github.com/golang/protobuf/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

$(PROTOC_GEN_GO_GRPC):
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

all: help

help:
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for help_line in $${help_lines[@]}; do \
		IFS=$$'#' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf "%-30s %s\n" $$help_command $$help_info ; \
	done

.PHONY: build-identity-api
build-identity-api: ## Build identity proto
	@echo "Building identity proto.."
	@protoc \
		--go_out=identity \
		--go-grpc_out=identity \
		coffeeshop.proto	
	@echo "Done!"

.PHONY: build-checkout-api
build-checkout-api: ## Build checkout proto
	@echo "Building checkout proto.."
	@protoc \
		--go_out=checkout \
		--go-grpc_out=checkout \
		coffeeshop.proto	
	@echo "Done!"

.PHONY: build-warehouse-api
build-warehouse-api: ## Build warehouse proto
	@echo "Building warehouse proto.."
	@protoc \
		--go_out=warehouse \
		--go-grpc_out=warehouse \
		coffeeshop.proto	
	@echo "Done!"

.PHONY: build-cart-api
build-cart-api: ## Build cart proto
	@echo "Building cart proto.."
	@protoc \
		--go_out=cart \
		--go-grpc_out=cart \
		coffeeshop.proto 
	@echo "Done!"

.PHONY: build-shipping-api
build-shipping-api: ## Build shipping proto
	@echo "Building shipping proto.."
	@protoc \
		--go_out=shipping \
		--go-grpc_out=shipping \
		coffeeshop.proto 
	@echo "Done!"

.PHONY: build-payment-api
build-payment-api: ## Build payment proto
	@echo "Building payment proto.."
	@protoc \
		--go_out=payment \
		--go-grpc_out=payment \
		coffeeshop.proto 
	@echo "Done!"

.PHONY: build-email-api
build-email-api: ## Build email proto
	@echo "Building email proto.."
	@protoc \
		--go_out=email \
		--go-grpc_out=email \
		coffeeshop.proto 
	@echo "Done!"

.PHONY: build-frontend-api
build-frontend-api: ## Build frontend proto
	@echo "Building frontend proto.."
	@protoc \
		--go_out=frontend \
		--go-grpc_out=frontend \
		--grpc-gateway_out frontend/proto \
		--grpc-gateway_opt paths=source_relative \
		-I ./ \
		coffeeshop.proto 
	@echo "Done!"


.PHONY: install-tools
install-tools: $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC) ## Install tools

.PHONY: remove-tools 
remove-tools: ## Remove tools
	@echo "Removing tools.."
	@rm -rf $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC)
	@echo "Done!"

