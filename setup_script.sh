#!/bin/bash

# Script para configurar el repositorio eSIM Go Client
# Ejecuta este script después de clonar tu repositorio de GitHub

set -e

echo "🚀 Configurando el repositorio eSIM Go Client..."

# Verificar que estamos en un directorio de Git
if [ ! -d ".git" ]; then
    echo "❌ Error: Este script debe ejecutarse en la raíz de un repositorio Git"
    echo "Por favor, primero clona tu repositorio:"
    echo "git clone https://github.com/tu-username/esimgo-client.git"
    echo "cd esimgo-client"
    exit 1
fi

# Obtener información del repositorio
REPO_URL=$(git config --get remote.origin.url)
MODULE_NAME=$(echo $REPO_URL | sed 's/.*github\.com\///; s/\.git$//')

echo "📦 Módulo detectado: $MODULE_NAME"

# Crear go.mod
echo "📝 Creando go.mod..."
cat > go.mod << EOF
module github.com/$MODULE_NAME

go 1.21

require ()
EOF

# Crear estructura de directorios
echo "📁 Creando estructura de directorios..."
mkdir -p examples/basic
mkdir -p examples/advanced  
mkdir -p examples/webhooks
mkdir -p docs

# Crear client.go
echo "📝 Creando client.go..."
cat > client.go << 'EOF'
package esimgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client represents the eSIM Go API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new eSIM Go API client
func NewClient(apiKey string) *Client {
	return &Client{
		baseURL:    "https://api.esim-go.com/v2.4",
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SetHTTPClient allows setting a custom HTTP client
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// SetBaseURL allows setting a custom base URL (useful for testing)
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = strings.TrimSuffix(baseURL, "/")
}

// makeRequest performs HTTP requests with proper authentication
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
		}
		apiErr.StatusCode = resp.StatusCode
		return &apiErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// APIError represents an API error response
type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// Common structures used across the API

// Country represents a country
type Country struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	ISO    string `json:"iso"`
}

