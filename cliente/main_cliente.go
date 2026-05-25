package main

import (
    "fmt"   // Formar e imprimir textos
    "net"   // Pacote de rede
    "os"    // 
    "sync"  // Ferramenta para trabalhar concorrência
    "time"  // Biblioteca para trabalhar com tempos
)

// =============== CONFIGURAÇÕES GERAIS ===================

const (
    defaultHost               = "echo-server"               // Nome do servidor no docker
    defaultPort               = "5000"                      // Porta de conexão
    defaultMessage            = "Olá do cliente echo!"
<<<<<<< HEAD:cliente/main.go
    bufferSize                = 1024                        // tamanho da leitura dos dados
    maxConnectionAttempts     = 15                          // Quantidade de tentativas de reconexão
    connectionTimeoutSeconds  = 5                           // Tempo máximo esperando a resposta do servidor
    clientCount               = 10                          // Quantidade de clientes simultâneos
)


func runClient(id int, address, message string, timeout time.Duration, wg *sync.WaitGroup, results chan<- bool) {

    // avisa que a função foi concluída
=======
    bufferSize                = 1024
    maxConnectionAttempts     = 15
    connectionTimeoutSeconds  = 5
    defaultStepDelay          = "1s"
    clientCount               = 10
)

func runClient(id int, address, message string, timeout, stepDelay time.Duration, wg *sync.WaitGroup, results chan<- bool) {
>>>>>>> 71c64619b026a4cbc70329c9e8ae23c7e8a77657:cliente/main_cliente.go
    defer wg.Done()

    // loop de tentativas caso o servidor demora a subir no Docker
    for attempt := 1; attempt <= maxConnectionAttempts; attempt++ {
<<<<<<< HEAD:cliente/main.go

        // Tenta conectar ao servidor TCP
=======
        fmt.Printf("[echo-client %d] Aguardando %s antes da tentativa %d/%d...\n", id, stepDelay, attempt, maxConnectionAttempts)
        time.Sleep(stepDelay)

>>>>>>> 71c64619b026a4cbc70329c9e8ae23c7e8a77657:cliente/main_cliente.go
        conn, err := net.DialTimeout("tcp", address, timeout)
        if err != nil {
            fmt.Printf("[echo-client %d] Tentativa %d/%d falhou: %v\n", id, attempt, maxConnectionAttempts, err)
            time.Sleep(1 * time.Second)
            continue
        }

        fmt.Printf("[echo-client %d] Conectado em %s\n", id, address)
<<<<<<< HEAD:cliente/main.go

        // Envia a mensagem de texto convertida em bytes para o servidor
=======
        fmt.Printf("[echo-client %d] Enviando mensagem...\n", id)
>>>>>>> 71c64619b026a4cbc70329c9e8ae23c7e8a77657:cliente/main_cliente.go
        _, err = conn.Write([]byte(message))
        if err != nil {
            conn.Close()
            fmt.Printf("[echo-client %d] Erro ao enviar mensagem: %v\n", id, err)
            results <- false // Envia sinal de falha pelo canal
            return
        }

        // Avisa ao servidor: 
        if tcpConn, ok := conn.(*net.TCPConn); ok {
            _ = tcpConn.CloseWrite()
        }

<<<<<<< HEAD:cliente/main.go

        // Cria o buffer (espaço em memória) para ler a resposta do servidor
=======
        fmt.Printf("[echo-client %d] Mensagem enviada, aguardando a resposta...\n", id)
        time.Sleep(stepDelay)

>>>>>>> 71c64619b026a4cbc70329c9e8ae23c7e8a77657:cliente/main_cliente.go
        responseBytes := make([]byte, 0, bufferSize)
        buffer := make([]byte, bufferSize)
        for {
            n, readErr := conn.Read(buffer)
            if n > 0 {
                responseBytes = append(responseBytes, buffer[:n]...)
            }
            if readErr != nil {
                break
            }
        }

        conn.Close() // fecha totalmente a conexão do cliente com o servidor
        response := string(responseBytes) // converte os bytes recebidos em texto 

        fmt.Printf("[echo-client %d] Enviado:  %s\n", id, message)
        fmt.Printf("[echo-client %d] Recebido: %s\n", id, response)

        // Certificação do Echo : verifica se a mensagem recebida é igual a enviada
        if response != message {
            fmt.Printf("[echo-client %d] Resposta recebida difere da mensagem enviada. Esperado: %s. Recebido: %s.\n", id, message, response)
            results <- false
            return
        }

        results <- true
        return
    }

    // Se ultrapassar 15 tentativas sem conseguir se conectar envia a mensagem e retorna false
    fmt.Fprintf(os.Stderr, "[echo-client %d] Não foi possível conectar ao servidor echo.\n", id)
    results <- false
}

// Lê o ambiente e orquestra o disparo das threads

func main() {


    host := os.Getenv("ECHO_HOST")
    if host == "" {
        host = defaultHost
    }

    port := os.Getenv("ECHO_PORT")
    if port == "" {
        port = defaultPort
    }

    message := os.Getenv("ECHO_MESSAGE")
    if message == "" {
        message = defaultMessage
    }

    address := net.JoinHostPort(host, port)
    timeout := time.Duration(connectionTimeoutSeconds) * time.Second

<<<<<<< HEAD:cliente/main.go
    
    var wg sync.WaitGroup // funciona como contador 
    results := make(chan bool, clientCount) // canal para receber os resultados
=======
    delayStr := os.Getenv("ECHO_STEP_DELAY")
    if delayStr == "" {
        delayStr = defaultStepDelay
    }

    stepDelay, err := time.ParseDuration(delayStr)
    if err != nil || stepDelay < 0 {
        fmt.Fprintf(os.Stderr, "[echo-client] ECHO_STEP_DELAY inválido: %s. Usando %s por padrão.\n", delayStr, defaultStepDelay)
        stepDelay = time.Second
    }

    var wg sync.WaitGroup
    results := make(chan bool, clientCount)
>>>>>>> 71c64619b026a4cbc70329c9e8ae23c7e8a77657:cliente/main_cliente.go

    // Aqui acontece o disparo das THREADS
    for i := 1; i <= clientCount; i++ {
<<<<<<< HEAD:cliente/main.go
        wg.Add(1) // Avisa a criação de mais uma THREADS

        // Cria uma Goroutine que permite executar a função em segundo plano enquanto o laço permanece em execução
        go runClient(i, address, message, timeout, &wg, results)
=======
        wg.Add(1)
        go runClient(i, address, message, timeout, stepDelay, &wg, results)
>>>>>>> 71c64619b026a4cbc70329c9e8ae23c7e8a77657:cliente/main_cliente.go
    }

    // garante que o programa não vai encerrar enquanto o todas threads não tiverem sido executadas
    wg.Wait()
    close(results)

    // verifica se alguma thread falhou
    success := true
    for ok := range results {
        if !ok {
            success = false
        }
    }

    // encerra o programa com o código da falha se houve algum erro
    if !success {
        os.Exit(1)
    }
}
