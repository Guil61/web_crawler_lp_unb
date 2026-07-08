# Web Crawler Concorrente em Go
 
**Wiki Técnica do Projeto**
 
- **Disciplina:**  Linguagens de Programação - Universidade de Brasília
- **Linguagem:** Go 1.25+
- **Versão do Documento:** 2.0 — Julho 2026
---
 
## 1. Visão Geral
 
O Web Crawler Concorrente em Go é um programa de linha de comando desenvolvido como projeto da disciplina de Paradigmas de Linguagens de Programação. A partir de uma URL semente, o crawler navega de forma automática pelas páginas da web, extrai links HTML e continua a exploração até atingir uma profundidade máxima configurável.
 
O objetivo principal do projeto é demonstrar na prática o modelo de concorrência de Go: goroutines, channels, worker pool e o paradigma CSP (*Communicating Sequential Processes* — "compartilhe memória comunicando"). Além do aspecto prático, o projeto também serve como estudo comparativo entre Go e C, avaliando critérios de legibilidade, escrita (*writability*) e confiabilidade de cada linguagem.
 
## 2. Contexto Histórico da Linguagem Go
 
A concorrência foi pensada desde o início da linguagem: não é uma biblioteca adicionada depois, é parte da própria linguagem. Essa decisão de projeto é o motivo pelo qual Go foi escolhida como base para este crawler.
 
| Ano | Marco |
|---|---|
| 2007 | Concebida no Google por Robert Griesemer, Rob Pike e Ken Thompson como resposta à lentidão de compilação e à complexidade de C++/Java em larga escala. |
| 2009 | Anunciada publicamente como projeto open source. |
| 2012 | Lançamento do Go 1.0, com garantia de compatibilidade da linguagem. |
| Hoje | Linguagem compilada, estática e com coletor de lixo. Este projeto usa Go 1.25. |
 
## 3. Premissas do Projeto: Usuário e Domínio
 
- **Domínio:** Web / rede. Percorrer páginas HTML a partir de uma URL semente, extraindo e seguindo links.
- **Usuário:** Quem opera pela linha de comando — desenvolvedor, pesquisador, administrador de site.
- **Premissa central:** Baixar páginas é uma operação lenta (espera de rede) → a solução é fazer várias requisições ao mesmo tempo, o que justifica o uso de concorrência.
- **Requisitos:** Não repetir páginas já visitadas, respeitar um limite de requisições por segundo e ser capaz de parar de forma ordenada (Ctrl+C).
## 4. Por Que Escolher Go?
 
- **Concorrência nativa:** Goroutines são tarefas leves — é possível ter milhares rodando ao mesmo tempo com custo mínimo de memória e escalonamento.
- **Comunicação segura:** Channels passam dados entre tarefas sem que o programador precise gerenciar travas manualmente.
- **Filosofia CSP:** "Não se comunique compartilhando memória; compartilhe memória se comunicando."
- **Biblioteca padrão forte:** Cliente HTTP, parser de HTML e primitivas de sincronização já vêm prontos na biblioteca padrão.
- **Binário único:** O programa compila para um binário único, o que o torna rápido, fácil de distribuir e de executar em qualquer sistema operacional suportado.
## 5. Legibilidade e Escrita (Readability & Writability) de Go
 
### 5.1 Legibilidade
 
- **Sintaxe minimalista:** A linguagem possui apenas 25 palavras-chave.
- **gofmt:** A ferramenta gofmt padroniza a formatação — todo código Go se parece, o que elimina discussões de estilo entre a equipe.
- **Simplicidade estrutural:** Sem sobrecarga de operadores nem herança complexa: há menos surpresa ao ler um trecho de código pela primeira vez.
- **Erros explícitos:** O tratamento de erros é explícito e visível diretamente no fluxo do código, em vez de exceções escondidas.
### 5.2 Escrita (Writability)
 
- **Concorrência em uma palavra:** Disparar uma tarefa concorrente é tão simples quanto escrever `go f()`.
- **Inferência de tipo:** O operador `:=` permite inferência de tipo, reduzindo código repetitivo.
- **Composição:** O projeto usa composição por structs e interfaces em vez de hierarquias complexas de herança.
- **Produtividade:** Grande parte da lógica é resolvida com poucas linhas, apoiando-se na biblioteca padrão (`net/http`, `html`, `sync`).
## 6. Confiabilidade de Go
 
