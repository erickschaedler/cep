package cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type viaCEPResp struct {
	CEP        string `json:"cep"`
	UF         string `json:"uf"`
	Localidade string `json:"localidade"`
	Bairro     string `json:"bairro"`
	Logradouro string `json:"logradouro"`
	Erro       bool   `json:"erro,omitempty"`
}

func FetchViaCEP(ctx context.Context, cep string) Result {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	var out Result
	out.Addr.Fonte = "viacep"

	cli := &http.Client{
		Timeout: 950 * time.Millisecond,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		out.Err = err
		return out
	}

	resp, err := cli.Do(req)
	if err != nil {
		out.Err = err
		return out
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		out.Err = fmt.Errorf("viacep status %d", resp.StatusCode)
		return out
	}

	var v viaCEPResp
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		out.Err = err
		return out
	}
	if v.Erro {
		out.Err = errors.New("viacep: cep não encontrado")
		return out
	}

	out.Addr = Address{
		CEP:        safeCEP(v.CEP),
		UF:         v.UF,
		Cidade:     v.Localidade,
		Bairro:     v.Bairro,
		Logradouro: v.Logradouro,
		Fonte:      "viacep",
	}

	if out.Addr.CEP == "" || out.Addr.UF == "" {
		out.Err = errors.New("viacep: resposta incompleta")
	}
	return out
}
