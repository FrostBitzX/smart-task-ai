.PHONY: tidy mod codegen codegen-tag lint

mod:
	go mod tidy

lint:
	golangci-lint run ./...

# OpenAPI code generation
OPENAPI_SPEC := openapi/openapi.yml
# ดึงรายชื่อ tag ทั้งหมดจากไฟล์ OPENAPI_SPEC (ดูจากบรรทัดที่มี `- name:` ใต้ section tags)
OPENAPI_TAGS := $(shell grep '^[[:space:]]*- name:' $(OPENAPI_SPEC) | awk '{print $$3}')

codegen:
	@if [ "$(fname)" != "" ]; then \
		echo "[codegen] Generating code for tag $(fname)..."; \
		$(MAKE) codegen-tag fname=$(fname); \
	else \
		for t in $(OPENAPI_TAGS); do \
			echo "[codegen] Generating code for tag $$t..."; \
			$(MAKE) codegen-tag fname=$$t; \
		done; \
	fi

# Generate code for a single tag (internal use only)
codegen-tag:
	@mkdir -p internal/interfaces/http/$(fname)
	oapi-codegen \
		-generate fiber,types,strict-server,spec \
		-include-tags $(fname) \
		-o internal/interfaces/http/$(fname)/spec.gen.go \
		-package $(fname) \
		$(OPENAPI_SPEC)