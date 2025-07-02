.PHONY: test build clean run-example run-webhook help

# Variables
MODULE := $(shell head -1 go.mod | cut -d' ' -f2)
BUILD_DIR := ./bin

help: ## Mostrar esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $1, $2}'

test: ## Ejecutar tests
	go test -v ./...

test-coverage: ## Ejecutar tests con cobertura
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

build: ## Compilar ejemplos
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/basic ./examples/basic/
	go build -o $(BUILD_DIR)/advanced ./examples/advanced/
	go build -o $(BUILD_DIR)/webhook-server ./examples/webhooks/

clean: ## Limpiar archivos generados
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

run-example: ## Ejecutar ejemplo básico
	@if [ -z "$(ESIM_GO_API_KEY)" ]; then \
		echo "❌ Error: ESIM_GO_API_KEY no está configurado"; \
		echo "Configúralo con: export ESIM_GO_API_KEY=tu-api-key"; \
		exit 1; \
	fi
	go run ./examples/basic/

run-advanced: ## Ejecutar ejemplo avanzado
	@if [ -z "$(ESIM_GO_API_KEY)" ]; then \
		echo "❌ Error: ESIM_GO_API_KEY no está configurado"; \
		echo "Configúralo con: export ESIM_GO_API_KEY=tu-api-key"; \
		exit 1; \
	fi
	go run ./examples/advanced/

run-webhook: ## Ejecutar servidor de webhooks
	go run ./examples/webhooks/

fmt: ## Formatear código
	go fmt ./...

lint: ## Ejecutar linter
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint no está instalado. Instalalo con:"; \
		echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

mod-tidy: ## Limpiar go.mod
	go mod tidy

docs: ## Generar documentación
	@if command -v godoc >/dev/null 2>&1; then \
		echo "📚 Documentación disponible en: http://localhost:6060/pkg/$(MODULE)"; \
		godoc -http=:6060; \
	else \
		echo "⚠️  godoc no está instalado. Instalalo con:"; \
		echo "go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

install-tools: ## Instalar herramientas de desarrollo
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/godoc@latest