- **Tipagem estática:** A tipagem estática forte permite que muitos erros sejam identificados ainda em tempo de compilação.
- **Coletor de lixo:** O coletor de lixo elimina vazamentos de memória e o risco de free duplicado, comuns em C.
- **Sem aritmética de ponteiros:** A ausência de aritmética de ponteiros elimina o risco de estouro de buffer (buffer overflow), um dos erros mais comuns em C.
Esses três fatores, somados ao detector de data races nativo (`go run -race`), tornam Go particularmente adequado para um programa que depende fortemente de concorrência, como este crawler.
 
## 7. Funcionalidades
 
- Crawl recursivo com profundidade configurável: navega N níveis a partir da URL inicial.
- Worker pool concorrente: múltiplos workers processam URLs em paralelo.
- Rate limiter global: controla o número de requisições por segundo para não sobrecarregar o servidor.
- Deduplicação thread-safe: cada URL é visitada no máximo uma vez, mesmo com muitos workers simultâneos.
- Encerramento gracioso: Ctrl+C cancela a operação de forma ordenada via `context`.
- Estatísticas ao final: exibe URLs descobertas, páginas baixadas, links encontrados e erros.
- Detector de data races: compatível com `go run -race` para validação de concorrência.
## 8. Requisitos
 
| Requisito | Versão / Detalhes |
|---|---|
| Go | 1.25 ou superior |
| golang.org/x/net/html | Obtida via `go mod tidy` |
| Conexão com a internet | Necessária para o crawl |
| Sistema Operacional | Linux, macOS ou Windows |
 
## 9. Instalação e Execução
 
### 9.1 Clonando e baixando dependências
 
```bash
go mod tidy # baixa a dependência golang.org/x/net/html
```
 
### 9.2 Executando diretamente
 
```bash
go run ./cmd -url=https://www.unb.br/ -depth=2 -workers=5
```
 
### 9.3 Compilando o binário
 
```bash
go build -o crawler ./cmd
./crawler -url=https://www.unb.br/ -depth=2 -workers=5
```
 
### 9.4 Executando com detector de data races
 
```bash
go run -race ./cmd -url=https://www.unb.br/ -depth=1 -workers=50
```
 
Para encerrar a execução a qualquer momento, pressione `Ctrl+C`.
 
## 10. Flags de Configuração
 
| Flag | Padrão | Obrigatório | Descrição |
|---|---|---|---|
| `-url` | (nenhum) | Sim | URL inicial (semente) do crawl |
| `-depth` | 2 | Não | Profundidade máxima de navegação |
| `-workers` | 5 | Não | Número de workers concorrentes |
| `-rate` | 10 | Não | Limite de requisições por segundo (global) |
 
## 11. Estrutura do Projeto
 
| Arquivo / Pacote | Responsabilidade |
|---|---|
| `cmd/main.go` | Ponto de entrada: leitura de flags, sinal de shutdown e inicialização |
| `internal/crawler` | Orquestração geral: cria channels, pool, scheduler e exibe estatísticas |
| `internal/models` | Tipos de dados compartilhados: Job, PageResult e Stats |
| `internal/parser` | Extração de links de documentos HTML |
| `internal/storage` | Conjunto de URLs visitadas, thread-safe |
| `internal/worker` | Worker pool com HTTP client e rate limiter |
| `internal/scheduler` | Fila de jobs, deduplicação e controle de profundidade |
 
## 12. Arquitetura e Fluxo de Dados
 
O sistema segue o padrão Pipeline com Worker Pool — um dos padrões clássicos de concorrência em Go. A comunicação entre componentes é feita exclusivamente via channels, sem memória compartilhada direta (exceto o `VisitedSet` e o `Stats`, ambos protegidos por primitivas atômicas ou mutex).
 
### 12.1 Fluxo principal
 
1. Scheduler injeta a URL semente no channel `jobs`.
2. Workers consomem o channel `jobs`, fazem o HTTP GET, extraem links e devolvem `PageResult` pelo channel `results`.
3. Scheduler consome o channel `results`, descarta URLs já visitadas e enfileira as novas.
4. Quando a fila esvazia e não há jobs em voo, o Scheduler fecha o channel `jobs`.
5. Os workers encerram ao receber o fechamento do channel e o `pool.Wait()` desbloqueia.
### 12.2 Channels utilizados
 
