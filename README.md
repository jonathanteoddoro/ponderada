# Ponderada — Servidor Go simples com CRUD

Resumo rápido
- Projeto mínimo em Go que expõe APIs HTTP para duas tabelas: `employees` e `products`.
- Usa GORM + driver Postgres para persistência e auto-migração das tabelas.
- Inclui páginas HTML estáticas (mock/cliente) para testar listagem e criação via browser.

O que há no repositório
- `main.go` — servidor HTTP, modelos `Employee` e `Product`, handlers para CRUD (criação e listagem) e middleware CORS.
- `conn/conn.go` — função de conexão com Postgres (usa GORM). Atualmente a URL de conexão está configurada no código.
- `.env` — arquivo de variáveis de ambiente (exemplo) com parâmetros do banco (opcional; atualmente o código não carrega `.env` automaticamente).
- `employees.html` — cliente front-end que consome `/employees` e `/products` (GET/POST).
- `index.html` — página mock para cadastro/listagem de produtos (dados mockados).

Rotas disponíveis
- `GET  /employees`  — retorna JSON com todos os funcionários
- `POST /employees`  — cria funcionário (JSON: `{ "name": "..", "address": ".." }`)
- `GET  /products`   — retorna JSON com todos os produtos
- `POST /products`   — cria produto (JSON: `{ "name": "..", "price": 12.5, "description": "..", "in_stock": true }`)

Observações sobre CORS
- O servidor já possui um middleware que adiciona os cabeçalhos CORS (em desenvolvimento está configurado para `Access-Control-Allow-Origin: *`).
- Se você abrir os arquivos HTML via `file://` ou outro servidor, pode ocorrer bloqueio por CORS; recomenda-se servir os arquivos via HTTP ou usar o mesmo host/porta do backend.

Como rodar localmente
1. Ajuste a string de conexão ao Postgres em `conn/conn.go` ou defina variáveis de ambiente conforme seu fluxo.

2. Instale dependências (opcional — Go modules resolverá automaticamente):

```bash
# na raiz do projeto
go mod tidy
```

3. Rodar o servidor:

```bash
go run main.go
# ou
go build -o ponderada && ./ponderada
```

4. Testar endpoints com curl:

```bash
# Preflight OPTIONS (simula browser)
curl -i -X OPTIONS http://localhost:8080/employees \
  -H "Origin: http://127.0.0.1:5500" \
  -H "Access-Control-Request-Method: POST"

# GET
curl -i -H "Origin: http://127.0.0.1:5500" http://localhost:8080/employees

# POST (criar employee)
curl -i -X POST http://localhost:8080/employees \
  -H "Content-Type: application/json" \
  -d '{"name":"Fulano","address":"Rua A"}'
```

Testar as páginas front-end
- Abra `employees.html` (recomendado servir via servidor estático para evitar políticas de navegador):

```bash
# na raiz do projeto (exemplo simples)
python3 -m http.server 8000
# então abra http://localhost:8000/employees.html
```

Boas práticas e próximos passos sugeridos
- Remover `Access-Control-Allow-Origin: *` em produção e restringir para origens permitidas.
- Mover a configuração do banco para variáveis de ambiente e carregar `.env` com `github.com/joho/godotenv` ou similar.
- Implementar rotas adicionais: GET/PUT/DELETE por ID, paginação, validação e autenticação.
- Servir os arquivos estáticos diretamente pelo servidor Go se preferir não depender de um servidor separado.

Status atual
- Tabelas `Employee` e `Product` são migradas automaticamente ao iniciar.
- Endpoints de criação e listagem para ambas as tabelas estão implementados.
- CORS habilitado para desenvolvimento.

Se quiser, eu:
- configuro a leitura de `.env` automaticamente e atualizo `conn/conn.go`, ou
- implemento rotas GET/PUT/DELETE por ID, ou
- sirvo `employees.html` diretamente pelo Go para evitar CORS.

