FROM golang:1.25-alpine AS builder

# Instala Node.js (necessário para cssnano/terser) e outras ferramentas úteis
RUN apk add --no-cache nodejs npm

WORKDIR /app

# Inicializa o módulo (se não existir go.mod)
RUN go mod init topstrem 2>/dev/null || true

# Instala a dependência do templ (versão estável)
RUN go get github.com/a-h/templ@v0.2.707

# Instala o CLI do templ
RUN go install github.com/a-h/templ/cmd/templ@v0.2.707

# Copiar go.mod primeiro para cache de dependências  
COPY go.mod go.sum ./  
RUN go mod download  

# Copia todo o código fonte
COPY . .

# Remove qualquer arquivo _templ.go antigo (evita conflitos na geração)
RUN find . -name "*_templ.go" -type f -delete

# Gera os arquivos Go a partir dos templates .templ
RUN templ generate

# Instala ferramentas de minificação globalmente
RUN npm install -g cssnano terser

# Minifica arquivos CSS e JS (ignora erros se não houver arquivos)
RUN find ./cmd/app/assets/static/css -name "*.css" -exec sh -c 'cssnano "$1" > "$1.tmp" && mv "$1.tmp" "$1"' _ {} \; || true
RUN find ./cmd/app/assets/static/js -name "*.js" -exec sh -c 'terser "$1" --compress --mangle -o "$1"' _ {} \; || true

# Atualiza as dependências (inclui o runtime do templ)
RUN go mod tidy

# Compila o binário (ajustado para o caminho correto)
RUN go build -o topstrem ./cmd/server

# Imagem final (leve)
FROM alpine:latest

# Instala certificados CA (útil para HTTPS, opcional)
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/topstrem .
COPY --from=builder /app/cmd/app/assets ./assets

EXPOSE 8080

CMD ["./topstrem"]