# Sentry — bbdw-catalog

Organização atual: `bbdw-workshop-jj`. Projeto: `bbdw-catalog`. Repositório vinculado: `bbdw-gh-action/bbdw-catalog`.

SDK `sentry-go` inicializado no servidor. Cada requisição recebe um hub isolado. Um panic gera um evento e resposta JSON 500; o processo continua atendendo. Respostas 404 não geram eventos.

A configuração mantém apenas erros nesta etapa, sem replay, profiling ou logs. Requests, usuários e breadcrumbs são removidos antes do envio; os testes verificam a ausência de cabeçalhos, query e dados sintéticos privados.

## Configuração

Copie `.env.example` para `.env` e configure o DSN do projeto próprio, ambiente, release e identificador de execução. Nunca grave `SENTRY_AUTH_TOKEN` no repositório. Sem DSN, o envio fica desativado.

Use `SENTRY_DSN`, `SENTRY_ENVIRONMENT`, `SENTRY_RELEASE`, `DEMO_SCENARIO` e `DEMO_RUN_ID` no ambiente do processo. Go não carrega `.env` automaticamente: exporte as variáveis no terminal ou use o Compose do storefront. Exemplo PowerShell para carregar seu `.env` local:

```powershell
Get-Content .env | ForEach-Object {
  if ($_ -match '^([A-Z_][A-Z0-9_]*)=(.*)$') {
    [Environment]::SetEnvironmentVariable($matches[1], $matches[2], 'Process')
  }
}
go run ./cmd/server
```

O binário contém arquivo, função e linha da stack. O evento local validado não inclui trechos do código Go; para contexto remoto de código, configure os [code mappings do Sentry](https://docs.sentry.io/integrations/source-code-mgmt/source-context/) associados à revisão efetivamente publicada. O vínculo GitHub do projeto já existe.

## Validação histórica na organização anterior — 21/10/2026

- Evento: [*errors.errorString: BBDW SDK validation catalog sdk-20261021-110837](https://bbdw-workshop.sentry.io/issues/7804712243/).
- ID: `9912388552844b64ae1d00f5bc73f4e0`.
- Release exclusiva de teste: `bbdw-catalog@sdk-20261021-110837`; ambiente `bbdw-sdk-validation`, cenário `SDK-SMOKE`.
- Um evento observado no projeto correto, sem erros de processamento do Sentry.
- Stack: `catalog.go:94`; `telemetry.go:44`; `catalog.go:54`.
- A falha foi introduzida somente em uma cópia temporária da aplicação e percorreu a inicialização real e o tratamento de erros. Nenhuma rota ou falha de teste foi adicionada à aplicação versionada.
- Os processos temporários foram encerrados. Este registro não representa deploy remoto nem execução do agente de reparo.

## Checks

`go vet ./...` e `go test ./...`: contratos HTTP e captura concorrente de panics. O teste usa um coletor HTTP local, confirma dois IDs isolados e nenhum evento para 404. `go test -race ./...` permanece no CI Linux; o ambiente Windows local não tem compilador C para o detector de corridas.

## Migração validada — 21/10/2026

Organização `bbdw-workshop-jj`, projeto próprio e repositório homônimo conectados. O `.env` local usa o novo DSN. A organização antiga e seus eventos foram preservados.

- Evento de migração (issue 7804862846 posteriormente excluída a pedido do responsável); ID `d9fdc75018bf435c98158d2fe91194c5`.
- Release: `bbw-catalog@migration-jj-20261021-122004`.
- Ambiente histórico: `bbw-migration-validation`; cenário `SDK-SMOKE`; execução `migration-jj-20261021-122004`.
- Um evento no projeto correto, sem erros de processamento.
- Panic HTTP recuperado com resposta 500, stack com arquivo, função e linha; build de validação com `-trimpath`. Frames Go classificados como externos pelo SDK nesta execução, sem contexto de código inline.
- Mapeamento: `github.com/bbdw-gh-action/bbdw-catalog/` → raiz do repositório, branch `main`.

A exceção foi gerada apenas em uma cópia isolada da aplicação, sem modificar o código saudável versionado. Todos os processos temporários foram encerrados. Os processos que já existiam antes da migração não foram reiniciados; precisam recarregar o ambiente para usar os novos DSNs.

## Nomes corrigidos para BBDW

Projeto atual: `bbdw-catalog`, vinculado a `bbdw-gh-action/bbdw-catalog`. O ID do projeto e o DSN foram preservados. Pacote/serviço, Compose, CI e documentação usam o novo nome. Os identificadores de releases acima pertencem ao ensaio anterior à renomeação. As instâncias locais antigas foram encerradas durante essa correção.
