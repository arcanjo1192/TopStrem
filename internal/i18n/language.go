package i18n

import (
    "net/http"
    "strings"
)

// idiomas suportados (códigos ISO 639-1)
var supportedLanguages = map[string]bool{
    "pt": true, // português
    "en": true, // inglês
    "es": true, // espanhol
    "fr": true, // francês
    "de": true, // alemão
    "it": true, // italiano
    "ja": true, // japonês
    "zh": true, // chinês (simplificado)
    "ru": true, // russo
    "ar": true, // árabe
    "hi": true, // hindi
    "ko": true, // coreano
}

// DetectLanguage retorna o idioma baseado na requisição:
// 1. Parâmetro 'lang' na query string (se suportado)
// 2. Cabeçalho Accept-Language (primeiro idioma suportado)
// 3. Padrão 'pt'
func DetectLanguage(r *http.Request) string {
    // 1. Prioridade máxima: parâmetro na URL
    if lang := r.URL.Query().Get("lang"); lang != "" {
        if supportedLanguages[lang] {
            return lang
        }
        // se não for suportado, tenta mapear (ex: "pt-BR" -> "pt")
        code := strings.Split(lang, "-")[0]
        if supportedLanguages[code] {
            return code
        }
        // senão, continua para o Accept-Language
    }

    // 2. Tenta o cabeçalho Accept-Language
    acceptLang := r.Header.Get("Accept-Language")
    if acceptLang == "" {
        return "pt"
    }

    // Exemplo: "en-US,en;q=0.9,pt-BR;q=0.8,pt;q=0.7"
    parts := strings.Split(acceptLang, ",")
    for _, part := range parts {
        // extrai o código do idioma (antes de ';' ou ',')
        langCode := strings.Split(part, ";")[0]
        langCode = strings.TrimSpace(langCode)
        // tenta o código completo (ex: "pt-BR")
        if supportedLanguages[langCode] {
            return langCode
        }
        // tenta apenas a parte principal (ex: "pt")
        mainCode := strings.Split(langCode, "-")[0]
        if supportedLanguages[mainCode] {
            return mainCode
        }
    }

    // 3. Fallback padrão
    return "pt"
}