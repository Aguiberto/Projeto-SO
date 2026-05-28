package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()

	// 1. Servidor detecta a conexão
	// (vai aparecer logo após o print "1" do cliente)
	// Não colocamos print aqui para não poluir, focamos no fluxo da mensagem:

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()

		// 3. Servidor confirma o recebimento e faz o eco
		fmt.Printf("[SERVIDOR] 3. Servidor confirma o recebimento da mensagem e faz o eco para o cliente.\n")
		fmt.Printf("[SERVIDOR]    Dados recebidos: %s\n", text)

		if strings.EqualFold(strings.TrimSpace(text), "SAIR") {
			_, _ = conn.Write([]byte("Conexão encerrada pelo servidor. Tchau!"))
			break
		}

		// Envia o eco
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Printf("[SERVIDOR] Erro ao responder %s: %v\n", clientAddr, err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[SERVIDOR] Erro de leitura em %s: %v\n", clientAddr, err)
	}

	// Dá um tempo para o cliente receber os dados e printar o passo "4"
	time.Sleep(1500 * time.Millisecond)

	// 5. Servidor informa a finalização
	fmt.Printf("[SERVIDOR] 5. Servidor informa a finalização da conexão com %s.\n", clientAddr)
}

func main() {
	port := os.Getenv("ECHO_PORT"); if port == "" { port = "5000" }
	address := ":" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("[SERVIDOR] Erro ao iniciar: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("[SERVIDOR SÍNCRONO] Escutando na porta %s... Linha do tempo controlada.\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[SERVIDOR] Erro ao aceitar conexão: %v\n", err)
			continue
		}
		handleConnection(conn)
	}
}