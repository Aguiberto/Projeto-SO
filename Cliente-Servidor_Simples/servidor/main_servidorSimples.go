package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Trata a conexão de forma SÍNCRONA. O servidor fica preso aqui até o cliente desconectar.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("[SERVIDOR SÍNCRONO] Novo cliente conectado: %s\n", clientAddr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("[SERVIDOR SÍNCRONO] Recebido de %s: %s\n", clientAddr, text)

		// O DELAY SOLICITADO: O servidor vai travar por 3 segundos para responder.
		// Como não há threads, o servidor INTEIRO fica congelado aqui para qualquer outro cliente.
		time.Sleep(3 * time.Second)

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

	fmt.Printf("[SERVIDOR SÍNCRONO] Conexão finalizada com %s\n", clientAddr)
}

func main() {
	port := os.Getenv("ECHO_PORT")
	if port == "" {
		port = "5000"
	}

	address := ":" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("[SERVIDOR] Erro ao iniciar: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("[SERVIDOR SÍNCRONO] Escutando na porta %s... Atendendo UM POR VEZ.\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[SERVIDOR] Erro ao aceitar conexão: %v\n", err)
			continue
		}
		
		// REMOVIDO O 'go': Agora a função roda na thread principal.
		// O servidor NÃO volta para o início do loop até que handleConnection termine!
		handleConnection(conn)
	}
}