# 🔐 Guia de Segurança - TopStrem

## ⚠️ PROBLEMAS CRÍTICOS CORRIGIDOS

### ✅ Credenciais Expostas
- **Problema**: Credenciais do Google OAuth estavam no `docker-compose.yml`
- **Solução**: Arquivo agora usa `.env.prod` (não versionado)
- **Ação**: Use `.env.example` e `.env.prod.example` como templates

### ✅ JWT Secret Não Validado
- **Problema**: JWT_SECRET podia estar vazio ou fraco
- **Solução**: Validação em `InitAuth()` exige mínimo 32 caracteres
- **Ação**: Configure JWT_SECRET forte em `.env`

### ✅ CORS Aberto Demais
- **Problema**: `Access-Control-Allow-Origin: *` permitia qualquer origem
- **Solução**: CORS agora valida contra lista de origens permitidas
- **Ação**: Configure `ALLOWED_ORIGINS` no `.env`

### ✅ CSRF Não Funcionava
- **Problema**: CSRF token sempre vazio, não protegia nada
- **Solução**: CSRF agora gera/valida tokens com cookies HttpOnly
- **Ação**: Cliente JS precisa enviar `X-CSRF-Token` header

### ✅ Token JWT Exposto na URL
- **Problema**: Token passado como query parameter em callback
- **Solução**: Token agora armazenado em HttpOnly cookie
- **Ação**: Cliente JS lê token do cookie (não da URL)

### ✅ Redis Sem Autenticação
- **Problema**: Qualquer um na rede podia acessar Redis
- **Solução**: `REDIS_PASSWORD` agora obrigatória
- **Ação**: Configure senha forte em `.env`

### ✅ JWT Algorithm "none" Não Rejeitado
- **Problema**: Token poderia ser forjado com algoritmo "none"
- **Solução**: Middleware valida explicitamente algoritmo HS256
- **Ação**: Nenhuma (automático)

### ✅ Memory Leak em RateLimiter
- **Problema**: IPs nunca eram removidos do mapa
- **Solução**: Delete IPs sem requisições ativas
- **Ação**: Nenhuma (automático)

### ✅ Porta Hardcoded
- **Problema**: Porta sempre 8080, sem flexibilidade
- **Solução**: Port parametrizável via `PORT` env
- **Ação**: Configure `PORT` conforme necessário

### ✅ Sem Graceful Shutdown
- **Problema**: Servidor fechava abruptamente com goroutines ativas
- **Solução**: Shutdown gracioso responde a SIGINT/SIGTERM
- **Ação**: Pressione Ctrl+C para shutdown seguro

---

## 📋 CHECKLIST DE SEGURANÇA PRÉ-DEPLOY

### 1. Variáveis de Ambiente
- [ ] Criar `.env` (local) baseado em `.env.example`
- [ ] Criar `.env.prod` (produção) baseado em `.env.prod.example`
- [ ] Verificar que `.env` e `.env.prod` NÃO estão versionados
- [ ] Confirmar que `GOOGLE_CLIENT_ID` foi alterado
- [ ] Confirmar que `GOOGLE_CLIENT_SECRET` foi alterado
- [ ] Confirmar que `JWT_SECRET` tem 32+ caracteres aleatórios
- [ ] Confirmar que `REDIS_PASSWORD` foi definida
- [ ] Configurar `ALLOWED_ORIGINS` com domínios reais

### 2. Redis
- [ ] Redis rodando com autenticação habilitada
- [ ] `REDIS_PASSWORD` configurada e forte
- [ ] Conexão testada: `redis-cli -a <password> ping`
- [ ] Firewall: Redis não acessível publicamente

### 3. Google OAuth
- [ ] Aplicativo registrado em Google Cloud Console
- [ ] OAuth 2.0 Client ID configurado
- [ ] OAuth 2.0 Client Secret configurado
- [ ] Authorized redirect URI = seu `REDIRECT_URL`
- [ ] Verificar que as credenciais são as de PRODUÇÃO

### 4. HTTPS
- [ ] Certificado SSL/TLS válido
- [ ] `ENVIRONMENT=production` no `.env.prod`
- [ ] Cookies com `Secure=true` (automático em produção)
- [ ] Redirecionar HTTP → HTTPS

