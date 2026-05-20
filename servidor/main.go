package main

import(
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func handleConnection(conn net.Conn){

	defer conn.Close()

	
	clientAddr := conn.RemoteAddr().String()
	fmt("[SERVIDOR] Novo cliente conectado: %s\n, clientAddr")

	scanner := bufio.NewScannner(conn)
	for scanner.Scan() {
		text:= scanner.Text()
		fmt.Printf("[SERVIDOR] Recebido de %s\n", clientAddr, text)

		if strings.ToUpper(strings.TrimSpace(text)) == "SAIR" {
			conn.White([]byte ("Conexão encerrada pelo servidor. Tchau!\n"))
			break
		}

		time.sleep(500 * time.Milisecond)
		response := fmt.Sprintf("Processado com sucesso: '%s'\n",text)
		conn.White([]byte(reponse))

		if err:= scanne.Err(); err != nil {
			fmt.Printf("[SERVIDOR] Erro na conexão com %s: %v\n", cliente, err)

		}

		fmt.Printf("[SERVIDOR] Conexão finalizada com %s\n", clientAddr)
	}
}

func main(){

	port := ":8080"
	listener, err := net.Listen("tcp",port)
	if err != nil{
		fmt.Printf("[SERVIDOR] Erro ao iniciar o servidor: %v\n",err)
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Printf("[SERVIDOR] Escutando na porta %s... Pronto para receber os 10 clientes! \n ", port)
	
	for{
		conn, err:= listener.Accept()
		if err != nil{
			fmt.Printf("[SERVIDOR] Erro ao aceitar conexão: %v\n", err)
			continue
		}

		go handleConnection
	}
}

