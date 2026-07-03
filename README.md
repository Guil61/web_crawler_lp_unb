# Web Crawler Concorrente em Go
# Wiki do projeto

Projeto da disciplina de Paradigmas de Linguagens de Programação.

É um web crawler de terminal: a partir de uma URL inicial, ele visita as
páginas, extrai os links e segue navegando até uma profundidade definida.
A ideia do trabalho é mostrar na prática a concorrência em Go — goroutines,
channels, worker pool e o modelo CSP ("compartilhe memória comunicando").

## Como rodar

Precisa de Go 1.25+ instalado (`go version`) e conexão com a internet.

```bash
go mod tidy    # baixa a dependência (golang.org/x/net/html)
go run ./cmd -url=https://www.unb.br/ -depth=2 -workers=5
```
\
Ou gerando o binário:

```bash
go build -o crawler ./cmd
./crawler -url=https://www.unb.br/ -depth=2 -workers=5
```

## Flags

| Flag       | Padrão          | Descrição                           |
|------------|-----------------|-------------------------------------|
| `-url`     | (obrigatório)   | URL inicial do crawl                |
| `-depth`   | `2`             | Profundidade máxima de navegação    |
| `-workers` | `5`             | Número de workers concorrentes      |
| `-rate`    | `10`            | Limite de requisições por segundo   |

Ctrl+C encerra o programa de forma ordenada. Para rodar com o detector de
data races: `go run ./cmd -url=https://www.unb.br/ -depth=1 -workers=50`.

## Estrutura

```
cmd/main.go          ponto de entrada (flags + CLI)
internal/models      tipos de dados (Job, PageResult, Stats)
internal/parser      extração de links do HTML
internal/storage     conjunto de URLs visitadas (thread-safe)
internal/worker      worker pool + rate limiter
internal/scheduler   fila, deduplicação e controle de profundidade
internal/crawler     junta tudo e orquestra a execução
```

## Como funciona

O scheduler envia URLs (`Job`) por um channel; um pool de workers consome
esses jobs, faz o HTTP GET, extrai os links e devolve o resultado por outro
channel. O scheduler recebe os resultados, descarta URLs já visitadas e
enfileira as novas. Quando a fila esvazia, ele fecha o channel de jobs e os
workers encerram.
