package cep

type Address struct {
	CEP        string `json:"cep"`
	UF         string `json:"uf"`
	Cidade     string `json:"cidade"`
	Bairro     string `json:"bairro"`
	Logradouro string `json:"logradouro"`
	Fonte      string `json:"fonte"`
}

type Result struct {
    Addr Address
    Err  error
}

func safeCEP(s string) string {
    if len(s) == 9 && s[5] == '-' {
        return s[:5] + s[6:]
    }

    return s
}

func NormalizeCEP(s string) string {
    out := make([]byte, 0, len(s))
    for i := range len(s) {
        c := s[i]
        if c >= '0' && c <= '9' {
            out = append(out, c)
        }
    }
    
    return string(out)
}
