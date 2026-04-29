# 🔐 Sistema de Autenticação Híbrido - Web + Mobile Nativo

## 📋 Visão Geral

O TopStrem agora suporta **3 cenários de autenticação** com o mesmo código backend:

1. **Web Desktop/Mobile Browser** - Cookie HttpOnly
2. **Mobile Nativo (iOS/Android)** - JWT em header Authorization
3. **Custom Tab Android** - Cookie HttpOnly (compatível com WebView)

---

## 🎯 Fluxo por Tipo de Cliente

### 1️⃣ Web Browser (Desktop + Mobile Chrome)

```mermaid
Web Browser → /auth/login → Redireciona para Google
                              ↓
                        Usuário autentica
                              ↓
                        /auth/callback?code=xxx
                              ↓
                        Servidor troca code por JWT
                        e armazena em HttpOnly cookie
                              ↓
                        Redireciona para /
                              ↓
                        fetchCurrentUser() detecta cookie
```

**Como funciona:**
- ✅ Browser envia automaticamente cookies em requisições (credentials: 'same-origin')
- ✅ Cookie HttpOnly é seguro contra XSS
- ✅ SameSite=Strict protege contra CSRF

---

### 2️⃣ Mobile Nativo (iOS/Android - sem WebView)

```mermaid
App Nativo → /auth/login (com header X-Client-Type: native)
                              ↓
                        Servidor retorna JSON:
                        { "authUrl": "https://accounts.google..." }
                              ↓
                        App abre URL em navegador externo
                              ↓
                        Usuário autentica no Google
                              ↓
                        Google redireciona para:
                        https://topstrem.com/auth/callback?code=xxx&client_type=native
                              ↓
                        Servidor retorna HTML com token
                        (página contém <script> que chama deepLink)
                              ↓
                        Deeplink: topstrem://auth?token=eyJhbGc...
                              ↓
                        App nativo intercepta deeplink
                        e armazena token em secure storage
                              ↓
                        App chama window.onAuthTokenReceived(token)
```

**Como funciona:**
- ✅ App nativo abre navegador externo (Safari, Chrome)
- ✅ Após autenticação, deeplink traz o token de volta
- ✅ App armazena token em Keychain (iOS) ou KeyStore (Android)
- ✅ Requisições subsequentes usam header: `Authorization: Bearer <token>`

---

### 3️⃣ Custom Tab Android (app que usa WebView/CustomTab)

```mermaid
CustomTab → /auth/login → Retorna JSON com URL
                              ↓
                        CustomTab abre URL
                              ↓
                        Usuário autentica
                              ↓
                        /auth/callback?code=xxx
                              ↓
                        Cookie HttpOnly é setado
                        e página redireciona
                              ↓
                        App notifica via AndroidInterface
```

---

## 🔄 Detecção de Tipo de Cliente

**No servidor (`internal/auth/auth.go`):**

```go
func detectClientType(r *http.Request) ClientType {
    // Prioridade 1: Header customizado
    if clientType := r.Header.Get("X-Client-Type"); clientType == "native" {
        return ClientTypeNativeApp
    }
    
    // Prioridade 2: User-Agent
    if strings.Contains(userAgent, "wv") || strings.Contains(userAgent, "WebView") {
        return ClientTypeCustomTab
    }
    
    // Padrão: Web Browser
    return ClientTypeWebBrowser
}
```

**No cliente (`login.js`):**

```javascript
let clientType = {
    isNativeApp: window.isMobileApp === true || !!window.getAuthToken,
    isWebBrowser: true
};
```

---

## 📝 Implementação no App Nativo

### iOS (Swift)

```swift
import SafariServices
import KeychainSwift

class AuthManager {
    let keychain = KeychainSwift()
    let baseURL = "https://topstrem.com"
    
    func login() {
        // 1. Requisitar authUrl
        fetch("\(baseURL)/auth/login", headers: ["X-Client-Type": "native"])
            .then { response in response.json() }
            .then { data in
                // 2. Abrir Safari com authUrl
                let safariVC = SFSafariViewController(url: URL(string: data.authUrl)!)
                safariVC.delegate = self
                present(safariVC)
            }
    }
    
    func handleDeeplink(_ url: URL) {
        // 3. App recebe deeplink: topstrem://auth?token=...
        if url.scheme == "topstrem" && url.host == "auth" {
            let token = URLComponents(url: url, resolvingAgainstBaseURL: true)?
                .queryItems?.first(where: { $0.name == "token" })?.value
            
            if let token = token {
                // 4. Armazenar em Keychain
                keychain.set(token, forKey: "authToken")
                // 5. Notificar que login foi bem-sucedido
                DispatchQueue.main.async {
                    NotificationCenter.default.post(name: NSNotification.Name("AuthSuccess"), object: nil)
                }
            }
        }
    }
    
    // Enviar requisições com token
    func request(_ endpoint: String) {
        var request = URLRequest(url: URL(string: "\(baseURL)\(endpoint)")!)
        
        if let token = keychain.get("authToken") {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        
        URLSession.shared.dataTask(with: request).resume()
    }
}
```

