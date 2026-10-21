# bbdw-catalog

Serviço de catálogo para a demonstração pública de fluxo agêntico no BBDW, apresentação do Banco do Brasil.

Fornecer produtos fictícios, preços em centavos, disponibilidade e imagens opcionais.

**Tecnologia prevista:** Go, inicialmente com net/http.

**Estado atual:** aplicação funcional e saudável, com testes, Dockerfile e CI próprio. O checkout contém CHK-001 isolado, instrumentação Sentry e workflow Codex. A conexão Sentry e a captura neste serviço foram validadas; o ensaio real até o PR continua pendente; o deploy remoto será definido depois.

## Escopo da apresentação de 50 minutos

A apresentação tem **50 minutos** e começa com a aplicação já executando a revisão defeituosa de **CHK-001**. A primeira compra revela a falha; a meta é iniciar Sentry → issue no checkout → workflow nos primeiros **5 minutos**.

Enquanto o agente trabalha em segundo plano, apresentar a arquitetura, o YAML do workflow, o contexto do incidente e os critérios de validação. Retomar os resultados por volta do minuto 30, sem acompanhar logs continuamente. Os cerca de 25 minutos disponíveis não são garantia de conclusão.

**Prioridades:** CHK-001 é o único cenário obrigatório; CAT-001 é reserva opcional, preparada apenas depois do fluxo principal; USR-001, USR-002, CAT-002, CHK-002, WEB-001 e WEB-002 são material complementar. Nenhum cenário adicional bloqueia a entrega.

Manter quatro aplicações pequenas, dados fictícios em memória, HTTP/JSON e CI/CD próprio por repositório. Autenticação real, banco de dados, pagamentos, mensageria, reparos entre repositórios e múltiplos agentes ficam fora desta entrega.

**Papel deste repositório:** Dependência saudável em Go: listar e consultar os produtos fictícios. CAT-001 é reserva opcional; CAT-002 é material complementar. Não é necessário executar um segundo reparo ao vivo para mostrar a diversidade de tecnologias.

[Roteiro completo e contingência](https://github.com/bbdw-gh-action/bbdw-storefront/blob/main/docs/ROTEIRO.md).

## Executar localmente

Pré-requisito: Go na versão indicada em `go.mod`. O catálogo usa apenas a biblioteca padrão; as fixtures são incorporadas no binário.

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

Para iniciar as quatro aplicações juntas, consulte [execução integrada](https://github.com/bbdw-gh-action/bbdw-storefront/blob/main/docs/EXECUCAO.md).

## Documentação

- [Especificação funcional e técnica](docs/SPEC.md)
- [Cenários de falha e reprodução](docs/FALHAS.md)
- [Causas e soluções de referência](docs/SOLUCOES.md)
- [Critérios e registro de avaliação](docs/AVALIACAO.md)
- [Fixtures sintéticas](specs/fixtures.json)
- [Catálogo estruturado de cenários](specs/scenarios.json)

## Fluxo da demonstração

Compra com CHK-001 → Sentry → issue em bbdw-checkout → GitHub Actions → Codex → PR com correção e teste → revisão → deploy → repetição da requisição.

Cada repositório terá seu próprio CI/CD. O workflow de reparo obrigatório nesta entrega é o do checkout; a expansão para outros repositórios é posterior. Codex é a opção inicial do checkout. As soluções públicas permitem comparar o reparo com os critérios esperados; este é um exercício guiado.

## Outros repositórios

- [bbdw-users](https://github.com/bbdw-gh-action/bbdw-users)
- [bbdw-checkout](https://github.com/bbdw-gh-action/bbdw-checkout)
- [bbdw-storefront](https://github.com/bbdw-gh-action/bbdw-storefront)

## Monitoramento Sentry

SDK configurado no projeto próprio `bbdw-catalog`, com captura real validada. Consulte [configuração e evidências](docs/SENTRY.md). A criação automática de issues e o reparo por agente continuam no escopo principal do checkout.
