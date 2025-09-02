package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "time"

    "cep-racer/internal/cep"
)

func main() {
    timeout := flag.Duration("timeout", time.Second, "tempo máximo total (ex.: 500ms, 2s)")
    jsonOut := flag.Bool("json", false, "imprime o resultado em JSON")
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "uso: go run ./cmd/racer [opções] <CEP>\n")
        fmt.Fprintf(os.Stderr, "opções:\n")
        flag.PrintDefaults()
    }
    flag.Parse()

    if flag.NArg() < 1 {
        flag.Usage()
        os.Exit(2)
    }

    cepInput := cep.NormalizeCEP(flag.Arg(0))
    if len(cepInput) != 8 {
        fmt.Fprintf(os.Stderr, "CEP inválido: informe 8 dígitos (recebido: %q)\n", flag.Arg(0))
        os.Exit(2)
    }

    ctx, cancel := context.WithTimeout(context.Background(), *timeout)
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
                output(res2.Addr, *jsonOut)
                return
            case <-ctx.Done():
                fmt.Printf("Timeout em %s (primeiro erro: %v)\n", (*timeout).String(), res.Err)
                os.Exit(1)
            }
        }
        cancel()
        output(res.Addr, *jsonOut)
    case <-ctx.Done():
        fmt.Printf("Timeout em %s — nenhuma API respondeu a tempo.\n", (*timeout).String())
        os.Exit(1)
    }
}

func output(a cep.Address, asJSON bool) {
    if asJSON {
        enc := json.NewEncoder(os.Stdout)
        enc.SetIndent("", "  ")
        _ = enc.Encode(a)
        return
    }
    fmt.Printf("[%s]\n", a.Fonte)
    fmt.Printf("CEP: %s\n", a.CEP)
    fmt.Printf("UF: %s\n", a.UF)
    fmt.Printf("Cidade: %s\n", a.Cidade)
    fmt.Printf("Bairro: %s\n", a.Bairro)
    fmt.Printf("Logradouro: %s\n", a.Logradouro)
}
