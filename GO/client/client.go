package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Pega os argumentos da linha de comando para host:porta
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Coloque as informações de host:port.")
		return
	}

	// Conecta ao servidor principal
	c, err := net.Dial("tcp", arguments[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	fmt.Println("Digite a senha que deseja descobrir.")
	for {
		// Leitura do teclado (entrada do usuário)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Senha: ")
		text, _ := reader.ReadString('\n')

		// Envia o texto para o servidor principal
		fmt.Fprintf(c, text+"\n")

		// Recebe a resposta do servidor principal
		message, _ := bufio.NewReader(c).ReadString('\n')

		// Exibe a resposta
		fmt.Print("Resposta: " + message)

		// Se o cliente enviar "EXIT", encerra a conexão
		if strings.ToUpper(strings.TrimSpace(string(text))) == "EXIT" {
			fmt.Println("Encerrando o cliente TCP...")
			return
		}

		// Verifica se a senha foi encontrada
		if strings.Contains(message, "found the password") {
			fmt.Println("A senha foi encontrada.")
			return
		}
	}
}
