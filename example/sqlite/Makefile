MIGRATION_DIRECTORY = migration

.DEFAULT_GOAL = build

build:
	for dir in $$(ls $(MIGRATION_DIRECTORY)); do \
		cd $(MIGRATION_DIRECTORY)/$$dir; \
		GO111MODULE=on go build -buildmode=plugin; \
	done
.PHONY: build