### Android (Kotlin)

```kotlin
import androidx.browser.customtabs.CustomTabsIntent
import androidx.security.crypto.EncryptedSharedPreferences

class AuthManager(context: Context) {
    private val baseURL = "https://topstrem.com"
    private val prefs = EncryptedSharedPreferences.create(
        context, "auth", masterKey, EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
    )
    
    fun login(context: Context) {
        // 1. Requisitar authUrl
        val client = OkHttpClient()
        val request = Request.Builder()
            .url("$baseURL/auth/login")
            .header("X-Client-Type", "native")
            .build()
        
        client.newCall(request).enqueue(object : Callback {
            override fun onResponse(call: Call, response: Response) {
                val json = JSONObject(response.body!!.string())
                val authUrl = json.getString("authUrl")
                
                // 2. Abrir em CustomTab
                CustomTabsIntent.Builder()
                    .build()
                    .launchUrl(context, Uri.parse(authUrl))
            }
            
            override fun onFailure(call: Call, e: IOException) {}
        })
    }
    
    fun handleDeeplink(uri: Uri) {
        // 3. App recebe deeplink: topstrem://auth?token=...
        if (uri.scheme == "topstrem" && uri.host == "auth") {
            val token = uri.getQueryParameter("token")
            
            if (token != null) {
                // 4. Armazenar em secure storage
                prefs.edit().putString("authToken", token).apply()
                // 5. Notificar UI
                EventBus.getDefault().post(AuthSuccessEvent())
            }
        }
    }
    
    // Enviar requisições com token
    fun request(endpoint: String) {
        val token = prefs.getString("authToken", null)
        
        val request = Request.Builder()
            .url("$baseURL$endpoint")
            .apply {
                if (token != null) {
                    header("Authorization", "Bearer $token")
                }
            }
            .build()
        
        OkHttpClient().newCall(request).execute()
    }
}
```

### Flutter (Dart)

```dart
import 'package:uni_links/uni_links.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:url_launcher/url_launcher.dart';

class AuthManager {
  final baseURL = "https://topstrem.com";
  final storage = FlutterSecureStorage();
  
  Future<void> login() async {
    // 1. Requisitar authUrl
    final response = await http.get(
      Uri.parse("$baseURL/auth/login"),
      headers: {"X-Client-Type": "native"}
    );
    
    final data = jsonDecode(response.body);
    final authUrl = data["authUrl"];
    
    // 2. Abrir em navegador externo
    if (await canLaunch(authUrl)) {
      await launch(authUrl, forceSafariVC: false, forceWebView: false);
    }
  }
  
  Future<void> setupDeeplinks() {
    // 3. Monitorar deeplinks: topstrem://auth?token=...
    deepLinkStream.listen((String? link) {
      if (link != null) {
        final uri = Uri.parse(link);
        if (uri.scheme == "topstrem" && uri.host == "auth") {
          final token = uri.queryParameters["token"];
          
          if (token != null) {
            // 4. Armazenar em secure storage
            storage.write(key: "authToken", value: token);
            // 5. Notificar app
            AuthNotifier().notifySuccess();
          }
        }
      }
    });
  }
  
  Future<http.Response> request(String endpoint) async {
    final token = await storage.read(key: "authToken");
    
    return http.get(
      Uri.parse("$baseURL$endpoint"),
      headers: token != null ? {"Authorization": "Bearer $token"} : {}
    );
  }
}
```

---

## 🔗 Deep Link Configuration

### AndroidManifest.xml (Android)