// Bundle represents a data bundle
type Bundle struct {
	ID          string `json:"id,omitempty"`
	Reference   string `json:"reference,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Bundle data details
	InitialQuantity   int64     `json:"initialQuantity,omitempty"`
	RemainingQuantity int64     `json:"remainingQuantity,omitempty"`
	StartTime         time.Time `json:"startTime,omitempty"`
	EndTime           time.Time `json:"endTime,omitempty"`
	Unlimited         bool      `json:"unlimited,omitempty"`
}

// Assignment represents a bundle assignment
type Assignment struct {
	ID                   string    `json:"id"`
	CallTypeGroup        string    `json:"callTypeGroup"`
	InitialQuantity      int64     `json:"initialQuantity"`
	RemainingQuantity    int64     `json:"remainingQuantity"`
	AssignmentDateTime   time.Time `json:"assignmentDateTime"`
	AssignmentReference  string    `json:"assignmentReference"`
	BundleState          string    `json:"bundleState"`
	Unlimited            bool      `json:"unlimited"`
}

// ESIM represents an eSIM
type ESIM struct {
	ICCID                   string    `json:"iccid"`
	PIN                     string    `json:"pin,omitempty"`
	PUK                     string    `json:"puk,omitempty"`
	MatchingID              string    `json:"matchingId,omitempty"`
	SMDPAddress             string    `json:"smdpAddress,omitempty"`
	ProfileStatus           string    `json:"profileStatus,omitempty"`
	FirstInstalledDateTime  int64     `json:"firstInstalledDateTime,omitempty"`
	CustomerRef             string    `json:"customerRef,omitempty"`
	LastAction              string    `json:"lastAction,omitempty"`
	ActionDate              string    `json:"actionDate,omitempty"`
	Physical                bool      `json:"physical,omitempty"`
	AssignedDate            string    `json:"assignedDate,omitempty"`
}

// ESIMGoClient is the main client that provides access to all services
type ESIMGoClient struct {
	*Client
	ESIMs        *ESIMService
	Organization *OrganizationService
	Orders       *OrdersService
	Inventory    *InventoryService
	Catalogue    *CatalogueService
	Networks     *NetworksService
}

// NewESIMGoClient creates a new complete eSIM Go API client
func NewESIMGoClient(apiKey string) *ESIMGoClient {
	baseClient := NewClient(apiKey)
	
	return &ESIMGoClient{
		Client:       baseClient,
		ESIMs:        NewESIMService(baseClient),
		Organization: NewOrganizationService(baseClient),
		Orders:       NewOrdersService(baseClient),
		Inventory:    NewInventoryService(baseClient),
		Catalogue:    NewCatalogueService(baseClient),
		Networks:     NewNetworksService(baseClient),
	}
}

// Common constants for API usage
const (
	// Bundle states
	BundleStateProcessing = "Processing"
	BundleStateQueued     = "Queued"
	BundleStateActive     = "Active"
	BundleStateDepleted   = "Depleted"
	BundleStateExpired    = "Expired"
	BundleStateRevoked    = "Revoked"
	BundleStateLapsed     = "Lapsed"
	
	// Order types
	OrderTypeValidate    = "validate"
	OrderTypeTransaction = "transaction"
	
	// Bundle types
	BundleTypeBundle = "bundle"
	
	// Direction constants
	DirectionAsc  = "asc"
	DirectionDesc = "desc"
	
	// Duration units
	DurationUnitDay   = "day"
	DurationUnitMonth = "month"
)
EOF

echo "✅ Archivos principales creados exitosamente!"

# Crear README.md mejorado
echo "📝 Creando README.md..."
cat > README.md << EOF
# eSIM Go API Client

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/$MODULE_NAME)](https://goreportcard.com/report/github.com/$MODULE_NAME)

Cliente Go completo para la API de eSIM Go v2.4. Este cliente proporciona una interfaz fácil de usar para todas las operaciones de la API de eSIM Go.

## 🚀 Características

- ✅ Soporte completo para la API v2.4 de eSIM Go
- ✅ Tipado fuerte con estructuras Go
- ✅ Manejo automático de autenticación con API Key
- ✅ Manejo de errores estructurado
- ✅ Contexto para cancelación y timeouts
- ✅ Sin dependencias externas
- ✅ Servicios organizados por funcionalidad

## 📦 Instalación

\`\`\`bash
go get github.com/$MODULE_NAME
\`\`\`

## 🔧 Uso Básico

\`\`\`go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/$MODULE_NAME"
)

func main() {
    // Crear cliente con tu API key
    client := esimgo.NewESIMGoClient("your-api-key-here")
    
    ctx := context.Background()
    
    // Obtener detalles de la organización
    org, err := client.Organization.GetDetails(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Organización: %s, Balance: %d %s\\n", 
        org.Name, org.Balance, org.Currency)
}
\`\`\`

## 📚 Servicios Disponibles

### ESIMs (\`client.ESIMs\`)
- \`ApplyBundle()\` - Aplicar bundle a eSIM
- \`GetDetails()\` - Obtener detalles de eSIM
- \`List()\` - Listar eSIMs
- \`GetHistory()\` - Obtener historial
- \`SendSMS()\` - Enviar SMS
- \`GetLocation()\` - Obtener ubicación
- \`ListBundles()\` - Listar bundles aplicados
- \`RevokeBundle()\` - Revocar bundle
- \`CheckCompatibility()\` - Verificar compatibilidad

### Organización (\`client.Organization\`)
- \`GetDetails()\` - Obtener detalles de organización
- \`GetBundleGroups()\` - Obtener grupos de bundles
- \`TopupBalance()\` - Recargar balance

### Órdenes (\`client.Orders\`)
- \`Create()\` - Crear orden
- \`List()\` - Listar órdenes
- \`GetByReference()\` - Obtener orden por referencia
- \`Validate()\` - Validar orden

### Inventario (\`client.Inventory\`)
- \`Get()\` - Obtener inventario
- \`Refund()\` - Reembolsar bundle

### Catálogo (\`client.Catalogue\`)
- \`List()\` - Listar catálogo
- \`GetBundleDetails()\` - Detalles de bundle
- \`SearchByCountry()\` - Buscar por país
- \`SearchByRegion()\` - Buscar por región

### Redes (\`client.Networks\`)
- \`GetCountryNetworks()\` - Obtener redes por país
- \`GetAllNetworks()\` - Obtener todas las redes

## 📖 Ejemplos

Ver el directorio \`examples/\` para ejemplos completos:

- \`examples/basic/\` - Uso básico del cliente
- \`examples/advanced/\` - Ejemplos avanzados
- \`examples/webhooks/\` - Manejo de webhooks

## 🧪 Tests

\`\`\`bash
go test -v ./...
\`\`\`

## 🤝 Contribuir

1. Fork el repositorio
2. Crea una rama para tu feature (\`git checkout -b feature/nueva-funcionalidad\`)
3. Commit tus cambios (\`git commit -am 'Agregar nueva funcionalidad'\`)
4. Push a la rama (\`git push origin feature/nueva-funcionalidad\`)
5. Abre un Pull Request

## 📄 Licencia

MIT License - ver el archivo LICENSE para detalles.

## 🆘 Soporte

- Documentación oficial de eSIM Go: https://docs.esim-go.com
- Issues: https://github.com/$MODULE_NAME/issues

EOF

# Crear ejemplo básico
echo "📝 Creando ejemplo básico..."
cat > examples/basic/main.go << EOF
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/$MODULE_NAME"
)

func main() {
    apiKey := os.Getenv("ESIM_GO_API_KEY")
    if apiKey == "" {
        log.Fatal("ESIM_GO_API_KEY environment variable is required")
    }

    client := esimgo.NewESIMGoClient(apiKey)
    ctx := context.Background()

    // Ejemplo básico: obtener detalles de organización
    fmt.Println("🏢 Obteniendo detalles de la organización...")
    org, err := client.Organization.GetDetails(ctx)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Printf("✅ Organización: %s\\n", org.Name)
    fmt.Printf("💰 Balance: %d %s\\n", org.Balance, org.Currency)

    // Listar algunos bundles del catálogo
    fmt.Println("\\n📦 Listando catálogo de bundles...")
    bundles, err := client.Catalogue.List(ctx, &esimgo.ListCatalogueRequest{
        PerPage: 5,
    })
    if err != nil {
        log.Printf("Error obteniendo catálogo: %v", err)
    } else {
        fmt.Printf("✅ Encontrados %d bundles\\n", len(bundles))
        for i, bundle := range bundles {
            if i >= 3 { // Mostrar solo los primeros 3
                break
            }
            fmt.Printf("  - %s: %s (Precio: %d)\\n", bundle.Name, bundle.Description, bundle.Price)
        }
    }
}
EOF

# Crear .env.example
echo "📝 Creando .env.example..."
cat > .env.example << EOF
# Tu API Key de eSIM Go
# Obtén tu API key en: https://sso.esim-go.com/login
ESIM_GO_API_KEY=your-api-key-here
EOF

# Ejecutar go mod tidy
echo "🔧 Ejecutando go mod tidy..."
go mod tidy

# Verificar que todo compila
echo "🔍 Verificando que el código compila..."
if go build ./...; then
    echo "✅ ¡Código compila correctamente!"
else
    echo "❌ Error en la compilación"
    exit 1
fi

echo ""
echo "🎉 ¡Configuración completada exitosamente!"
echo ""
echo "📋 Próximos pasos:"
echo "1. Revisa los archivos generados"
echo "2. Copia el contenido restante de los servicios (esims.go, orders.go, etc.)"
echo "3. Configura tu API key: cp .env.example .env"
echo "4. Ejecuta el ejemplo: cd examples/basic && ESIM_GO_API_KEY=tu-key go run main.go"
echo "5. Commit y push los cambios:"
echo "   git add ."
echo "   git commit -m 'Initial commit: eSIM Go client library'"
echo "   git push origin main"
echo ""
echo "🔗 Tu repositorio estará disponible en: https://github.com/$MODULE_NAME"
