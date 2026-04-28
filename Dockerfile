FROM golang:1.25-alpine AS builder
WORKDIR /app

# Inicializa o módulo
RUN go mod init topstrem 2>/dev/null || true

# Instala a dependência do templ (versão estável)
RUN go get github.com/a-h/templ@v0.2.707

# Instala o CLI do templ (mesma versão)
RUN go install github.com/a-h/templ/cmd/templ@v0.2.707

# Copia todo o código fonte
COPY . .
RUN rm -f cmd/app/*.go

# Remove qualquer arquivo _templ.go antigo (evita conflitos)
RUN find . -name "*_templ.go" -type f -delete

# Gera os arquivos Go a partir dos templates .templ
RUN templ generate

# Atualiza as dependências (inclui o runtime, que agora está no go.mod)
RUN go mod tidy

# Compila o binário
RUN go build -o topstrem ./cmd/server

# Imagem final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/topstrem .
COPY --from=builder /app/cmd/app/assets ./assets
EXPOSE 8080
CMD ["./topstrem"]