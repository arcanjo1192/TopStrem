package handlers

import (
	"net/http"
	"strings"
)

// ResponseFormat define o formato de resposta desejado
type ResponseFormat string

const (
	FormatHTML ResponseFormat = "html"
	FormatJSON ResponseFormat = "json"
)

// NegotiateFormat detecta o formato de resposta preferido baseado no header Accept
// Prioridade:
// 1. application/json -> JSON
// 2. text/html ou */* (padrão) -> HTML
func NegotiateFormat(r *http.Request) ResponseFormat {
	accept := r.Header.Get("Accept")
	
	// Se explicitamente pede JSON
	if strings.Contains(accept, "application/json") {
		return FormatJSON
	}
	
	// Padrão: HTML (para compatibilidade com browsers)
	return FormatHTML
}

// IsJSONRequest verifica se a requisição é um pedido de JSON
func IsJSONRequest(r *http.Request) bool {
	return NegotiateFormat(r) == FormatJSON
}
