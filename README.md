# Policy Inference Decider (PID)

O **PID** é um serviço projetado para avaliar políticas de decisão descritas em grafos (**DOT**). Ele processa um grafo de estados, aplica os inputs fornecidos e retorna o nó de destino alcançado com seus respectivos atributos.

O serviço é executado como uma **AWS Lambda** em arquitetura **arm64**.
<img width="11346" height="5321" alt="image" src="https://github.com/user-attachments/assets/9856da19-b147-496e-ad69-1f89fa66ed69" />

## 🚀 Como Funciona

A comunicação ocorre via **API HTTP (Lambda Function URL)** através de um método `POST`.

### Estrutura do Payload

* **`policy_dot`**: O grafo no formato DOT (ex: `digraph { start -> ok [cond="age>=18"]; }`).
* **`input`**: Um mapa de variáveis para validação (ex: `{"age": 20}`).

**Resposta:** Um JSON contendo o `output` do nó atingido após a avaliação das condições nas arestas.

---

## 🛠 Desenvolvimento e Makefile

Utilizei o `Makefile` na raiz do projeto para padronizar o fluxo de trabalho:

| Comando | Descrição |
| --- | --- |
| `make test` | Executa a suíte de testes unitários. |
| `make coverage` | Valida a cobertura de testes (falha se for `< 90%`). |
| `make build-lambda` | Compila o binário `bootstrap` e gera o `.zip` para deploy (linux/arm64). |
| `make format` | Executa `gofmt`, `gofumpt` e `go mod tidy`. |
| `make sort-imports` | Organiza os imports (requer `make install-tools`). |
| `make run-all` | Executa formatação, ordenação e testes (ideal para pre-commit). |

---

## ☁️ Infraestrutura e CI/CD

### AWS & Monitoramento

* **Deployment:** Localizado em `us-east-1`**.
* **Logs:** Centralizados no **CloudWatch Logs** via `slog`, integrados ao log group da função.

### Pipelines (GitHub Actions)

1. **Continuous Deployment (`deploy.yml`)**: Disparado em todo push na `main`. Realiza o build em Go, gera o artefato e atualiza a Lambda via `update-function-code` utilizando secrets do ambiente de **prod**.
<img width="2946" height="890" alt="image (16)" src="https://github.com/user-attachments/assets/23e00d60-4608-4a65-9839-0640470f86b0" />

<img width="2958" height="1096" alt="image (18)" src="https://github.com/user-attachments/assets/a6e1efb4-71a2-439b-854f-08ad38d0bf85" />

<img width="2930" height="432" alt="image (17)" src="https://github.com/user-attachments/assets/0ca1a140-9719-450e-bee5-854832ee10f8" />


2. **Quality Gate (`coverage.yml`)**: Disparado em PRs para `main` ou `develop`. Bloqueia o merge caso a cobertura de código seja inferior a **90%**.

<img width="1702" height="554" alt="image (15)" src="https://github.com/user-attachments/assets/ac65bbea-bffc-4910-9b1c-a9bb69b5f13c" />


### Proteção de Branch

As branches `main` e `develop` possuem regras de proteção ativas:

* Pull Requests obrigatórios (proibido push direto).
* Status checks de **Coverage** obrigatórios para merge.
<img width="2376" height="994" alt="image (20)" src="https://github.com/user-attachments/assets/67305a2e-684d-4a34-8db5-ec7401ddc5db" />
<img width="2406" height="822" alt="image (21)" src="https://github.com/user-attachments/assets/aa332168-8d12-4289-a636-fe8c683bc47e" />

---

## 🧪 Testes de Carga

O projeto utiliza o **Artillery** com execução distribuída via AWS Lambda.

**Pré-requisitos:** Artillery instalado globalmente e credenciais AWS configuradas.

```bash
# Instalação
npm install -g artillery

# Execução
./loadtest/run-lambda.sh "https://SUA_LAMBDA_URL" [arquivo.yaml]

```

*Sugestões de cenários em `./loadtest`: `test-50rps.yaml`, `test-100rps-mixed.yaml`.*

<img width="1384" height="1624" alt="image (19)" src="https://github.com/user-attachments/assets/850f90dc-b1d8-4d60-b6f8-7d7944aea69d" />


---

## 📂 Estrutura do Projeto

* `main.go`: Ponto de entrada da Lambda e configuração do handler.
* `internal/handler`: Tradução de eventos HTTP/Lambda e binding de dados.
* `internal/policy`: Core engine (parsing de DOT e avaliação de expressões com `govaluate`).
* `internal/apierror`: Padronização de erros e códigos de retorno.
* `loadtest/`: Manifestos e scripts de teste de performance.

---

## 📈 Roadmap & Próximos Passos (acho interessante agregar)

* [ ] **Custom Metrics:** Finalizar branch `feature/add-metrics` para dashboards no CloudWatch.
* [ ] **Observabilidade:** Configurar alarmes de taxa de erro (>1% por 5min) e latência.
* [ ] **Automação de Testes:** Integrar os testes de carga diretamente no workflow do GitHub Actions.
* [ ] **Multi-environment:** Implementar segregação de ambientes (Staging/Prod) via variáveis de ambiente no CI/CD.
