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
	// Removemos o defer conn.Close() para controlar o fechamento manualmente no momento exato
	clientAddr := conn.RemoteAddr().String()

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

		// Envia o eco de volta para o cliente
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Printf("[SERVIDOR] Erro ao responder %s: %v\n", clientAddr, err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[SERVIDOR] Erro de leitura em %s: %v\n", clientAddr, err)
	}

	// === SINCRONIZAÇÃO CAUSAL ===
	// Fechamos a conexão com o cliente primeiro. 
	// Isso força o container do cliente a terminar de ler, sair do loop dele e printar o Passo 4.
	conn.Close()

	// Uma pausa de meio segundo garante que o Docker processe e renderize o log do cliente primeiro
	time.Sleep(500 * time.Millisecond)

	// 5. Servidor informa a finalização (Garantido no final de tudo)
	fmt.Printf("[SERVIDOR] 5. Servidor informa a finalização da conexão com %s.\n", clientAddr)
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

	fmt.Printf("[SERVIDOR SÍNCRONO] Escutando na porta %s... Linha do tempo controlada.\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[SERVIDOR] Erro ao aceitar conexão: %v\n", err)
			continue
		}
		// Execução estritamente síncrona (Sem a palavra-chave 'go')
		handleConnection(conn)
	}
}