```xml
<activity android:name=".MainActivity">
    <intent-filter android:label="@string/app_name">
        <action android:name="android.intent.action.MAIN" />
        <category android:name="android.intent.category.LAUNCHER" />
    </intent-filter>
    
    <!-- Deep link para topstrem://auth?token=... -->
    <intent-filter android:autoVerify="true">
        <action android:name="android.intent.action.VIEW" />
        <category android:name="android.intent.category.DEFAULT" />
        <category android:name="android.intent.category.BROWSABLE" />
        <data android:scheme="topstrem" android:host="auth" />
    </intent-filter>
</activity>
```

### Info.plist (iOS)

```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>topstrem</string>
        </array>
    </dict>
</array>
```

---

## 🔐 Endpoints de Autenticação

### GET /auth/login

**Detecta tipo de cliente e retorna:**

- **Web Browser**: 302 Redirect para Google OAuth
- **Mobile Nativo**: 
  ```json
  {
    "status": "ok",
    "authUrl": "https://accounts.google.com/o/oauth2/v2/auth?..."
  }
  ```

### GET /auth/callback?code=xxx&client_type=native

**Web Browser:**
- Sets HttpOnly cookie com JWT
- 302 Redirect para /

**Mobile Nativo (client_type=native):**
- Returns HTML com JavaScript que abre deeplink
- Exemplo: `topstrem://auth?token=eyJhbGc...`

### GET /api/me

**Autenticação suportada:**
- Cookie HttpOnly (web)
- Header: `Authorization: Bearer <token>` (mobile nativo)

**Response:**
```json
{
  "email": "user@example.com",
  "name": "User Name"
}
```

### POST /auth/logout

**Autenticação suportada:**
- Cookie HttpOnly (web) - limpa cookie
- Header: `Authorization: Bearer <token>` (mobile) - apenas confirma logout

**Response:**
```json
{
  "status": "logged out"
}
```

---

## 🛡️ Segurança

### ✅ Web Browser
- ✅ Cookie HttpOnly - não acessível via JS
- ✅ Secure flag - apenas HTTPS
- ✅ SameSite=Strict - CSRF protection
- ✅ 72 horas de expiração

### ✅ Mobile Nativo
- ✅ JWT armazenado em secure storage (Keychain/KeyStore)
- ✅ Deeplink apenas entre apps do mesmo bundle
- ✅ Token em header (não em URL)
- ✅ 72 horas de expiração

### ✅ Geral
- ✅ JWT assinado com HMAC-256
- ✅ Validação de signature em toda requisição
- ✅ Algoritmo "none" rejeitado explicitamente
- ✅ ONLY HMAC permitido

---

## 🧪 Testar Localmente

### Web Browser
```bash
curl http://localhost:8080/auth/login
# → 302 Redirect para Google
```

### Mobile Nativo
```bash
curl -H "X-Client-Type: native" http://localhost:8080/auth/login
# → JSON com authUrl
```

### Verificar Autenticação (Web)
```bash
curl -b "auth_token=..." http://localhost:8080/api/me
# → {"email": "...", "name": "..."}
```

### Verificar Autenticação (Mobile)
```bash
curl -H "Authorization: Bearer eyJhbGc..." http://localhost:8080/api/me
# → {"email": "...", "name": "..."}
```

---

## 📱 Migração de WebView para Nativo

Se você tinha código usando a antiga WebView:

**Antes (WebView):**
```javascript
window.AndroidInterface.openCustomTab(authUrl);
// Depois que volta:
window.AndroidInterface.onLoginSuccess();
```

**Agora (Mobile Nativo):**
```javascript
// 1. App envia header X-Client-Type: native ao chamar login
// 2. Abre URL em navegador externo
// 3. Intercepta deeplink topstrem://auth?token=...
// 4. Armazena em secure storage
// 5. login.js detecta automaticamente via window.onAuthTokenReceived(token)
```

---

## 🐛 Troubleshooting

### "Chrome mobile não está autenticando"
- **Causa**: Cookie não estava sendo enviado em redirecionamentos cross-origin
- **Solução**: Agora usa navegador externo + deeplink para mobile nativo

### "Token JWT inválido"
- **Verificar**: JWT_SECRET está configurado corretamente?
- **Verificar**: Token não expirou (72 horas)?

### "Deeplink não abre"
- **Verificar**: AndroidManifest.xml tem a intent-filter correta?
- **Verificar**: Bundle identifier em Info.plist está correto?

---

## 📚 Referências

- [OAuth2 Flow](https://tools.ietf.org/html/rfc6749)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8949)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
