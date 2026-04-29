# Configuração de Autenticação - Webview Android

## 📋 Pré-requisitos

1. **Google Cloud Console registrado**
   - OAuth 2.0 Client ID e Secret
   - Redirect URI configurada: `http://localhost:8080/auth/callback`

2. **Redis rodando localmente**
   ```bash
   # macOS com Homebrew
   brew install redis
   brew services start redis
   
   # Linux
   sudo apt-get install redis-server
   sudo systemctl start redis-server
   
   # Docker
   docker run -d -p 6379:6379 redis:latest
   ```

3. **Go 1.25+** instalado

## 🔧 Configuração do Projeto

### 1. Criar arquivo `.env` na raiz do projeto

```bash
cp .env.example .env
```

### 2. Preencher variáveis de ambiente

```env
GOOGLE_CLIENT_ID=seu_client_id_do_gcp.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=sua_client_secret_do_gcp
REDIRECT_URL=http://localhost:8080/auth/callback
JWT_SECRET=sua_chave_secreta_minimo_32_caracteres_aleatorios
REDIS_ADDR=localhost:6379
```

### 3. Gerar JWT Secret seguro

```bash
# macOS/Linux
openssl rand -base64 32

# Windows PowerShell
[System.Convert]::ToBase64String((1..32 | ForEach-Object {Get-Random -Maximum 256}) -as [byte[]])
```

## 🚀 Executar o Servidor

```bash
# Carregar variáveis de .env (opcional, dependendo do setup)
go run ./cmd/app/main.go
```

O servidor estará disponível em: **http://localhost:8080**

## 🔐 Fluxo de Autenticação no Webview Android

### 1. Cliente (JS) inicia o login

```javascript
// src: cmd/app/assets/static/js/mobile-auth.js
async function handleMobileLogin() {
    const response = await fetch('/auth/login');
    const data = await response.json();
    
    if (data.url) {
        // Abre Custom Tab do Android
        if (window.AndroidInterface) {
            window.AndroidInterface.openCustomTab(data.url);
        }
    }
}
```

### 2. Servidor responde com URL de autenticação

```go
// /auth/login retorna JSON com URL do Google OAuth
{
    "url": "https://accounts.google.com/o/oauth2/v2/auth?...",
    "success": "true"
}
```

### 3. Usuário autentica no Google (em Custom Tab)

- Google valida credenciais
- Redireciona para: `http://localhost:8080/auth/callback?code=...`

### 4. Servidor troca code por token JWT

```go
// /auth/callback:
// 1. Exchange OAuth code por ID Token do Google
// 2. Valida ID Token
// 3. Gera JWT assinado com JWT_SECRET
// 4. Armazena JWT em HttpOnly secure cookie (NÃO na URL!)
// 5. Redireciona para / (sem expor o token)
```

### 5. App nativo recebe token via cookie

```javascript
// callback.html não recebe mais token na URL
// Token está armazenado em cookie HttpOnly
// JavaScript NÃO consegue acessar (por segurança)
// Servidor envia automaticamente em próximas requisições

// Para app nativo, token está disponível via:
fetch('http://localhost:8080/api/watch/...', {
    credentials: 'include'  // Envia cookies automaticamente
})
```

### 6. App nativo valida sessão

```kotlin
// O servidor envia o cookie auth_token automaticamente
// App nativo valida se sessão está ativa via /health ou qualquer endpoint
// Se receber 401, fazer re-login
```

## 🛡️ Segurança para Webview Android

### ✅ O que está implementado:

1. **JWT com HMAC-SHA256** - Tokens assinados e validáveis
2. **Detecção de Webview** - Retorna JSON em vez de HTML
3. **Validação de ID Token** - Verifica assinatura do Google
4. **CORS habilitado** - Aceita requisições da webview
5. **Rate Limiting** - Proteção contra brute force
6. **HTTPS Ready** - Pode usar HTTPS em produção

### 📝 Recomendações adicionais:

1. **No Android, NUNCA armazene tokens em `localStorage`**
   ```kotlin
   // ✅ BOM - SharedPreferences encriptado
   val encryptedSharedPreferences = EncryptedSharedPreferences.create(
       context,
       "secret_shared_prefs",
       MasterKey.Builder(context).setKeyScheme(MasterKey.KeyScheme.AES256_GCM).build(),
       EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
       EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
   )
   encryptedSharedPreferences.edit().putString("auth_token", token).apply()
   ```

2. **Implementar logout seguro**
   ```javascript
   function logout() {
       localStorage.removeItem('auth_token');
       if (window.AndroidInterface && typeof window.AndroidInterface.clearToken === 'function') {
           window.AndroidInterface.clearToken();
       }
       window.location.href = '/';
   }
   ```

3. **Validar token expiração**
   ```javascript
   function isTokenExpired(token) {
       const payload = parseJWT(token);
       const now = Math.floor(Date.now() / 1000);
       return payload.exp < now;
   }
   ```

4. **Implementar refresh token** (opcional)
   - Armazenar refresh token separado
   - Renovar JWT antes de expirar (72h)

## 🧪 Teste Local

### 1. Iniciar servidor
```bash
go run ./cmd/app/main.go
```

### 2. Acessar página inicial
```
http://localhost:8080
```

### 3. Clique em "Login"
- Verá URL do Google OAuth em JSON
- Custom Tab abre (em app Android)
- Após autenticação, retorna a callback.html

### 4. Verificar token
```javascript
// No console do browser
const token = localStorage.getItem('auth_token');
console.log(token);

// Decodificar
function parseJWT(token) {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => 
        '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
    ).join(''));
    return JSON.parse(jsonPayload);
}

parseJWT(token);
// { email: "user@example.com", name: "User", iat: ..., exp: ... }
```

## 🔍 Troubleshooting

### ❌ "Variáveis Google OAuth não configuradas"
- Verifique se `.env` existe e está preenchido
- Verifique se `GOOGLE_CLIENT_ID` e `GOOGLE_CLIENT_SECRET` não estão vazios

### ❌ Erro "GoogleClientID não configurado"
- `JWT_SECRET` pode estar vazio
- Gere uma nova chave: `openssl rand -base64 32`

### ❌ Redis connection refused
- Verifique se Redis está rodando: `redis-cli ping`
- Ajuste `REDIS_ADDR` se em outra porta

### ❌ CORS error no webview
- Verificar se middleware CORS está aplicado
- Testar com `curl -H "Origin: ..." http://localhost:8080`

### ❌ Token não chega ao app Android
- Verifique se `AndroidInterface` existe na webview
- Verifique logs com `adb logcat`
- Teste fallback em `localStorage` primeiro

## 📚 Estrutura de arquivos relacionados

```
cmd/
├── app/
│   └── assets/
│       ├── callback.html          ← Página que recebe token
│       └── static/js/
│           ├── login.js           ← Login normal
│           └── mobile-auth.js     ← Login específico para webview
internal/
├── auth/
│   ├── auth.go                   ← OAuth2 + JWT
│   ├── claims.go                 ← Validação de tokens
│   └── middleware.go             ← Proteção de rotas
mobile/
└── mobile.go                      ← Servidor principal para app
.env.example                       ← Template de variáveis
```

## 🎯 Próximos passos

1. Registrar aplicativo Android no Google Cloud
2. Implementar WebViewClient no Android
3. Adicionar interface `AndroidInterface` com métodos:
   - `onLoginSuccess(token)`
   - `openCustomTab(url)`
   - `clearToken()`
4. Armazenar token de forma segura no Android
5. Implementar refresh de tokens (optional)
