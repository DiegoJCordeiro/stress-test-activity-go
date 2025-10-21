# Sistema de Teste de Carga em Go

Sistema CLI desenvolvido em Go para realizar testes de carga em serviços web, permitindo avaliar a performance e estabilidade de aplicações através de requisições HTTP concorrentes.

## Funcionalidades

- ✅ Execução de múltiplas requisições HTTP simultâneas
- ✅ Controle preciso de concorrência
- ✅ Relatório detalhado com métricas de performance
- ✅ Distribuição de códigos de status HTTP
- ✅ Indicador de progresso em tempo real
- ✅ Execução via Docker

## Requisitos

- Docker (para execução via container)
- Go 1.21+ (para execução local)

## Instalação e Uso

### Opção 1: Executar com Docker (Recomendado)

#### 1. Build da imagem Docker

```bash
docker build -t loadtest .
```

#### 2. Executar teste de carga

```bash
docker run loadtest --url=http://google.com --requests=1000 --concurrency=10
```

### Opção 2: Executar localmente

#### 1. Instalar dependências

```bash
go mod download
```

#### 2. Compilar aplicação

```bash
go build -o loadtest main.go
```

#### 3. Executar teste

```bash
./loadtest --url=http://google.com --requests=1000 --concurrency=10
```

## Parâmetros

| Parâmetro | Descrição | Obrigatório | Exemplo |
|-----------|-----------|-------------|---------|
| `--url` | URL do serviço a ser testado | Sim | `--url=http://google.com` |
| `--requests` | Número total de requisições | Sim | `--requests=1000` |
| `--concurrency` | Número de requisições simultâneas | Sim | `--concurrency=10` |

## Exemplos de Uso

### Teste básico com 100 requisições e concorrência 10

```bash
docker run loadtest --url=https://api.example.com --requests=100 --concurrency=10
```

### Teste de alta carga com 10000 requisições e concorrência 100

```bash
docker run loadtest --url=https://api.example.com --requests=10000 --concurrency=100
```

### Teste em endpoint específico

```bash
docker run loadtest --url=https://api.example.com/v1/users --requests=500 --concurrency=50
```

## Relatório de Saída

O sistema gera um relatório completo contendo:

- **Tempo total de execução**: Duração completa do teste
- **Quantidade total de requests**: Total de requisições realizadas
- **Requests com status 200**: Requisições bem-sucedidas
- **Distribuição de status HTTP**: Contagem detalhada de todos os códigos de status retornados
- **Erros de conexão**: Requisições que falharam antes de receber resposta

### Exemplo de Relatório

```
Iniciando teste de carga...
URL: http://google.com
Total de requests: 1000
Concorrência: 10

Progresso: 1000/1000 requests completados

========== RELATÓRIO DE TESTE DE CARGA ==========
Tempo total de execução: 15.234s
Quantidade total de requests: 1000
Requests com status 200: 985

Distribuição de status HTTP:
  Status 200: 985 requests
  Status 301: 10 requests
  Status 500: 3 requests
  Erros de conexão: 2
================================================
```