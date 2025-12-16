# Password Validator API

API RESTful para validação de senhas desenvolvida em Go, seguindo princípios SOLID e clean architecture.

## 🎯 Visão Geral

Esta aplicação expõe uma API web que valida senhas de acordo com critérios específicos de segurança. A solução foi desenvolvida com foco em:

- **Clean Code**: Código limpo, legível e bem documentado
- **SOLID Principles**: Aplicação dos princípios SOLID
- **Testabilidade**: Cobertura completa de testes unitários e de integração
- **Extensibilidade**: Fácil adição de novas regras de validação
- **Observabilidade**: Métricas para monitoramento em produção

## 🔐 Requisitos de Senha

Uma senha é considerada válida quando possui:

- ✅ **9 ou mais caracteres**
- ✅ **Ao menos 1 dígito** (0-9)
- ✅ **Ao menos 1 letra minúscula** (a-z)
- ✅ **Ao menos 1 letra maiúscula** (A-Z)
- ✅ **Ao menos 1 caractere especial** (!@#$%^&*()-+)
- ✅ **Não possuir caracteres repetidos**
- ❌ **Não conter espaços em branco** (espaços são considerados inválidos)

### Exemplos

```go
IsValid("")           // false - vazia
IsValid("aa")         // false - muito curta, caracteres repetidos
IsValid("ab")         // false - muito curta, falta dígito, maiúscula, especial
IsValid("AAAbbbCc")   // false - caracteres repetidos, falta dígito e especial
IsValid("AbTp9!foo")  // false - 'o' repetido
IsValid("AbTp9!foA")  // false - 'A' repetido
IsValid("AbTp9 fok")  // false - contém espaço
IsValid("AbTp9!fok")  // true  - válida!
```

## 🏗️ Arquitetura e Decisões de Design

### Arquitetura em Camadas

A aplicação segue uma **arquitetura em camadas** (layered architecture) para separar responsabilidades:

```
┌─────────────────────────────────────┐
│        API Layer (HTTP)             │  ← Handlers, Middleware, DTOs
├─────────────────────────────────────┤
│     Application Layer (Service)     │  ← Orquestração de validadores
├─────────────────────────────────────┤
│      Domain Layer (Business)        │  ← Regras de validação
└─────────────────────────────────────┘
```

**Benefícios:**
- Separação clara de responsabilidades
- Testabilidade isolada de cada camada
- Facilita manutenção e evolução do código

### Princípios SOLID Aplicados

#### 1. **Single Responsibility Principle (SRP)**
Cada validador tem uma única responsabilidade:
- `MinLengthValidator`: valida apenas o comprimento mínimo
- `DigitValidator`: valida apenas a presença de dígitos
- `NoDuplicatesValidator`: valida apenas duplicatas e espaços

#### 2. **Open/Closed Principle (OCP)**
O sistema é aberto para extensão, mas fechado para modificação:
- Novas regras podem ser adicionadas criando novos validadores
- Não é necessário modificar código existente para adicionar regras

```go
// Adicionar nova regra é simples:
validators = append(validators, rules.NewCustomValidator())
```

#### 3. **Liskov Substitution Principle (LSP)**
Todos os validadores implementam a mesma interface `PasswordValidator`:
```go
type PasswordValidator interface {
    Validate(password string) error
}
```
Qualquer validador pode ser substituído por outro sem quebrar o sistema.

#### 4. **Interface Segregation Principle (ISP)**
Interface minimalista com apenas um método necessário:
```go
type PasswordValidator interface {
    Validate(password string) error  // Apenas o essencial
}
```

#### 5. **Dependency Inversion Principle (DIP)**
Dependências são injetadas via construtor:
```go
service := application.NewPasswordService(validators)  // DI
handler := handlers.NewPasswordHandler(service)        // DI
```

### Padrões de Design

#### Strategy Pattern
Cada regra de validação é uma estratégia independente que pode ser composta:
```go
validators := []domain.PasswordValidator{
    rules.NewMinLengthValidator(9),
    rules.NewDigitValidator(),
    // ... outras estratégias
}
```

## 📁 Estrutura do Projeto

```
itau-backend-challenge/
├── cmd/
│   └── api/
│       └── main.go                  # Entry point da aplicação
├── internal/
│   ├── domain/                      # Camada de domínio (regras de negócio)
│   │   ├── validator.go             # Interface PasswordValidator
│   │   └── rules/                   # Implementações de regras
│   │       ├── min_length.go        # Validador de comprimento mínimo
│   │       ├── digit.go             # Validador de dígitos
│   │       ├── lowercase.go         # Validador de minúsculas
│   │       ├── uppercase.go         # Validador de maiúsculas
│   │       ├── special_char.go      # Validador de caracteres especiais
│   │       ├── no_duplicates.go     # Validador de duplicatas
│   │       └── *_test.go            # Testes unitários
│   ├── application/                 # Camada de aplicação (orquestração)
│   │   ├── password_service.go      # Serviço de validação
│   │   └── password_service_test.go # Testes do serviço
│   └── api/                         # Camada de API (HTTP)
│       ├── handlers/
│       │   └── password_handler.go  # HTTP handlers
│       ├── middleware/
│       │   ├── logging.go           # Middleware de logging
│       │   └── cors.go              # Middleware de CORS
│       └── models/
│           └── request.go           # DTOs (Request/Response)
├── pkg/
│   └── metrics/
│       └── metrics.go               # Métricas Prometheus
├── tests/
│   └── integration/
│       └── api_test.go              # Testes de integração
├── go.mod                           # Dependências Go
├── go.sum                           # Checksums de dependências
└── README.md                        # Este arquivo
```

## 🚀 Como Executar

### Pré-requisitos

- **Go 1.21+** instalado ([Download](https://golang.org/dl/))

### Instalação

1. **Clone o repositório**:
```bash
git clone https://github.com/willherrera/itau-backend-challenge.git
```

2. **Baixe as dependências**:
```bash
go mod tidy
```

3. **Gere a documentação Swagger** (opcional, mas recomendado):
```bash
# Instale o swag CLI (apenas uma vez)
go install github.com/swaggo/swag/cmd/swag@latest

# Gere os arquivos de documentação
swag init -g cmd/api/main.go -o docs
```

### Executar a Aplicação

```bash
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:8080`

**Saída esperada:**
```
Starting password validation API on :8080
Endpoints:
  POST   http://localhost:8080/api/v1/validate-password
  GET    http://localhost:8080/health
  GET    http://localhost:8080/metrics
  GET    http://localhost:8080/swagger/index.html
```

### Executar Testes

**Todos os testes:**
```bash
go test ./... -v
```

**Apenas testes unitários:**
```bash
go test ./internal/... -v
```

**Apenas testes de integração:**
```bash
go test ./tests/integration/... -v
```

**Com cobertura:**
```bash
go test ./... -cover
```

## 📡 Documentação da API

### 📚 Swagger UI (Documentação Interativa)

A API possui documentação interativa completa via **Swagger UI**, que permite:
- ✅ Visualizar todos os endpoints disponíveis
- ✅ Ver schemas de request/response
- ✅ Testar a API diretamente no navegador
- ✅ Ver exemplos de uso

**Acesse:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

![Swagger UI](https://img.shields.io/badge/Swagger-Interactive_Docs-85EA2D?logo=swagger&logoColor=white)

> **Nota:** Certifique-se de gerar a documentação com `swag init` antes de executar a aplicação (veja seção de instalação).

---

### POST /api/v1/validate-password

Valida uma senha de acordo com os critérios estabelecidos.

**Request:**
```json
{
  "password": "AbTp9!fok"
}
```

**Response (Senha Válida):**
```json
{
  "isValid": true
}
```

**Response (Senha Inválida):**
```json
{
  "isValid": false,
  "errors": [
    "password must have at least 9 characters",
    "password must contain at least one digit"
  ]
}
```

**Status Codes:**
- `200 OK`: Validação executada com sucesso
- `400 Bad Request`: JSON inválido
- `405 Method Not Allowed`: Método HTTP não permitido

### GET /health

Verifica o status da aplicação.

**Response:**
```json
{
  "status": "healthy",
  "service": "password-validator"
}
```

### GET /metrics

Expõe métricas Prometheus para monitoramento.

**Response:** Formato Prometheus text-based

## 🧪 Testes

### Estratégia de Testes

A aplicação possui **3 níveis de testes**:

#### 1. Testes Unitários (Domain Layer)
Cada validador é testado isoladamente:
- `min_length_test.go`: 5 casos de teste
- `digit_test.go`: 6 casos de teste
- `lowercase_test.go`: 6 casos de teste
- `uppercase_test.go`: 6 casos de teste
- `special_char_test.go`: 8 casos de teste
- `no_duplicates_test.go`: 13 casos de teste

#### 2. Testes de Serviço (Application Layer)
Testa a orquestração de múltiplos validadores:
- `password_service_test.go`: 11 casos de teste + todos os exemplos dos requisitos

#### 3. Testes de Integração (API Layer)
Testa o fluxo completo HTTP:
- `api_test.go`: 
  - Endpoint de validação
  - Todos os exemplos dos requisitos
  - Health check
  - Tratamento de erros (JSON inválido, método não permitido)

### Cobertura de Testes

Todos os **8 exemplos dos requisitos** são testados:
```go
✓ IsValid("")           // false
✓ IsValid("aa")         // false
✓ IsValid("ab")         // false
✓ IsValid("AAAbbbCc")   // false
✓ IsValid("AbTp9!foo")  // false
✓ IsValid("AbTp9!foA")  // false
✓ IsValid("AbTp9 fok")  // false
✓ IsValid("AbTp9!fok")  // true
```

## 📊 Métricas e Observabilidade

### Métricas Prometheus

A aplicação expõe as seguintes métricas em `/metrics`:

#### Contadores
- `password_validation_requests_total{result="valid|invalid"}`: Total de requisições por resultado
- `password_validation_errors_total{rule="min_length|digit|..."}`: Total de erros por regra

#### Histogramas
- `password_validation_duration_seconds`: Latência das requisições

#### Gauges
- `password_validation_in_progress`: Validações em andamento (concorrência)

### Exemplos de Uso

**Consultar métricas:**
```bash
curl http://localhost:8080/metrics
```

**Exemplo de saída:**
```
# HELP password_validation_requests_total Total number of password validation requests
# TYPE password_validation_requests_total counter
password_validation_requests_total{result="valid"} 42
password_validation_requests_total{result="invalid"} 15

# HELP password_validation_duration_seconds Duration of password validation requests
# TYPE password_validation_duration_seconds histogram
password_validation_duration_seconds_bucket{le="0.005"} 50
password_validation_duration_seconds_sum 0.123
password_validation_duration_seconds_count 57
```

## 🤔 Premissas e Decisões

### Premissas Assumidas

1. **Retornar todos os erros**: Optei por retornar **todos** os erros de validação, não apenas o primeiro. Isso melhora a experiência do usuário, que recebe feedback completo.

2. **HTTP 200 para validação**: Mesmo quando a senha é inválida, retorno `200 OK` porque a **operação de validação** foi bem-sucedida. O campo `isValid` indica o resultado.

3. **Case-sensitive para duplicatas**: 'A' e 'a' são considerados caracteres diferentes para fins de duplicação.

### Decisões Técnicas

#### Por que Go?
- Performance excelente para APIs
- Concorrência nativa (goroutines)
- Tipagem estática (menos bugs)
- Excelente suporte para testes
- Deploy simples (binário único)

#### Por que Gorilla Mux?
- Router HTTP robusto e maduro
- Suporte a middleware
- Fácil definição de rotas
- Amplamente utilizado na comunidade Go

#### Por que Prometheus?
- Padrão de mercado para métricas
- Integração nativa com Kubernetes
- Queries poderosas (PromQL)
- Ecossistema rico (Grafana, AlertManager)

#### Estrutura de Pastas
- `internal/`: Código privado da aplicação (não pode ser importado)
- `pkg/`: Código reutilizável (pode ser importado)
- `cmd/`: Entry points da aplicação
- `tests/`: Testes de integração separados

#### Tratamento de Erros
Cada validador retorna um erro descritivo:
```go
var ErrNoDigit = errors.New("password must contain at least one digit")
```
Isso facilita debugging e fornece mensagens claras ao usuário.

#### Dependency Injection
Todas as dependências são injetadas via construtor:
```go
service := application.NewPasswordService(validators)
```
Isso facilita testes (mock injection) e segue o princípio DIP.

## 🔧 Possíveis Melhorias Futuras

- [ ] Configuração via variáveis de ambiente (comprimento mínimo, caracteres especiais)
- [ ] Rate limiting para proteção contra abuso
- [ ] Cache de validações (para senhas já validadas)
- [ ] Suporte a i18n (internacionalização de mensagens de erro)
- [ ] Docker container para deploy facilitado
- [ ] CI/CD pipeline (GitHub Actions)

## 📝 Exemplos de Uso

### cURL

**Validar senha válida:**
```bash
curl -X POST http://localhost:8080/api/v1/validate-password \
  -H "Content-Type: application/json" \
  -d '{"password":"AbTp9!fok"}'
```

**Validar senha inválida:**
```bash
curl -X POST http://localhost:8080/api/v1/validate-password \
  -H "Content-Type: application/json" \
  -d '{"password":"abc123"}'
```

**Health check:**
```bash
curl http://localhost:8080/health
```

