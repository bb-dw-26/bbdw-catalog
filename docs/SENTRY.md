# Sentry — bbdw-catalog

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

O binário contém arquivo, função e linha da stack. O evento local validado não inclui trechos do código Go; para contexto remoto de código, configure os [code mappings do Sentry](https://docs.sentry.io/integrations/source-code-mgmt/source-context/) associados à revisão efetivamente publicada. Confira o vínculo GitHub do projeto antes de configurar o contexto remoto.
