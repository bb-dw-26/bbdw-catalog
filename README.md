# bbw-catalog

Serviço de catálogo para a demonstração pública de fluxo agêntico no BBW, apresentação do Banco do Brasil.

Fornecer produtos fictícios, preços em centavos, disponibilidade e imagens opcionais.

**Tecnologia prevista:** Go, inicialmente com net/http.

**Estado atual:** especificações iniciais. A aplicação, os defeitos executáveis, a integração Sentry e os workflows de CI/CD e reparo ainda serão implementados.

## Documentação

- [Especificação funcional e técnica](docs/SPEC.md)
- [Cenários de falha e reprodução](docs/FALHAS.md)
- [Causas e soluções de referência](docs/SOLUCOES.md)
- [Critérios e registro de avaliação](docs/AVALIACAO.md)
- [Fixtures sintéticas](specs/fixtures.json)
- [Catálogo estruturado de cenários](specs/scenarios.json)

## Fluxo da demonstração

Requisição com defeito → Sentry → issue neste repositório → GitHub Actions → Copilot ou Codex → PR com correção e teste → revisão → deploy → repetição da requisição.

Cada repositório terá seu próprio CI/CD e sua própria configuração de agente. A escolha entre Copilot e Codex permanece aberta. As soluções públicas permitem comparar o reparo com os critérios esperados; este é um exercício guiado.

## Outros repositórios

- [bbw-users](https://github.com/bbdw-gh-action/bbw-users)
- [bbw-checkout](https://github.com/bbdw-gh-action/bbw-checkout)
- [bbw-storefront](https://github.com/bbdw-gh-action/bbw-storefront)
