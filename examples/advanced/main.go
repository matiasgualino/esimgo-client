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

	// Ejemplo avanzado: búsqueda de bundles por país
	fmt.Println("🔍 Buscando bundles para España...")
	spainBundles, err := client.Catalogue.List(ctx, &esimgo.ListCatalogueRequest{
		Countries: "ES",
		PerPage:   10,
	})
	if err != nil {
		log.Printf("Error buscando bundles para España: %v", err)
	} else {
		fmt.Printf("✅ Encontrados %d bundles para España\n", len(spainBundles.Bundles))
		for i, bundle := range spainBundles.Bundles {
			if i >= 3 {
				break
			}
			fmt.Printf("  - %s: %dMB por %d días (€%.2f)\n",
				bundle.Name, bundle.DataAmount, bundle.Duration, float64(bundle.Price)/100)
		}
	}

	// Ejemplo: validar una orden
	fmt.Println("\n💰 Validando una orden...")
	if len(spainBundles.Bundles) > 0 {
		orderItems := []esimgo.OrderItem{
			{
				Type:     esimgo.BundleTypeBundle,
				Quantity: 1,
				Item:     spainBundles.Bundles[0].Name,
			},
		}

		validation, err := client.Orders.Validate(ctx, orderItems, false)
		if err != nil {
			log.Printf("Error validando orden: %v", err)
		} else {
			fmt.Printf("✅ Orden válida: %t\n", validation.Valid)
			fmt.Printf("💵 Total: %.2f %s\n", validation.Total, validation.Currency)
		}
	}

	// Ejemplo: obtener inventario
	fmt.Println("\n📦 Consultando inventario...")
	inventory, err := client.Inventory.Get(ctx)
	if err != nil {
		log.Printf("Error obteniendo inventario: %v", err)
	} else {
		fmt.Printf("✅ Bundles en inventario: %d\n", len(inventory.Bundles))
	}

	fmt.Println("\n🎉 ¡Ejemplos completados!")
}
