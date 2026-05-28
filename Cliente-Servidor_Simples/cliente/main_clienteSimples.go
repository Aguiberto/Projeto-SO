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

func runClient(id int, address, message string, timeout, stepDelay time.Duration) bool {
	for attempt := 1; attempt <= maxConnectionAttempts; attempt++ {
		
		// 1. Cliente tenta se conectar
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			fmt.Printf("[echo-client %d] Tentativa %d/%d falhou: %v\n", id, attempt, maxConnectionAttempts, err)
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Printf("[echo-client %d] 1. Cliente se conecta.\n", id)
		time.Sleep(stepDelay) 

		// 2. Cliente envia a mensagem e já confirma o envio
		_, err = conn.Write([]byte(message))
		if err != nil {
			conn.Close()
			fmt.Printf("[echo-client %d] Erro ao enviar mensagem: %v\n", id, err)
			return false
		}
		fmt.Printf("[echo-client %d] 2. Cliente envia mensagem e confirma o envio.\n", id)
		
		// Fecha o fluxo de escrita para o servidor saber que o texto acabou
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			_ = tcpConn.CloseWrite()
		}

		// Pequena pausa para a mensagem trafegar na rede sem atropelos
		time.Sleep(190 * time.Millisecond) 

		// 4. Cliente lê o retorno do Eco
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

		fmt.Printf("[echo-client %d] 4. Cliente confirma o retorno do eco: %s\n", id, response)
		
		// Espera o delay padrão antes de permitir o encerramento do ciclo do cliente
		time.Sleep(stepDelay)

		if response != message {
			fmt.Printf("[echo-client %d] Erro: Resposta diferente da enviada.\n", id)
			return false
		}

		return true
	}
	return false
}

func main() {
	host := os.Getenv("ECHO_HOST"); if host == "" { host = defaultHost }
	port := os.Getenv("ECHO_PORT"); if port == "" { port = defaultPort }
	message := os.Getenv("ECHO_MESSAGE"); if message == "" { message = defaultMessage }

	address := net.JoinHostPort(host, port)
	timeout := time.Duration(connectionTimeoutSeconds) * time.Second

	delayStr := os.Getenv("ECHO_STEP_DELAY"); if delayStr == "" { delayStr = defaultStepDelay }
	stepDelay, _ := time.ParseDuration(delayStr)

	fmt.Println("[SISTEMA] Iniciando simulação sequencial perfeita (Sem Threads)...")

	success := true
	
	// Loop sequencial linear puro (Sem concorrência)
	for i := 1; i <= clientCount; i++ {
		fmt.Printf("\n=============================================\n")
		fmt.Printf("          INICIANDO CICLO DO CLIENTE %d      \n", i)
		fmt.Printf("=============================================\n")
		
		if !runClient(i, address, message, timeout, stepDelay) {
			success = false
			break 
		}
		
		// Margem para o terminal organizar as quebras de linha entre clientes
		time.Sleep(1 * time.Second) 
	}

	if !success {
		os.Exit(1)
	}
	fmt.Println("\n[SISTEMA] Todos os 10 clientes foram atendidos com sucesso, um por vez!")
}