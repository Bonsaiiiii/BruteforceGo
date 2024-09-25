package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// Função para gerar todas as combinações de senhas possíveis (força bruta)
func bruteForcePassword(password string, maxLength int) string {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var attempt string

	// Gerar combinações de diferentes tamanhos de caracteres
	for length := 1; length <= maxLength; length++ {
		attempt = generatePasswordRecursive(charSet, "", length, password)
		if attempt == password {
			return attempt
		}
	}
	return ""
}

// Função recursiva para gerar as combinações de caracteres
func generatePasswordRecursive(charSet, current string, length int, target string) string {
	if len(current) == length {
		fmt.Println("Tentando senha:", current)
		if current == target {
			return current
		}
		time.Sleep(50 * time.Millisecond) // Simulação de tempo de processamento
		return ""
	}
	for _, char := range charSet {
		found := generatePasswordRecursive(charSet, current+string(char), length, target)
		if found != "" {
			return found
		}
	}
	return ""
}

// Função que processa cada conexão em uma goroutine separada
func handleConnection(conn net.Conn) {
	addr := conn.RemoteAddr()
	fmt.Println(addr, ": conectado ao servidor de processamento...")
	defer fmt.Println(addr, ": desconectado do servidor de processamento...")
	defer conn.Close()

	// Ler a senha enviada pelo servidor principal
	request, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler requisição:", err)
		return
	}

	password := strings.TrimSpace(request)
	maxLength := len(password) // Determina o tamanho máximo da senha
	result := bruteForcePassword(password, maxLength)

	// Se a senha for encontrada, retornar para o servidor principal
	if result != "" {
		conn.Write([]byte("Senha encontrada: " + result + "\n"))
	} else {
		conn.Write([]byte("Senha não encontrada\n"))
	}
}

func main() {
	l, err := net.Listen("tcp", ":9002")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor de processamento:", err)
		return
	}
	defer l.Close()

	fmt.Println("Servidor de processamento iniciado na porta 9002")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleConnection(conn)
	}
}
