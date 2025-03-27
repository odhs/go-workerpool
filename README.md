# _Worker Pool_ Example

## Introdução

O funcionamento desse programa é uma arte!

Uma explição de maneira simples é: A ideia é ter uma função e uma fila de entradas.
Essa fila de entradas irá executar para cada entrada a função dada, porém, paralelamente.
Essa abordagem é muito interessante para aproveitarmos vários processadores ao mesmo tempo, assim conseguindo o resultado muito mais rápido.

Como exemplo podemos citar, o processamento de conversão de vídeo, onde o vídeo pode ser dividido em várias partes menores, processado e depois unido para termos o vídeo convertido. Pode também por exemplo, processar algum conjunto de dados, por exemplo, juntar tabelas, ou listas de emails, pedadinho por pedacinho paralelamente, resultando em um processamento mais rápido.

A seguir mais detalhes e conceitos.

## Uso

```sh
go run main.go
```

Exemplo de saída do programa:

```txt
Starting the worker pool with 20 numbers to be processed
2025/03/27 17:07:41 INFO worker started worker_id=1
2025/03/27 17:07:41 INFO worker started worker_id=0
2025/03/27 17:07:41 INFO worker started worker_id=2
Number: 2, WorkerID: 2, Timestamp: 17:07:42.753
Number: 1, WorkerID: 1, Timestamp: 17:07:42.754
Number: 0, WorkerID: 0, Timestamp: 17:07:42.976
Number: 4, WorkerID: 1, Timestamp: 17:07:43.594
Number: 5, WorkerID: 2, Timestamp: 17:07:43.786
Number: 3, WorkerID: 0, Timestamp: 17:07:43.904
Number: 6, WorkerID: 0, Timestamp: 17:07:44.613
Number: 7, WorkerID: 1, Timestamp: 17:07:44.705
Number: 8, WorkerID: 2, Timestamp: 17:07:45.074
Number: 9, WorkerID: 0, Timestamp: 17:07:45.676
Number: 10, WorkerID: 1, Timestamp: 17:07:45.822
Number: 11, WorkerID: 2, Timestamp: 17:07:45.899
Number: 12, WorkerID: 0, Timestamp: 17:07:46.609
Number: 13, WorkerID: 1, Timestamp: 17:07:46.921
Number: 14, WorkerID: 2, Timestamp: 17:07:47.041
Number: 15, WorkerID: 0, Timestamp: 17:07:47.529
Number: 17, WorkerID: 2, Timestamp: 17:07:48.027
Number: 16, WorkerID: 1, Timestamp: 17:07:48.029
2025/03/27 17:07:48 INFO Worker interrupted by input channel closed worker_id=1
2025/03/27 17:07:48 INFO Worker interrupted by input channel closed worker_id=0
Number: 18, WorkerID: 0, Timestamp: 17:07:48.660
Number: 19, WorkerID: 1, Timestamp: 17:07:49.004
2025/03/27 17:07:49 INFO Worker interrupted by input channel closed worker_id=2
All jobs are done
Every the 20 numbers were processed.
```

## Código

O código fonte encontra-se com comentários em inglês (mas é facilmente traduzível), já os conceitos aqui no `README.md` estão em Língua Portuguesa, destinado aos mesu alunos.

Este exemplo será mantido simples para melhor entendimento. Em um outro repositório está um exemplo prático de uso para converter vídeos em um sistema web, e um outro para processamento de listas de emails.

Este projeto explora os conceitos de **Worker Pool**, **Wait Group** e **Mutex** no contexto de programação concorrente em Go. Para entender melhor você tambpém precisa do conceito de _go routines_ e _canais_.

## Conceitos

Abaixo os conceitos utilizados no software.

### _Worker Pool_

Um **Worker Pool** é um padrão de design usado para gerenciar e processar tarefas de forma eficiente utilizando múltiplas goroutines. Ele é útil para:

- **Controlar a concorrência**: Limitar o número de goroutines em execução simultaneamente, evitando sobrecarregar o sistema.
- **Melhorar o desempenho**: Dividir o trabalho em tarefas menores e distribuí-las entre múltiplos "workers" (trabalhadores), que podem beneficiar-se de um sistema com multiprocessamento.
- **Gerenciar filas de tarefas**: Garantir que todas as tarefas sejam processadas de forma organizada.

#### Como funciona um _Worker Pool_?

1. Um conjunto fixo de _goroutines_ (os "workers") é criado.
2. As tarefas são enviadas para um canal de trabalho.
3. Cada _worker_ consome as tarefas do canal e as processa.
4. Quando todas as tarefas são concluídas, o programa pode encerrar.

Este padrão é amplamente utilizado em sistemas que precisam processar grandes volumes de dados ou realizar operações paralelas, como servidores web, processamento de arquivos ou tarefas de background.

---

### _Wait Group_

O **Wait Group** é uma estrutura fornecida pelo pacote `sync` em Go que permite sincronizar a execução de _goroutines_. Ele é usado para:

- **Aguardar a conclusão de múltiplas goroutines**: O programa principal pode esperar até que todas as goroutines tenham terminado antes de continuar.
- **Evitar condições de corrida (_race conditions_)**: Garante que o programa não finalize antes que todas as tarefas concorrentes sejam concluídas.

#### Como funciona um _Wait Group_?

1. Um `Wait Group` é criado.
2. O método `Add(n)` é chamado para indicar o número de tarefas (_goroutines_) que serão aguardadas. Então teremos o número de tarefas a serem executadas dentro do `Wait Group`.
3. Cada tarefa (_goroutine_) chama `Done()` ao terminar sua execução e nisso subtrai 1 do total no `Wait Group`.
4. O método `Wait()` é chamado para bloquear a execução do restante do software até que todas as tarefas tenham chamado `Done()`.