| Channel | Tipo | Capacidade | Produtor | Consumidor |
|---|---|---|---|---|
| `jobs` | `models.Job` | = workers | Scheduler | Workers |
| `results` | `models.PageResult` | = workers | Workers | Scheduler |
| `done` | `struct{}` | 0 (sync) | goroutine do Scheduler | `Crawler.Run` |
 
## 13. Detalhamento dos Componentes
 
### 13.1 `main.go` — Ponto de Entrada
 
Responsável por parsear as flags de CLI (`-url`, `-depth`, `-workers`, `-rate`), configurar o `context` com cancelamento por sinal (SIGINT/SIGTERM) e delegar a execução ao `crawler.New(cfg).Run(ctx)`.
 
### 13.2 `crawler.go` — Orquestrador
 
Cria todos os channels e componentes, dispara o Scheduler em uma goroutine separada, aguarda o encerramento via channel `done` e exibe as estatísticas finais. É o único ponto que conhece todos os subsistemas.
 
### 13.3 `models.go` — Tipos de Dados
 
| Tipo | Campos | Descrição |
|---|---|---|
| `Job` | `URL string`, `Depth int` | Unidade de trabalho enviada aos workers |
| `PageResult` | `Job`, `Links []string`, `WorkerID int`, `Err error` | Resultado retornado pelo worker ao Scheduler |
| `Stats` | `PagesVisited`, `PagesFetched`, `LinksFound`, `Errors int64` | Contadores atômicos de métricas |
 
O tipo `Stats` usa `sync/atomic` para garantir leituras e escritas seguras de múltiplos workers sem mutex.
 
### 13.4 `parser.go` — Extração de Links
 
Usa o tokenizador de `golang.org/x/net/html` para percorrer o HTML token a token e extrair atributos `href` de tags `<a>`. A função `resolve()` converte URLs relativas em absolutas e descarta esquemas não-HTTP/HTTPS e fragmentos (`#anchor`).
 
### 13.5 `visited.go` — Conjunto de URLs Visitadas
 
Implementa um `map[string]struct{}` protegido por `sync.Mutex`. O método `Visit(url)` é atômico: verifica e insere em uma única operação, retornando `true` apenas na primeira vez que a URL é vista.
 
### 13.6 `worker.go` — Pool de Workers e Rate Limiter
 
- **Rate Limiter:** usa `time.Ticker` com intervalo calculado como `1s / rate`. Workers chamam `Wait(ctx)`, que bloqueia até o próximo tick ou até o contexto ser cancelado.
- **Worker Pool:** cada worker roda em uma goroutine independente, consumindo o channel `jobs`. O cliente HTTP tem timeout de 10 segundos. O User-Agent é identificado como `GoConcurrentCrawler/1.0` (projeto acadêmico).
### 13.7 `scheduler.go` — Fila e Controle de Profundidade
 
Mantém uma fila em memória (slice de `Job`) e um contador `inFlight` de jobs em processamento. Usa o padrão de `select` com canal nulo: quando não há pending, `sendCh` é `nil` e o case de envio é ignorado automaticamente, evitando busy-wait.
 
Ao receber um resultado, verifica se `nextDepth <= maxDepth` e só enfileira links novos (não visitados). Quando `len(pending) == 0 && inFlight == 0`, o loop termina e o channel `jobs` é fechado via `defer`.
 
## 14. Construtores de Go Utilizados no Projeto
 
A tabela a seguir resume os principais construtores da linguagem empregados na implementação e o papel de cada um dentro do crawler.
 
| Construtor | Papel no Projeto |
|---|---|
| goroutine (`go f()`) | Tarefa concorrente leve; usada para cada worker e para o loop do scheduler. |
| channel (`chan T`) | Esteira de comunicação entre goroutines (`jobs`, `results`, `done`). |
| `select` | Espera em vários channels ao mesmo tempo; usada pelo scheduler e pelo rate limiter. |
| `defer` | Adia execução (fechar channel, destravar mutex) de forma segura. |
| struct / interface | Modelagem de dados (`Job`, `PageResult`, `Stats`) e comportamento. |
| slice e map | Coleções seguras com verificação de limites (fila de jobs, conjunto de URLs visitadas). |
| `sync.Mutex` / `sync/atomic` | Sincronização do `VisitedSet` e dos contadores de `Stats`. |
| `context` | Cancelamento propagado a todos os workers ao pressionar Ctrl+C. |
 
## 15. Concorrência e Segurança
 
