---
name: service-regression
description: Reproduzir e validar correções de incidentes em bbdw-catalog usando os contratos e testes específicos deste serviço.
---

Leia `docs/SPEC.md` e `docs/FALHAS.md`. O escopo autorizado de produção é `catalog.go`; adicione testes novos que correspondam a `^patchwork_[a-z0-9_]+_test\.go$`. A configuração confiável do Patchwork fica em `.github/patchwork.json`.

Produtos sem imagem retornam imageUrls=[] e thumbnailUrl=null; produto inexistente retorna 404 PRODUCT_NOT_FOUND. Use httptest sem dependências reais. Novos testes devem usar funções TestPatchwork*; execute go test ./... -run TestPatchwork -count=1 antes e depois.

Comandos de verificação: gofmt nos arquivos alterados; go vet ./...; go test -race ./...; go build ./.... Execute primeiro o teste novo contra o código original e registre a falha funcional. Depois aplique a correção e execute as verificações completas. Não substitua testes por alterações em fixtures ou na telemetria.

Use o MCP `sentry` para o projeto `bbdw-catalog` da organização `bbdw-workshop-jj`; consulte o evento e a release informados no prompt. A skill compartilhada `patchwork-investigation` é disponibilizada pelo workflow junto com esta pasta.

As falhas planejadas na documentação não autorizam introduzir defeitos na aplicação. Corrija somente o incidente recebido. Relate se consultou `docs/SOLUCOES.md` e encerre processos de teste iniciados.
