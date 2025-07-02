package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/matiasgualino/esimgo-client"
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

	fmt.Printf("✅ Organización: %s\n", org.Name)
	fmt.Printf("💰 Balance: %d %s\n", org.Balance, org.Currency)

	// Listar algunos bundles del catálogo
	fmt.Println("\n📦 Listando catálogo de bundles...")
	bundles, err := client.Catalogue.List(ctx, &esimgo.ListCatalogueRequest{
		PerPage: 5,
	})
	if err != nil {
		log.Printf("Error obteniendo catálogo: %v", err)
	} else {
		fmt.Printf("✅ Encontrados %d bundles\n", len(bundles))
		for i, bundle := range bundles {
			if i >= 3 { // Mostrar solo los primeros 3
				break
			}
			fmt.Printf("  - %s: %s (Precio: %d)\n", bundle.Name, bundle.Description, bundle.Price)
		}
	}
}
