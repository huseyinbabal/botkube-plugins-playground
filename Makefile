EXECUTORS_PLUGIN_DIRS=$(wildcard ./contrib/executors/*)
SOURCES_PLUGIN_DIRS=$(wildcard ./contrib/sources/*)

all: build-plugins

clean: clean-plugins

build-plugins: $(EXECUTORS_PLUGIN_DIRS) $(SOURCES_PLUGIN_DIRS)

clean-plugins:
	rm -f ./contrib/build/*

protoc:
	# Clean
	rm -rf source/proto/*.pb.go executor/proto/*pb.go

	# Compile
	cd plugin/source
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative plugin/source/**/*.proto

	cd plugin/executor
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative plugin/executor/**/*.proto

	## Refresh dependencies
	cd plugin/source/proto && go mod tidy
	cd plugin/executor/proto && go mod tidy

$(EXECUTORS_PLUGIN_DIRS):
	$(MAKE) -C $@

$(SOURCES_PLUGIN_DIRS):
	$(MAKE) -C $@

.PHONY: all $(EXECUTORS_PLUGIN_DIRS) $(SOURCES_PLUGIN_DIRS)