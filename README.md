# cep-racer

CLI em Go para consulta de CEP que faz uma "corrida" entre duas fontes públicas (BrasilAPI e ViaCEP) e retorna a primeira resposta válida. Útil para obter endereço rapidamente com baixa latência.

## Requisitos

- Go 1.22 ou superior
- Acesso à internet (as consultas são feitas em tempo real)

## Como executar sem build

Execute com `go run` passando o CEP (com ou sem hífen):

```bash
go run ./cmd/racer 01001000
go run ./cmd/racer 01001-000
```

### Saída de exemplo (texto)

```
[brasilapi]
CEP: 01001000
UF: SP
Cidade: São Paulo
Bairro: Sé
Logradouro: Praça da Sé
```

## Flags úteis

- `-timeout` (padrão: `1s`): tempo máximo total da consulta.
  - Exemplos: `-timeout=500ms`, `-timeout=2s`.
- `-json`: imprime a saída em JSON (uma linha por consulta).

Exemplos:

```bash
go run ./cmd/racer -json 01001000
go run ./cmd/racer -timeout=2s 01001-000
```

### Saída de exemplo (JSON)

```json
{
  "cep": "01001000",
  "uf": "SP",
  "cidade": "São Paulo",
  "bairro": "Sé",
  "logradouro": "Praça da Sé",
  "fonte": "brasilapi"
}
```

## Build opcional (binário)

```bash
go build -o cep-racer ./cmd/racer
./cep-racer 01001000
./cep-racer -json -timeout=2s 01001000
```

## Como funciona

- Duas requisições são disparadas em paralelo: BrasilAPI e ViaCEP.
- A primeira resposta válida vence; a outra é cancelada.
- Timeout global configurável via `-timeout` (padrão 1s).
- O CEP de entrada é normalizado para conter apenas dígitos e deve ter 8 dígitos.

## Tratamento de erros

- CEP inválido: o programa exige exatamente 8 dígitos após normalização e encerra com código 2.
- Timeout: se nenhuma fonte responder a tempo, uma mensagem de timeout é exibida e o programa encerra com erro.
- Erro nas fontes: se ambas falharem (ex.: CEP inexistente em ambas), são exibidos os dois erros.

## Observações

- Formatos de entrada aceitos: `12345678` ou `12345-678` (outros caracteres são ignorados na normalização, mas ainda é necessário totalizar 8 dígitos).
- Campos mínimos esperados de cada fonte: `cep` e `uf`. Respostas incompletas são tratadas como erro.

---

Projeto de exemplo para demonstrar concorrência e corrida de requests em Go.

