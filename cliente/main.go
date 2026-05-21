package main

import (
    "fmt"
    "net"
    "os"
    "time"
)

const (
    defaultHost               = "echo-server"
    defaultPort               = "5000"
    defaultMessage            = "Olá do cliente echo!"
    bufferSize                = 1024
    maxConnectionAttempts     = 15
    connectionTimeoutSeconds  = 5
)

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

    for attempt := 1; attempt <= maxConnectionAttempts; attempt++ {
        conn, err := net.DialTimeout("tcp", address, timeout)
        if err != nil {
            fmt.Printf("[echo-client] Tentativa %d/%d falhou: %v. Aguardando servidor...\n", attempt, maxConnectionAttempts, err)
            time.Sleep(1 * time.Second)
            continue
        }

        fmt.Printf("[echo-client] Conectado em %s\n", address)
        _, err = conn.Write([]byte(message))
        if err != nil {
            conn.Close()
            fmt.Printf("[echo-client] Erro ao enviar mensagem: %v\n", err)
            os.Exit(1)
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

        fmt.Printf("[echo-client] Enviado:  %s\n", message)
        fmt.Printf("[echo-client] Recebido: %s\n", response)

        if response != message {
            fmt.Printf("[echo-client] Resposta recebida difere da mensagem enviada. Esperado: %s. Recebido: %s.\n", message, response)
            os.Exit(1)
        }

        return
    }

    fmt.Fprintln(os.Stderr, "[echo-client] Não foi possível conectar ao servidor echo.")
    os.Exit(1)
}