Exemplo básico:

```go
var wg sync.WaitGroup
var total_tarefas = 3

wg.Add(total_tarefas)

for(i=0; i<total_tarefas; i++){
  go func() {
      defer wg.Done()
      fmt.Println("Goroutine concluída!")
  }()
}

wg.Wait()
fmt.Println("Todas as goroutines foram concluídas.")
```

### _Mutex_

O _Mutex_ _(Mutual Exclusion_) é uma estrutura usada para evitar condições de corrida ao acessar recursos compartilhados em um ambiente concorrente. Ele garante que apenas uma _goroutine_ possa acessar um recurso crítico por vez.

#### Para que serve o _Mutex_?

- **Evitar condições de corrida**: Impede que múltiplas goroutines modifiquem ou leiam um recurso compartilhado simultaneamente.
- **Proteger seções críticas**: Bloqueia o acesso a uma seção do código até que a goroutine atual termine sua execução.

#### Como funciona o _Mutex_?

1. O método `Lock()` é chamado para bloquear o acesso ao recurso.
2. O método `Unlock()` é chamado para liberar o recurso.

Exemplo básico:

```go
var mu sync.Mutex
var contador int

func incrementar() {
    mu.Lock()
    defer mu.Unlock()
    contador++
}
```

### _Context_

É um conjunto de dados que permitem um ambiente (estado atual) ou uma operação serem executados sob regras. Um contexto no geral pode ser definido como um ponto onde recursos estarão disponíveis e vinculados.

Na linaguagem _Go_, um **Contexto** (**_Context_**) é uma estrutura fornecida pelo pacote `context` que permite controlar o ciclo de vida de _goroutines_ e gerenciar dados compartilhados entre elas. Ele é amplamente utilizado para:

- **Cancelar operações**: Permitir que _goroutines_ sejam notificadas para encerrar sua execução quando não forem mais necessárias.
- **Definir prazos (_deadlines_)**: Especificar um tempo limite para a execução de uma operação.
- **Passar valores entre goroutines**: Compartilhar informações relevantes de forma segura e eficiente.

#### Para que serve o _Context_?

O `Context` é útil em situações onde múltiplas _goroutines_ precisam cooperar para realizar uma tarefa, mas podem ser interrompidas caso o trabalho não seja mais necessário. Ele é frequentemente usado em servidores web, APIs e sistemas distribuídos para gerenciar requisições e evitar vazamentos de recursos.

#### Como funciona o _Context_?

1. Um `Context` é criado usando funções como `context.Background()` ou `context.WithCancel()`.
2. O `Context` é passado para funções e _goroutines_ que precisam ser controladas.
3. As _goroutines_ verificam o canal `Done()` do `Context` para saber se devem encerrar sua execução.
4. O `Context` também pode armazenar valores e prazos (_deadlines_) que podem ser acessados pelas _goroutines_.

Exemplo básico:

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    go func(ctx context.Context) {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Goroutine encerrada.")
                return
            default:
                fmt.Println("Goroutine em execução...")
                time.Sleep(500 * time.Millisecond)
            }
        }
    }(ctx)

    time.Sleep(2 * time.Second)
    cancel() // Cancela o contexto, notificando a goroutine para encerrar
    time.Sleep(1 * time.Second)
}
```

## Estrutura do Código

### Pacote _workerpoll_

#### Interfaces e Tipos

- **`Job`**: Representa uma tarefa a ser processada.
- **`Result`**: Representa o resultado de uma tarefa processada.
- **`ProcessFunc`**: Uma função que processa um `Job` e retorna um `Result`.

```go
type Job interface{}
type Result interface{}
type ProcessFunc func(ctx context.Context, job Job) Result
```

- **`WorkerPool`**: Interface que define os métodos principais do pool:
  - `Start(ctx context.Context, inputCh <-chan Job) (<-chan Result, error)`: Inicia o pool de workers.
  - `Stop() error`: Para o pool de workers.
  - `IsRunning() bool`: Verifica se o pool está em execução.

## Referências dos Conceitos

Apresentadas na ordem em que o artigo foi escrito:

- [Worker Pools - Go by Example](https://gobyexample.com/worker-pools)
- [sync.WaitGroup - Documentação Oficial](https://pkg.go.dev/sync#WaitGroup)
- [sync.Mutex - Documentação Oficial](https://pkg.go.dev/sync#Mutex)
- [context - Documentação Oficial](https://pkg.go.dev/context)
- [Contextos em Go - Blog Oficial](https://go.dev/blog/context)

### Referências Adicionais

- [A Tour of Go - Goroutines](https://go.dev/tour/concurrency/1)  
  Introdução oficial às goroutines e como elas permitem a execução de funções de forma concorrente.

- [Go by Example - Channels](https://gobyexample.com/channels)  
  Exemplos práticos de como usar canais para comunicação entre goroutines.

- [A Tour of Go - Channels](https://go.dev/tour/concurrency/2)  
  Explicação oficial sobre canais em Go, usados para comunicação segura entre goroutines.

### Referências Complementares

- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)  
  Guia oficial sobre como escrever código concorrente eficiente e idiomático em Go.

- [log/slog - Documentação Oficial](https://pkg.go.dev/log/slog)  
  Documentação do pacote `slog`, usado para registrar logs estruturados no software.

- [Go Memory Model](https://go.dev/ref/mem)  
  Referência oficial sobre o modelo de memória do Go, essencial para entender como goroutines e canais interagem.

### Referências para a arquitetura do Código

- [Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_ptBR.md)
