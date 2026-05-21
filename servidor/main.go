package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

    clientAddr := conn.RemoteAddr().String()
    fmt.Printf("[SERVIDOR] Novo cliente conectado: %s\n", clientAddr)

    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        text := scanner.Text()
        fmt.Printf("[SERVIDOR] Recebido de %s: %s\n", clientAddr, text)

        if strings.EqualFold(strings.TrimSpace(text), "SAIR") {
            _, _ = conn.Write([]byte("Conexão encerrada pelo servidor. Tchau!"))
            break
        }

        _, err := conn.Write([]byte(text))
        if err != nil {
            fmt.Printf("[SERVIDOR] Erro ao responder %s: %v\n", clientAddr, err)
            break
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("[SERVIDOR] Erro de leitura em %s: %v\n", clientAddr, err)
    }

    fmt.Printf("[SERVIDOR] Conexão finalizada com %s\n", clientAddr)
}

func main() {
    port := os.Getenv("ECHO_PORT")
    if port == "" {
        port = "5000"
    }

    address := ":" + port
    listener, err := net.Listen("tcp", address)
    if err != nil {
        fmt.Printf("[SERVIDOR] Erro ao iniciar o servidor: %v\n", err)
        os.Exit(1)
    }
    defer listener.Close()

    fmt.Printf("[SERVIDOR] Escutando na porta %s... Pronto para receber conexões.\n", port)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("[SERVIDOR] Erro ao aceitar conexão: %v\n", err)
            continue
        }
        go handleConnection(conn)
    }
}
