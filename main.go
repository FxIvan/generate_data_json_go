package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Estructura del cliente
type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Estructura de la cuenta bancaria
type Account struct {
	ID         int     `json:"id"`
	CustomerID int     `json:"customer_id"`
	Type       string  `json:"type"`
	Balance    float64 `json:"balance"`
}

// Estructura de la transacción
type Transaction struct {
	ID          int     `json:"id"`
	FromAccount int     `json:"from_account"`
	ToAccount   int     `json:"to_account"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
}

func main() {
	// Crear archivo de salida
	file, err := os.Create("bank_data.json")
	if err != nil {
		log.Fatalf("Error creando el archivo: %v", err)
	}
	defer file.Close()

	// Configurar generador de números aleatorios
	rand.Seed(time.Now().UnixNano())

	// Escribir el inicio del JSON
	_, _ = file.WriteString("{\n\"customers\": [")

	totalSize := int64(0)
	maxSize := int64(1 * 1024 * 1024 * 1024) // 1GB en bytes
	customerCount := 200000                  // Incrementar a 200,000 clientes
	accountCount := 400000                   // Incrementar a 400,000 cuentas
	transactionCount := 10000000             // Incrementar a 10,000,000 transacciones

	// Generar clientes
	var buffer bytes.Buffer
	for i := 0; i < customerCount; i++ {
		customer := Customer{
			ID:      i,
			Name:    "Customer_" + strconv.Itoa(i),
			Address: "Address_" + strconv.Itoa(i),
		}
		if i > 0 {
			_, _ = file.WriteString(",")
		}
		json.NewEncoder(&buffer).Encode(customer)
		_, _ = file.Write(buffer.Bytes())
		buffer.Reset()
	}
	_, _ = file.WriteString("],\n\"accounts\": [")

	// Generar cuentas
	for i := 0; i < accountCount; i++ {
		account := Account{
			ID:         i,
			CustomerID: rand.Intn(customerCount),
			Type:       randomAccountType(),
			Balance:    randomFloat(1000, 100000),
		}
		if i > 0 {
			_, _ = file.WriteString(",")
		}
		json.NewEncoder(&buffer).Encode(account)
		_, _ = file.Write(buffer.Bytes())
		buffer.Reset()
	}
	_, _ = file.WriteString("],\n\"transactions\": [")

	// Generar transacciones
	for i := 0; i < transactionCount; i++ {
		transaction := Transaction{
			ID:          i,
			FromAccount: rand.Intn(accountCount),
			ToAccount:   rand.Intn(accountCount),
			Amount:      randomFloat(10, 5000),
			Date:        randomDate(),
		}
		if i > 0 {
			_, _ = file.WriteString(",")
		}
		json.NewEncoder(&buffer).Encode(transaction)
		_, _ = file.Write(buffer.Bytes())
		buffer.Reset()

		// Verificar tamaño acumulado
		totalSize += int64(len(buffer.Bytes()))
		if totalSize > maxSize {
			break
		}
	}

	// Escribir el cierre del JSON
	_, _ = file.WriteString("]\n}")
	log.Printf("Archivo generado con tamaño aproximado de %d bytes", totalSize)
}

// Generar un tipo de cuenta aleatorio
func randomAccountType() string {
	types := []string{"Savings", "Checking", "Credit"}
	return types[rand.Intn(len(types))]
}

// Generar un número aleatorio en un rango
func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Generar una fecha aleatoria
func randomDate() string {
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()
	diff := end.Sub(start).Seconds()
	randomSeconds := rand.Int63n(int64(diff))
	return start.Add(time.Duration(randomSeconds) * time.Second).Format("2006-01-02")
}
