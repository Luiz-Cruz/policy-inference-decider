# Policy Inference Decider (PID)

Serviço que avalia políticas descritas em grafo (DOT). Você manda o grafo, o input e recebe o output da inferência (qual
nó foi atingido e quais atributos). Roda como Lambda na AWS.

## O que é

API HTTP (Lambda Function URL) que recebe POST com:

- `policy_dot`: grafo em DOT (ex.: `digraph { start -> ok [cond="age>=18"]; ... }`)
- `input`: mapa de variáveis (ex.: `{"age": 20}`)

Resposta: JSON com o `output` do nó atingido após avaliar as condições das arestas.

A Lambda usa o runtime **provided** (Amazon Linux 2023). O binário no deploy se chama `bootstrap` para o runtime
encontrar.

## Makefile

Na raiz do projeto, `make` + alvo:

| Comando              | O que faz                                           |
|----------------------|-----------------------------------------------------|
| `make test`          | Roda os testes                                      |
| `make coverage`      | Testes com cobertura; falha se &lt; 90%             |
| `make build-lambda`  | Gera o zip para a Lambda (linux/arm64, `bootstrap`) |
| `make format`        | gofmt + gofumpt + go mod tidy                       |
| `make sort-imports`  | Ordena imports (exige `make install-tools` antes)   |
| `make install-tools` | Instala gci e gofumpt                               |
| `make run-all`       | format, sort-imports e test (bom antes de commit)   |

## AWS

- **Lambda** em us-east-1, arquitetura arm64.
- **Logs**: CloudWatch Logs com o nome da função. Os logs da aplicação (slog) vão para o mesmo log group.

## CI/CD (GitHub Actions)

- **Deploy** (`.github/workflows/deploy.yml`): em todo push na `main` o workflow faz build (Go, linux/arm64), empacota
  em zip e atualiza o código da Lambda com `update-function-code`. Usa o ambiente **prod** do repositório; os secrets
  `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` e `AWS_LAMBDA_FUNCTION_NAME` vêm desse ambiente.
- **Coverage** (`.github/workflows/coverage.yml`): em PRs para `main` ou `develop` (quando o PR não está em draft) roda
  os testes com cobertura e exige **≥ 90%** para passar. O check deve ser configurado como obrigatório na proteção de
  branch para bloquear merge sem cobertura mínima.

## Proteção de branch

Em Settings → Branches → branch protection, foi configurado para **main** e **develop**:

- **Exigir pull request**: não permitir push direto; todo código entra via PR.
- **Exigir status checks**: marque o check **Coverage** como obrigatório antes do merge.

Assim ninguém faz push direto em main/develop e ninguém mergeia sem bater os 90% de cobertura.

## Load tests

Testes de carga com Artillery, rodando os workers na AWS (Lambda). Pré-requisitos: Artillery instalado e credenciais AWS
configuradas. Na raiz do projeto:

```bash
npm install -g artillery
./loadtest/run-lambda.sh "https://SUA_LAMBDA_URL" [arquivo.yaml]
```

Exemplos: `test-50rps.yaml`, `test-100rps-mixed.yaml`. O script substitui o placeholder da URL no YAML e chama o
Artillery em modo run-lambda.

## Estrutura

- `main.go`: entrada Lambda, monta handler com parser e executor.
- `internal/handler`: HTTP/Lambda handler, binding do body e chamada ao parser + executor.
- `internal/policy`: parsing do DOT, execução do grafo (condições govaluate, resultado por nó).
- `internal/apierror`: códigos e formatos de erro da API.
- `loadtest/`: scripts e YAMLs do Artillery.
