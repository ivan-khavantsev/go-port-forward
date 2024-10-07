// main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func handleConnection(clientConn net.Conn, remote string) {
	defer clientConn.Close()

	// Подключаемся к удаленному серверу
	remoteConn, err := net.Dial("tcp", remote)
	if err != nil {
		log.Printf("Ошибка при подключении к удаленному серверу %s: %v", remote, err)
		return
	}
	defer remoteConn.Close()

	// Перенаправляем данные между клиентом и сервером
	go io.Copy(remoteConn, clientConn)
	io.Copy(clientConn, remoteConn)
}

func startPortForwarder(local string, remote string) {
	listener, err := net.Listen("tcp", local)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер на порту %s: %v", local, err)
	}
	defer listener.Close()

	log.Printf("Портфорвардер запущен на порту %s, перенаправление на %s", local, remote)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка при принятии подключения: %v", err)
			continue
		}

		go handleConnection(clientConn, remote)
	}
}

func main() {
	// Параметры командной строки
	listenPort := flag.String("local", ":8080", "Хост и Порт для прослушивания")
	remoteHost := flag.String("remote", "127.0.0.1:80", "Удаленный хост и порт для перенаправления")

	flag.Parse()

	if *listenPort == "" || *remoteHost == "" {
		fmt.Println("Usage: port-forward --local <[host]:port> --remote <host:port>")
		os.Exit(1)
	}

	// Запускаем портфорвардер
	startPortForwarder(*listenPort, *remoteHost)
}
