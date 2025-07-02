# eSIM Go API Client

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/matiasgualino/esimgo-client)](https://goreportcard.com/report/github.com/matiasgualino/esimgo-client)

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

```bash
go get github.com/matiasgualino/esimgo-client
```

## 🔧 Uso Básico

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/matiasgualino/esimgo-client"
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
    
    fmt.Printf("Organización: %s, Balance: %d %s\n", 
        org.Name, org.Balance, org.Currency)
}
```

## 📚 Servicios Disponibles

### ESIMs (`client.ESIMs`)
- `ApplyBundle()` - Aplicar bundle a eSIM
- `GetDetails()` - Obtener detalles de eSIM
- `List()` - Listar eSIMs
- `GetHistory()` - Obtener historial
- `SendSMS()` - Enviar SMS
- `GetLocation()` - Obtener ubicación
- `ListBundles()` - Listar bundles aplicados
- `RevokeBundle()` - Revocar bundle
- `CheckCompatibility()` - Verificar compatibilidad

### Organización (`client.Organization`)
- `GetDetails()` - Obtener detalles de organización
- `GetBundleGroups()` - Obtener grupos de bundles
- `TopupBalance()` - Recargar balance

### Órdenes (`client.Orders`)
- `Create()` - Crear orden
- `List()` - Listar órdenes
- `GetByReference()` - Obtener orden por referencia
- `Validate()` - Validar orden

### Inventario (`client.Inventory`)
- `Get()` - Obtener inventario
- `Refund()` - Reembolsar bundle

### Catálogo (`client.Catalogue`)
- `List()` - Listar catálogo
- `GetBundleDetails()` - Detalles de bundle
- `SearchByCountry()` - Buscar por país
- `SearchByRegion()` - Buscar por región

### Redes (`client.Networks`)
- `GetCountryNetworks()` - Obtener redes por país
- `GetAllNetworks()` - Obtener todas las redes

## 📖 Ejemplos

Ver el directorio `examples/` para ejemplos completos:

- `examples/basic/` - Uso básico del cliente
- `examples/advanced/` - Ejemplos avanzados
- `examples/webhooks/` - Manejo de webhooks

## 🧪 Tests

```bash
go test -v ./...
```

## 🤝 Contribuir

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Abre un Pull Request

## 📄 Licencia

MIT License - ver el archivo LICENSE para detalles.

## 🆘 Soporte

- Documentación oficial de eSIM Go: https://docs.esim-go.com
- Issues: https://github.com/matiasgualino/esimgo-client/issues

