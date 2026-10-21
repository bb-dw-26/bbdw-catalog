# Especificação de bbdw-catalog

## Escopo da apresentação de 50 minutos

A apresentação tem **50 minutos** e começa com a aplicação já executando a revisão defeituosa de **CHK-001**. A primeira compra revela a falha; a meta é iniciar Sentry → issue no checkout → workflow nos primeiros **5 minutos**.

Enquanto o agente trabalha em segundo plano, apresentar a arquitetura, o YAML do workflow, o contexto do incidente e os critérios de validação. Retomar os resultados por volta do minuto 30, sem acompanhar logs continuamente. Os cerca de 25 minutos disponíveis não são garantia de conclusão.

**Prioridades:** CHK-001 é o único cenário obrigatório; CAT-001 é reserva opcional, preparada apenas depois do fluxo principal; USR-001, USR-002, CAT-002, CHK-002, WEB-001 e WEB-002 são material complementar. Nenhum cenário adicional bloqueia a entrega.

Manter quatro aplicações pequenas, dados fictícios em memória, HTTP/JSON e CI/CD próprio por repositório. Autenticação real, banco de dados, pagamentos, mensageria, reparos entre repositórios e múltiplos agentes ficam fora desta entrega.

**Papel deste repositório:** Dependência saudável em Go: listar e consultar os produtos fictícios. CAT-001 é reserva opcional; CAT-002 é material complementar. Não é necessário executar um segundo reparo ao vivo para mostrar a diversidade de tecnologias.

[Roteiro completo e contingência](https://github.com/bbdw-gh-action/bbdw-storefront/blob/main/docs/ROTEIRO.md).

## Objetivo e responsabilidade

Fornecer produtos fictícios, preços em centavos, disponibilidade e imagens opcionais.

É dono dos produtos e dos preços. Não mantém carrinhos, pedidos ou reservas de estoque nesta primeira versão.

Tecnologia: Go, inicialmente com net/http. Porta local proposta: `3002`. Clientes e produtos vêm de fixtures em memória. A criação do pedido retorna a confirmação sem persistência ou histórico nesta versão. Persistência durável fica fora desta primeira entrega.

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

O scaffold implementa os endpoints e fluxos obrigatórios. Users e catalog carregam as fixtures versionadas; checkout e interface usam os contratos HTTP. Os endpoints marcados como complementares permanecem fora da implementação.

## Contratos comuns

- HTTP/JSON, IDs como strings, valores monetários em centavos inteiros e moeda BRL.
- Envelope de erro: `{"error":{"code":"CODIGO_ESTAVEL","message":"Mensagem legível"},"requestId":"identificador"}`.
- Propagar `X-Request-ID` nas chamadas e respostas; gerar um identificador quando ausente.
- Dados exclusivamente fictícios. Não incluir tokens, cabeçalhos de autenticação ou dados pessoais reais nos eventos.
- O checkout deve ter projeto Sentry próprio e identificar ambiente, release, serviço e cenário. Os demais serviços também possuem SDK e projetos próprios; consulte docs/SENTRY.md para as evidências.
- Contratos de erro esperados (4xx) não são falhas inesperadas por padrão. Os defeitos descritos devem gerar exceções reais, instrumentadas explicitamente.
- A implementação deve evitar captura duplicada da mesma exceção.

## Observabilidade

Para CHK-001, instrumentar o checkout, publicar source maps associados à release e comprovar o envio de uma exceção real ao Sentry. Propagar o request ID pelos serviços e evitar eventos duplicados. Testes automatizados usam clientes simulados, sem credenciais reais.

Users e catalog devem permanecer saudáveis; a interface mostra o erro retornado pelo checkout e permite repetir a mesma compra após o reparo. A instrumentação Sentry própria desses três repositórios foi autorizada e validada como suporte à demonstração.

Se CAT-001 for preparado como reserva, instrumentar a recuperação do panic HTTP no Go para capturar stack e retornar 500 sem derrubar o serviço. Essa preparação é opcional e ocorre depois da conclusão do cenário principal.

## CI/CD independente

**CI previsto:** Verificação de gofmt, go vet, go test ./..., testes HTTP com httptest e go build ./.... Fixar a versão do Go ao implementar.

**CD previsto:** Binário compilado em imagem própria, health check, configuração, deploy e rollback independentes. DSN próprio somente se a reserva CAT-001 for preparada.

O destino de hospedagem será escolhido posteriormente. O workflow de reparo do checkout recebe o incidente CHK-001 e abre um PR no checkout. Workflows equivalentes nos outros repositórios ficam para uma etapa posterior. Os comandos estão no README e o CI em `.github/workflows/ci.yml`. Dockerfile disponível; ainda não há deploy para hospedagem configurado.

## Critérios de entrega para os 50 minutos

- Dependência saudável em Go: listar e consultar os produtos fictícios. CAT-001 é reserva opcional; CAT-002 é material complementar. Não é necessário executar um segundo reparo ao vivo para mostrar a diversidade de tecnologias.
- Validar os contratos usados pela compra principal e os casos vizinhos necessários à correção.
- Manter build, testes e deploy independentes deste repositório.
- Ensaiar a compra com as quatro versões fixadas e registrar request ID e releases.
- Preparar um resultado anterior de CHK-001, identificado como ensaio, para a contingência.
- Documentar execução local e encerramento dos processos.

Para checkout, a entrega exige CHK-001 reproduzível, evento Sentry, issue, execução do agente, PR com teste falhando antes e passando depois e repetição da compra corrigida. Para os demais repositórios, a entrega exige suportar essa jornada de forma saudável. As soluções e falhas complementares continuam públicas para exploração posterior, sem obrigação de implementação nesta etapa.

## Instrumentação adicional autorizada

O serviço possui SDK e projeto Sentry próprio. A validação e a configuração estão em [SENTRY.md](SENTRY.md). Isso não adiciona novos cenários defeituosos à apresentação: CHK-001 continua sendo o único obrigatório.