### 5. CORS
- [ ] `ALLOWED_ORIGINS` configurado com seus domínios
- [ ] Não usar `*` como origem
- [ ] Testar CORS em ambiente local

### 6. CSRF
- [ ] Cliente JS envia `X-CSRF-Token` em POST/PUT/DELETE
- [ ] Server valida token contra cookie
- [ ] Testar com curl: `curl -X POST http://localhost:8080/favorites -H "X-CSRF-Token: invalid" -H "Cookie: csrf_token=..."`

### 7. Tokens JWT
- [ ] Access token com vida curta (15 min)
- [ ] Refresh token com vida longa (7 dias)
- [ ] Tokens em HttpOnly cookies (não localStorage)
- [ ] Testar expiração de token

### 8. Docker/Kubernetes
- [ ] `docker-compose.yml` usa `.env.prod` (não versionado)
- [ ] Secrets do Kubernetes para credenciais sensíveis
- [ ] Liveness probe: `/health` endpoint
- [ ] Readiness probe: Redis ping

### 9. Logging
- [ ] Erros logados sem expor stack traces
- [ ] Senhas/tokens NUNCA nos logs
- [ ] Logs centralizados (CloudWatch, ELK, Loki)
- [ ] Rotação de logs configurada

### 10. Rate Limiting
- [ ] Rate limit ativado (100 req/min geral)
- [ ] Rate limit por API (30 req/min)
- [ ] Testar com ab ou vegeta

### 11. Testes
- [ ] [ ] Testes unitários para auth criados
- [ ] [ ] Testes de segurança CORS/CSRF
- [ ] [ ] Penetration testing básico realizado

### 12. Monitoramento
- [ ] Alertas para erros críticos
- [ ] Alertas para taxa alta de 401/403
- [ ] Alertas para memory usage
- [ ] Alertas para taxa de erro > 1%

---

## 🔒 Recomendações Adicionais

### Implementar Refresh Tokens
O código atual usa JWT com 72 horas. Adicione:
```go
// Gerar tokens com vida curta
accessToken := generateAccessToken(email, name)    // 15 min
refreshToken := generateRefreshToken(email, name)  // 7 dias

// Endpoint para renovação
POST /auth/refresh
```

### Rate Limiting por User
Adicione limite por usuário autenticado (mais generoso que anônimo):
```go
// Anônimo: 10 req/min
// Autenticado: 100 req/min
```

### Logging Estruturado
Adicione biblioteca como `zap`:
```go
logger.Info("user login",
    zap.String("email", user.Email),
    zap.Time("timestamp", time.Now()),
)
```

### WAF (Web Application Firewall)
Em produção, considere:
- CloudFlare WAF
- AWS WAF
- Nginx ModSecurity

### Secrets Management
Para produção:
- AWS Secrets Manager
- Kubernetes Secrets
- Vault
- Google Secret Manager

---

## 🚨 Se Algo Foi Comprometido

1. **Credenciais Google Expostas**:
   ```bash
   # Invalidar OAuth credentials em https://console.cloud.google.com
   # Gerar novos Client ID/Secret
   # Atualizar .env
   # Fazer redeploy
   ```

2. **JWT_SECRET Comprometido**:
   ```bash
   # Invalidar todos os JWTs forçando re-login
   # Gerar novo JWT_SECRET
   # Limpar cache Redis
   ```

3. **REDIS_PASSWORD Comprometida**:
   ```bash
   # Trocar senha no Redis
   # Limpar dados antigos
   # Atualizar .env
   ```

---

## 📞 Contatos de Segurança

Se encontrar vulnerabilidade:
1. NÃO publique em issues públicas
2. Envie email para: [seu-email-seguranca@domain.com]
3. Inclua: descrição, passos para reproduzir, impacto

---

## 📚 Referências

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OAuth 2.0 Security](https://tools.ietf.org/html/rfc6749)
- [JWT Best Practices](https://tools.ietf.org/html/rfc7519)
- [CORS Spec](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [CSP Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)

---

Última atualização: 28 de Abril de 2026
