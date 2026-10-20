# Especificação de bbw-catalog

## Objetivo e responsabilidade

Fornecer produtos fictícios, preços em centavos, disponibilidade e imagens opcionais.

É dono dos produtos e dos preços. Não mantém carrinhos, pedidos ou reservas de estoque nesta primeira versão.

Tecnologia: Go, inicialmente com net/http. Porta local proposta: `3002`. O armazenamento inicial, quando necessário, é em memória, recriado a partir de fixtures em cada inicialização. Persistência durável fica fora desta primeira entrega.

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

`thumbnailUrl` é a primeira imagem quando houver imagens; caso contrário, é `null`. Os preços são inteiros não negativos. Disponibilidade é um booleano; a demo não implementa controle concorrente de estoque.

As fixtures em `specs/fixtures.json` são dados de entrada para implementação e testes, não endpoints já disponíveis.

## Contratos comuns

- HTTP/JSON, IDs como strings, valores monetários em centavos inteiros e moeda BRL.
- Envelope de erro: `{"error":{"code":"CODIGO_ESTAVEL","message":"Mensagem legível"},"requestId":"identificador"}`.
- Propagar `X-Request-ID` nas chamadas e respostas; gerar um identificador quando ausente.
- Dados exclusivamente fictícios. Não incluir tokens, cabeçalhos de autenticação ou dados pessoais reais nos eventos.
- Cada aplicação usa seu próprio projeto Sentry e identifica ambiente, release, serviço e cenário de ensaio.
- Contratos de erro esperados (4xx) não são falhas inesperadas por padrão. Os defeitos descritos devem gerar exceções reais, instrumentadas explicitamente.
- A implementação deve evitar captura duplicada da mesma exceção.

## Observabilidade

Inicializar o SDK apropriado antes da aplicação. Em TypeScript, publicar source maps associados à release. Em Go, capturar panics no limite HTTP, registrar a stack e retornar 500 sem derrubar o serviço; não presumir que a recuperação padrão de net/http envia eventos ao Sentry. No React, conectar o limite de erro ao SDK.

Correlacionar a requisição com a release e o cenário. Verificar a entrega do evento em ensaio real, aguardando a ingestão. Testes unitários não devem depender de Sentry ou credenciais reais.


## CI/CD independente

**CI previsto:** Verificação de gofmt, go vet, go test ./..., testes HTTP com httptest e go build ./.... Fixar a versão do Go ao implementar.

**CD previsto:** Binário compilado em imagem própria, health check, configuração e DSN próprios, deploy e rollback independentes.

O destino de hospedagem será escolhido posteriormente. O workflow de reparo recebe contexto do incidente, trabalha neste repositório e abre um PR. Os comandos de instalação e verificação serão adicionados junto à implementação; não há workflow de deploy executável nesta etapa.

## Critérios de entrega

- Implementar o caminho saudável de acordo com este contrato.
- Criar testes dos comportamentos esperados, incluindo entradas inválidas.
- Instrumentar falhas inesperadas e comprovar captura em um ambiente de demonstração.
- Implementar cenários isolados conforme `FALHAS.md`, preservando revisões reproduzíveis.
- Executar CI com ferramentas próprias desta tecnologia.
- Documentar execução local, configuração e encerramento de processos.
