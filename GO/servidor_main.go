package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func handleClient(conn net.Conn, processor1Addr string, processor2Addr string) {
	addr := conn.RemoteAddr()
	fmt.Println(addr, ": connected...")
	defer fmt.Println(addr, ": disconnected...")
	defer conn.Close()

	// Recebe a senha do cliente
	netData, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler dados do cliente:", err)
		return
	}

	// Redireciona a senha para os dois servidores de processamento
	processorConn1, err1 := net.Dial("tcp", processor1Addr)
	processorConn2, err2 := net.Dial("tcp", processor2Addr)
	if err1 != nil || err2 != nil {
		fmt.Println("Erro ao conectar aos servidores de processamento:", err1, err2)
		return
	}
	defer processorConn1.Close()
	defer processorConn2.Close()

	// Envia a senha para ambos os servidores
	fmt.Fprintf(processorConn1, netData)
	fmt.Fprintf(processorConn2, netData)

	// Aguarda a resposta de qualquer um dos dois
	processorResponse1 := make(chan string)
	processorResponse2 := make(chan string)

	go func() {
		response, _ := bufio.NewReader(processorConn1).ReadString('\n')
		processorResponse1 <- response
	}()

	go func() {
		response, _ := bufio.NewReader(processorConn2).ReadString('\n')
		processorResponse2 <- response
	}()

	// Retorna a resposta do primeiro que encontrar a senha correta
	select {
	case result := <-processorResponse1:
		conn.Write([]byte("Processor 1 found the password: " + result))
	case result := <-processorResponse2:
		conn.Write([]byte("Processor 2 found the password: " + result))
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Enter with port number in argument")
		return
	}

	// Endereços dos servidores de processamento
	processor1Addr := "localhost:9001"
	processor2Addr := "localhost:9002"

	// Inicia o servidor principal
	l, err := net.Listen("tcp", ":"+arguments[1])
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		return
	}
	defer l.Close()
	fmt.Println("Servidor principal iniciado na porta", arguments[1])

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			return
		}
		go handleClient(c, processor1Addr, processor2Addr)
	}
}