| Componente | Mecanismo de Proteção | Justificativa |
|---|---|---|
| `VisitedSet` | `sync.Mutex` | Operação check-and-set deve ser atômica |
| `Stats` | `sync/atomic` | Contadores independentes; atomic é mais eficiente que mutex |
| `jobs` channel | Fechamento pelo Scheduler | Sinaliza fim de trabalho aos workers (Go idiomático) |
| Cancelamento | `context.Context` | Propagação de shutdown por todos os subsistemas |
 
O projeto não usa `sync.WaitGroup` para workers diretamente no `main` — o pool encapsula internamente o `wg.Wait()`. O cancelamento via `context` garante que todos os workers parem de forma ordenada ao receber Ctrl+C.
 
## 16. Comparação Go × C
 
Como parte da avaliação da disciplina, o projeto também documenta a diferença entre a abordagem adotada em Go e a abordagem equivalente que seria necessária em C para os mesmos problemas.
 
| Aspecto | Go — Como Fizemos no Projeto | C — Abordagem Equivalente |
|---|---|---|
| Concorrência | `go` + channel — fila e sincronização prontas. | `pthread_create` + fila manual com mutex e variáveis de condição. |
| Memória | Coletor de lixo cuida de tudo. | `malloc`/`free` na mão — risco de vazamento. |
| Erros | Retorno de valor (`links, err`). | Códigos de retorno / `errno`. |
| Coleções | slice/map com limites verificados. | Arrays + ponteiros — risco de estouro de buffer. |
 
## 17. Saída do Programa
 
### 17.1 Durante a execução
 
```
[Worker 3] Crawling https://www.unb.br/ (depth 0)
[Scheduler] https://www.unb.br/ -> 42 links (38 novos enfileirados)
```
 
### 17.2 Estatísticas finais
 
```
======== ESTATÍSTICAS FINAIS ========
URLs únicas descobertas : 128
Páginas baixadas (200)  : 115
Links encontrados       : 3204
Erros                   : 13
Tempo total              : 4.321s
=====================================
```
 
## 18. Dependências
 
| Pacote | Uso | Como Obter |
|---|---|---|
| `golang.org/x/net/html` | Tokenização e parsing de HTML | `go mod tidy` |
| `sync`, `sync/atomic` | Mutex e operações atômicas (stdlib) | Nativo do Go |
| `net/http` | Cliente HTTP para buscar páginas (stdlib) | Nativo do Go |
| `context` | Propagação de cancelamento (stdlib) | Nativo do Go |
| `time` | Rate limiter via Ticker (stdlib) | Nativo do Go |
 
## 19. Referências
 
- Go Tour — Concorrência: https://go.dev/tour/concurrency/1
- Go Blog — Pipelines e cancelamento: https://go.dev/blog/pipelines
- CSP — Hoare, C.A.R. (1978): *Communicating Sequential Processes*, CACM
- golang.org/x/net/html: https://pkg.go.dev/golang.org/x/net/html
- Go Race Detector: https://go.dev/doc/articles/race_detector
## 20. Equipe do Projeto e Divisão de Responsabilidades
 
O desenvolvimento do crawler foi dividido entre os quatro integrantes do grupo, cada um responsável por um subconjunto coeso de pacotes e funcionalidades, conforme detalhado a seguir.
 
| Integrante | Área de Responsabilidade |
|---|---|
| **André Silva Guil** | Orquestração geral do sistema (`cmd/main.go` e `internal/crawler`): leitura de flags, inicialização do context de cancelamento, criação dos channels `jobs`/`results`/`done` e coordenação entre todos os subsistemas. |
| **Arthur Gonçalves Maia Lima da Silva** | Worker pool e extração de conteúdo (`internal/worker` e `internal/parser`): implementação do rate limiter com `time.Ticker`, dos workers concorrentes com o cliente HTTP e do parser de tokens HTML para extração de links. |
| **Pedro Henrique Araujo Silva** | Fila e deduplicação (`internal/scheduler` e `internal/storage`): lógica de enfileiramento com o padrão select/canal nulo, controle de profundidade máxima e o conjunto thread-safe de URLs visitadas. |
| **Davi Alemar de Souza Guimaraes** | Modelagem de dados e validação de concorrência (`internal/models` e testes): definição dos tipos `Job`, `PageResult` e `Stats` com contadores atômicos, e verificação de ausência de data races com `go run -race`. |
 
