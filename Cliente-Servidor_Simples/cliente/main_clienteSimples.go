package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	defaultHost              = "echo-server"
	defaultPort              = "5000"
	defaultMessage           = "Olá do cliente echo!"
	bufferSize               = 1024
	maxConnectionAttempts    = 15
	connectionTimeoutSeconds = 5
	defaultStepDelay         = "2s" // Tempo de espera entre os passos
	clientCount              = 10   // Executaremos 10 clientes, um após o outro
)

// Agora é uma função comum, sem WaitGroup e sem Canais
func runClient(id int, address, message string, timeout, stepDelay time.Duration) bool {
	for attempt := 1; attempt <= maxConnectionAttempts; attempt++ {
		
		// 1. Tenta conectar
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			fmt.Printf("[echo-client %d] Tentativa %d/%d falhou: %v\n", id, attempt, maxConnectionAttempts, err)
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Printf("[echo-client %d] ---> CONECTADO em %s\n", id, address)
		time.Sleep(stepDelay) // Pausa para leitura

		// 2. Envia a mensagem
		fmt.Printf("[echo-client %d] Enviando mensagem...\n", id)
		_, err = conn.Write([]byte(message))
		if err != nil {
			conn.Close()
			fmt.Printf("[echo-client %d] Erro ao enviar mensagem: %v\n", id, err)
			return false
		}

		// Avisa o fechamento da escrita
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			_ = tcpConn.CloseWrite()
		}

		fmt.Printf("[echo-client %d] Mensagem enviada. Aguardando resposta do servidor...\n", id)

		// 3. Lê o retorno (Echo)
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

		fmt.Printf("[echo-client %d] <--- RECEBIDO: %s\n", id, response)
		time.Sleep(stepDelay) // Pausa para leitura antes de fechar o ciclo do cliente

		if response != message {
			fmt.Printf("[echo-client %d] Erro: Resposta diferente da enviada.\n", id)
			return false
		}

		return true
	}
	return false
}

func main() {
	host := os.Getenv("ECHO_HOST")
	if host == "" { host = defaultHost }
	port := os.Getenv("ECHO_PORT")
	if port == "" { port = defaultPort }
	message := os.Getenv("ECHO_MESSAGE")
	if message == "" { message = defaultMessage }

	address := net.JoinHostPort(host, port)
	timeout := time.Duration(connectionTimeoutSeconds) * time.Second

	delayStr := os.Getenv("ECHO_STEP_DELAY")
	if delayStr == "" { delayStr = defaultStepDelay }
	stepDelay, _ := time.ParseDuration(delayStr)

	fmt.Println("[SISTEMA] Iniciando simulação sequencial (Sem Threads)...")

	success := true
	
	// O LOOP AGORA É TOTALMENTE LINEAR!
	// O cliente 'i+1' só vai começar quando a função runClient(i) terminar completamente e retornar.
	for i := 1; i <= clientCount; i++ {
		fmt.Printf("\n=== INICIANDO CICLO DO CLIENTE %d ===\n", i)
		
		if !runClient(i, address, message, timeout, stepDelay) {
			success = false
			break // Para o loop se algum falhar
		}
		
		// Uma folga de 1 segundo entre um cliente e outro para o terminal respirar
		time.Sleep(1 * time.Second) 
	}

	if !success {
		os.Exit(1)
	}
	fmt.Println("\n[SISTEMA] Todos os 10 clientes foram atendidos com sucesso, um por vez!")
}