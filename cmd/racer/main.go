package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cep-racer/internal/cep"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("uso: go run ./cmd/racer <CEP>")
		os.Exit(2)
	}
	cepInput := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	results := make(chan cep.Result, 2)

	go func() { results <- cep.FetchBrasilAPI(ctx, cepInput) }()
	go func() { results <- cep.FetchViaCEP(ctx, cepInput) }()

	select {
	case res := <-results:
		if res.Err != nil {
			select {
			case res2 := <-results:
				if res2.Err != nil {
					fmt.Printf("Falhou nas duas fontes: %v | %v\n", res.Err, res2.Err)
					os.Exit(1)
				}
				printAddress(res2.Addr)
				return
			case <-ctx.Done():
				fmt.Printf("Timeout em 1s (primeiro erro: %v)\n", res.Err)
				os.Exit(1)
			}
		}
		cancel()
		printAddress(res.Addr)
	case <-ctx.Done():
		fmt.Println("Timeout em 1s — nenhuma API respondeu a tempo.")
		os.Exit(1)
	}
}

func printAddress(a cep.Address) {
	fmt.Printf("[%s]\n", a.Fonte)
	fmt.Printf("CEP: %s\n", a.CEP)
	fmt.Printf("UF: %s\n", a.UF)
	fmt.Printf("Cidade: %s\n", a.Cidade)
	fmt.Printf("Bairro: %s\n", a.Bairro)
	fmt.Printf("Logradouro: %s\n", a.Logradouro)
}
