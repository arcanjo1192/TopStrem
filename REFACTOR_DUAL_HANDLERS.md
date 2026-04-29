# Refatoração: Handlers Duplos (HTML + JSON)

## 📋 Visão Geral

Esta refatoração permite que **cada handler sirva tanto HTML quanto JSON** do mesmo código, sem duplicação. O formato é negociado via header `Accept`:

- **`Accept: application/json`** → JSON puro para mobile/APIs
- **Padrão (sem header ou `Accept: text/html`)** → HTML com template para web

## 🎯 Estratégia

### 1. **Helper de Negociação** (`format.go`)
```go
NegotiateFormat(r *http.Request) ResponseFormat
IsJSONRequest(r *http.Request) bool
```

Detecta automaticamente o formato desejado baseado no header `Accept`.

### 2. **Funções de Dados Internas**
Cada handler tem uma função `get<Nome>Data()` que extrai a lógica pura:

- `getCatalogData()` - dados brutos do catálogo
- `getDetailData()` - dados enriquecidos de detalhe  
- `getEpisodesData()` - episódios agrupados por temporada
- `getWatchData()` - streams disponíveis

**Vantagem**: A lógica é executada uma vez, independente do formato de saída.

### 3. **Resposta Unificada**
Cada handler retorna dados em ambos formatos usando a mesma fonte:

```go
if IsJSONRequest(r) {
    // JSON para mobile/frontend
    json.NewEncoder(w).Encode(data)
} else {
    // HTML com template para web
    templates.PageName(data, lang).Render(r.Context(), w)
}
```

## 📱 Como Usar no Mobile Nativo

### **Requisição JSON**
```bash
# Catálogo (JSON)
curl -H "Accept: application/json" https://topstrem.com/catalog/movie/top

# Resposta
{
  "type": "movie",
  "id": "top",
  "metas": [...]
}
```

### **Requisição HTML** (padrão, para web)
```bash
# Sem header Accept
curl https://topstrem.com/catalog/movie/top

# Resposta: HTML renderizado
<!DOCTYPE html>
<html>...</html>
```

## 🔄 Handlers Refatorados

| Handler | Rota | Dados | Formatos |
|---------|------|-------|----------|
| **CatalogHandler** | `/catalog/{type}/{id}` | `CatalogDataResponse` | HTML ✓ JSON ✓ |
| **DetailHandler** | `/detail/{type}/{id}` | `DetailDataResponse` | HTML ✓ JSON ✓ |
| **EpisodesHandler** | `/api/episodes/{id}` | `EpisodesDataResponse` | HTML ✓ JSON ✓ |
| **WatchHandler** | `/api/watch/{type}/{id}` | `WatchDataResponse` | HTML ✓ JSON ✓ |

## 📦 Estruturas de Resposta JSON

### CatalogDataResponse
```go
{
  "type": "movie|series",
  "id": "catalog-id",
  "metas": [{ id, name, type, poster, ... }]
}
```

### DetailDataResponse
```go
{
  "mediaType": "movie|series",
  "id": "tt1234567",
  "meta": { id, name, description, poster, videos, trailers, ... }
}
```

### EpisodesDataResponse
```go
{
  "seriesId": "tt1234567",
  "seasons": [
    {
      "season": 1,
      "episodes": [{ season, episode, name, description, ... }]
    }
  ]
}
```

### WatchDataResponse
```go
{
  "mediaType": "movie|series",
  "id": "tt1234567",
  "streams": [{ name, url, ... }]
}
```

## 🚀 Exemplos de Requisição (Frontend Mobile)

### React Native / Flutter
```javascript
// Requisitar JSON
const response = await fetch(
  'https://topstrem.com/detail/movie/tt0133093',
  { headers: { 'Accept': 'application/json' } }
);
const data = await response.json();
console.log(data.meta); // Dados estruturados para renderizar UI
```

### Mantém Compatibilidade Web
```html
<!-- Requisição padrão em navegador: obtém HTML renderizado -->
<iframe src="https://topstrem.com/detail/movie/tt0133093"></iframe>
```

## ✅ Benefícios

✓ **Zero duplicação de código** - Lógica executada uma vez  
✓ **Reutilização máxima** - Mesma função para HTML e JSON  
✓ **Negociação automática** - Header `Accept` define o formato  
✓ **Compatível com web** - Padrão continua sendo HTML  
✓ **Escalável** - Adicionar novos handlers é trivial  

## 🔧 Adicionando Novo Handler

1. Criar função `get<Nome>Data()` com a lógica
2. Criar struct `<Nome>DataResponse` para JSON
3. Criar template `<Nome>Page()` (se necessário)
4. No handler, usar:
```go
if IsJSONRequest(r) {
    json.NewEncoder(w).Encode(DataResponse{...})
} else {
    templates.Page(...).Render(r.Context(), w)
}
```

## 📝 Rotas no main.go

Nenhuma mudança necessária! As rotas funcionam exatamente como antes:

```go
http.HandleFunc("/catalog/", middleware.CORS(handlers.CatalogHandler(cachedApiClient)))
http.HandleFunc("/detail/", middleware.CORS(handlers.DetailHandler(cachedApiClient, cachedTmdbClient)))
http.HandleFunc("/api/episodes/", middleware.CORS(handlers.EpisodesHandler(cachedApiClient, cachedTmdbClient)))
http.HandleFunc("/api/watch/", middleware.CORS(handlers.WatchHandler(cachedWatchClient)))
```

A negociação de formato é **automática e transparente** para as rotas! 🎉
