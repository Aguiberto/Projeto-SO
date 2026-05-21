package main

import (
    "fmt"
    "net"
    "os"
    "sync"
    "time"
)

const (
    defaultHost               = "echo-server"
    defaultPort               = "5000"
    defaultMessage            = "Olá do cliente echo!"
    bufferSize                = 1024
    maxConnectionAttempts     = 15
    connectionTimeoutSeconds  = 5
    clientCount               = 10
)

func runClient(id int, address, message string, timeout time.Duration, wg *sync.WaitGroup, results chan<- bool) {
    defer wg.Done()

    for attempt := 1; attempt <= maxConnectionAttempts; attempt++ {
        conn, err := net.DialTimeout("tcp", address, timeout)
        if err != nil {
            fmt.Printf("[echo-client %d] Tentativa %d/%d falhou: %v\n", id, attempt, maxConnectionAttempts, err)
            time.Sleep(1 * time.Second)
            continue
        }

        fmt.Printf("[echo-client %d] Conectado em %s\n", id, address)
        _, err = conn.Write([]byte(message))
        if err != nil {
            conn.Close()
            fmt.Printf("[echo-client %d] Erro ao enviar mensagem: %v\n", id, err)
            results <- false
            return
        }

        if tcpConn, ok := conn.(*net.TCPConn); ok {
            _ = tcpConn.CloseWrite()
        }

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

        conn.Close()
        response := string(responseBytes)

        fmt.Printf("[echo-client %d] Enviado:  %s\n", id, message)
        fmt.Printf("[echo-client %d] Recebido: %s\n", id, response)

        if response != message {
            fmt.Printf("[echo-client %d] Resposta recebida difere da mensagem enviada. Esperado: %s. Recebido: %s.\n", id, message, response)
            results <- false
            return
        }

        results <- true
        return
    }

    fmt.Fprintf(os.Stderr, "[echo-client %d] Não foi possível conectar ao servidor echo.\n", id)
    results <- false
}

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

    var wg sync.WaitGroup
    results := make(chan bool, clientCount)

    for i := 1; i <= clientCount; i++ {
        wg.Add(1)
        go runClient(i, address, message, timeout, &wg, results)
    }

    wg.Wait()
    close(results)

    success := true
    for ok := range results {
        if !ok {
            success = false
        }
    }

    if !success {
        os.Exit(1)
    }
}
