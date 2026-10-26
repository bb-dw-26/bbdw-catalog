# bbdw-catalog

Serviço de catálogo em Go. Fornece produtos fictícios, preços em centavos, disponibilidade e imagens opcionais.

## Executar localmente

Pré-requisito: Go na versão indicada em `go.mod`. O catálogo usa net/http e o SDK Sentry; as fixtures são incorporadas no binário.

```sh
go run ./cmd/server
```

Porta padrão: **3002**. `PORT` e `CORS_ORIGIN` são variáveis de ambiente; `.env.example` documenta os valores, mas o Go não carrega esse arquivo automaticamente.

```sh
go vet ./...
go test ./...
go build -o bin/catalog ./cmd/server
```

Encerre o processo com **Ctrl+C** ao terminar. Cada aplicação tem seu próprio lockfile ou módulo Go, testes e workflow em `.github/workflows/ci.yml`. O CI também constrói a imagem; a publicação e o deploy serão configurados quando a hospedagem for escolhida.

Para iniciar as três aplicações juntas, consulte [execução integrada](https://github.com/bbdw-gh-action/bbdw-storefront/blob/main/docs/EXECUCAO.md).

## Documentação

- [Especificação funcional e técnica](docs/SPEC.md)
- [Fixtures sintéticas](specs/fixtures.json)

## Outros repositórios

- [bbdw-checkout](https://github.com/bbdw-gh-action/bbdw-checkout)
- [bbdw-storefront](https://github.com/bbdw-gh-action/bbdw-storefront)

## Monitoramento Sentry

SDK configurado no projeto próprio `bbdw-catalog`, com SDK de captura. Consulte [configuração](docs/SENTRY.md).

## Automação compartilhada

Integração com o Patchwork configurada no workflow local. Consulte [operação, skills e MCP](docs/PATCHWORK.md).
