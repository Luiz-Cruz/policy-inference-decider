# Policy Inference Decider (PID)

Serviço que avalia políticas descritas em grafo (DOT). Você manda o grafo, o input e recebe o output da inferência (qual
nó foi atingido e quais atributos). Roda como Lambda na AWS.

<img width="725" height="369" alt="image" src="https://github.com/user-attachments/assets/33a05066-a68f-432f-87de-0adfda66a800" />


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

<img width="2946" height="890" alt="image (16)" src="https://github.com/user-attachments/assets/ff519933-fd6e-4e72-964c-4c9eedf7a08d" />
<img width="2958" height="1096" alt="image (18)" src="https://github.com/user-attachments/assets/b8f53d2e-cadf-43f9-84fb-309eec19a546" />
<img width="2930" height="432" alt="image (17)" src="https://github.com/user-attachments/assets/b08c3302-7f10-41b4-8faf-b300ac2fb5e2" />
  
- **Coverage** (`.github/workflows/coverage.yml`): em PRs para `main` ou `develop` (quando o PR não está em draft) roda
  os testes com cobertura e exige **≥ 90%** para passar. O check deve ser configurado como obrigatório na proteção de
  branch para bloquear merge sem cobertura mínima.
  
<img width="1702" height="554" alt="image (15)" src="https://github.com/user-attachments/assets/4787d192-d3ce-42f5-8e18-756a5204e463" />


## Proteção de branch

Em Settings → Branches → branch protection, foi configurado para **main** e **develop**:

- **Exigir pull request**: não permitir push direto; todo código entra via PR.
- **Exigir status checks**: foi selecionado o check **Coverage** como obrigatório antes do merge.

Assim ninguém faz push direto em main/develop e ninguém mergeia sem bater os 90% de cobertura.

<img width="2376" height="994" alt="image (20)" src="https://github.com/user-attachments/assets/e3df7269-c142-450c-b63b-5b246719ea7a" />
<img width="2406" height="822" alt="image (21)" src="https://github.com/user-attachments/assets/06eeb90a-20d8-49c3-b1ec-4096c64e10b1" />


## Load tests

Testes de carga com Artillery, rodando os workers na AWS (Lambda). Pré-requisitos: Artillery instalado e credenciais AWS
configuradas. Na raiz do projeto:

```bash
npm install -g artillery
./loadtest/run-lambda.sh "https://SUA_LAMBDA_URL" [arquivo.yaml]
```

Exemplos: `test-50rps.yaml`, `test-100rps-mixed.yaml`. O script substitui o placeholder da URL no YAML e chama o
Artillery em modo run-lambda.

<img width="1384" height="1624" alt="image (19)" src="https://github.com/user-attachments/assets/4b1d7580-b6e1-46e0-8431-5c71d03cd6b3" />


## Next steps (acho interessante agregar no futuro)

- **Métricas custom**: concluir o trabalho em andamento na branch `feature/add-metrics` e montar um dashboard no CloudWatch com as métricas da Lambda (sucesso/erro por causa).
- **Alertas**: configurar alarmes no CloudWatch quando o error rate ultrapassar um limite (ex.: &gt; 1% por pelo menos 5 minutos), para reagir a picos de erro na Lambda.
-- **Monitoramento**: configurar monitoramento da Lambda no CloudWatch (ex.: CPU usage, memory usage, error rate, etc.).
-- **Github Actions**: 
- configurar Github Actions para rodar os testes de carga e cobertura automaticamente.
- verificar possibilidade de segregar os deploys para ambientes de teste desde o CLI AWS ou usando o Github Actions.

## Estrutura

- `main.go`: entrada Lambda, monta handler com parser e executor.
- `internal/handler`: HTTP/Lambda handler, binding do body e chamada ao parser + executor.
- `internal/policy`: parsing do DOT, execução do grafo (condições govaluate, resultado por nó).
- `internal/apierror`: códigos e formatos de erro da API.
- `loadtest/`: scripts e YAMLs do Artillery.
