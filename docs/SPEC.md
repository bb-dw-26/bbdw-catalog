# Especificação de bbdw-catalog

## Escopo técnico

Aplicação com dados fictícios em memória, contratos HTTP/JSON, build e CI próprios. Autenticação real, banco de dados, pagamentos e mensageria ficam fora do escopo.

## Objetivo e responsabilidade

Fornecer produtos fictícios, preços em centavos, disponibilidade e imagens opcionais.

É dono dos produtos e dos preços. Não mantém carrinhos, pedidos ou reservas de estoque nesta primeira versão.

Tecnologia: Go com net/http. Porta local padrão: `3002`. Os produtos de `specs/fixtures.json` são incorporados ao binário e consultados em memória.

## Interface e regras de negócio

| Método e rota | Resposta esperada |
| --- | --- |
| `GET /health` | 200, `{"status":"ok"}` |
| `GET /products` | 200, `{"items":[...]}` |
| `GET /products/:id` | 200, produto; 404, `PRODUCT_NOT_FOUND` |

Produto de exemplo:

```json
{"id":"product-no-image","name":"Caderno Demo","priceCents":2500,"currency":"BRL","available":true,"imageUrls":[],"thumbnailUrl":null}
```

`thumbnailUrl` é a primeira imagem quando houver imagens; caso contrário, é `null`. Os preços são inteiros não negativos. Disponibilidade é um booleano; o serviço não implementa controle concorrente de estoque.

## Contratos comuns

- HTTP/JSON, IDs como strings, valores monetários em centavos inteiros e moeda BRL.
- Envelope de erro: `{"error":{"code":"CODIGO_ESTAVEL","message":"Mensagem legível"},"requestId":"identificador"}`.
- Propagar `X-Request-ID` nas chamadas e respostas; gerar um identificador quando ausente.
- Dados exclusivamente fictícios. Não incluir tokens, cabeçalhos de autenticação ou dados pessoais reais nos eventos.
- Cada serviço possui projeto Sentry próprio e identifica ambiente, release e serviço. Consulte [SENTRY.md](SENTRY.md).
- Contratos de erro esperados (4xx) não são falhas inesperadas por padrão. Exceções inesperadas devem ser instrumentadas explicitamente.
- A implementação deve evitar captura duplicada da mesma exceção.

## Observabilidade

Cada serviço usa projeto Sentry próprio. Identificar ambiente, release e serviço; propagar o request ID e evitar eventos duplicados. Publicar source maps associados à mesma release nos projetos TypeScript. Testes de telemetria usam coletores locais, sem credenciais reais.

## Build e CI

A versão do Go vem de `go.mod`. O CI verifica gofmt, executa `go vet ./...`, `go test -race ./...` e `go build ./cmd/server`. O workflow [ci.yml](../.github/workflows/ci.yml) também verifica o build da imagem Docker.

Os comandos locais estão no [README](../README.md). A imagem é construída pelo [Dockerfile](../Dockerfile); não há workflow de deploy remoto configurado neste repositório.

## Automação compartilhada

Integração com o Patchwork configurada no workflow local. Consulte [operação, skills e MCP](PATCHWORK.md).
