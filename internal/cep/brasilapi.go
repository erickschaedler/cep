package cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type brasilAPIResp struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

func FetchBrasilAPI(ctx context.Context, cep string) Result {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	var out Result
	out.Addr.Fonte = "brasilapi"

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
		out.Err = fmt.Errorf("brasilapi status %d", resp.StatusCode)
		return out
	}

	var b brasilAPIResp
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		out.Err = err
		return out
	}

	out.Addr = Address{
		CEP:        safeCEP(b.CEP),
		UF:         b.State,
		Cidade:     b.City,
		Bairro:     b.Neighborhood,
		Logradouro: b.Street,
		Fonte:      "brasilapi",
	}

	if out.Addr.CEP == "" || out.Addr.UF == "" {
		out.Err = errors.New("brasilapi: resposta incompleta")
	}
	return out
}
