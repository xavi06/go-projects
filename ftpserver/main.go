package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type client struct {
	pwd string
}

func handleConn(c net.Conn) {
	client := new(client)
	pwdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		c.Close()
	}
	client.pwd = pwdir
	defer c.Close()
	fmt.Fprintln(c, "220 FTP Server ready.")
	scan := bufio.NewScanner(c)
	for scan.Scan() {
		if err := scan.Err(); err != nil {
			log.Print(err)
			continue
		}
		cmd := scan.Text()
		switch cmd {
		case "close":
			c.Close()
		case "ls":
			client.listDir(c)
		case "pwd":
			client.getpwd(c)
		case "cd":
			//client.pwd = pwdir + "/1001"
			client.changeDir(c)
		default:
			fmt.Fprintln(c, "unknown command!")
		}
	}
}

func (cl *client) listDir(c net.Conn) {
	dir, err := os.Open(cl.pwd)
	if err != nil {
		log.Fatal(err)
		c.Close()
	}
	infos, err := dir.Readdir(0)
	if err != nil {
		log.Fatal(err)
		c.Close()
	}
	fmt.Fprintln(c, "150 Opening ASCII mode data connection for file list")
	for _, fi := range infos {
		fmt.Fprint(c, fi.Name()+"\n")
	}
	fmt.Fprintln(c, "226 Transfer complete.")

}

func (cl *client) getpwd(c net.Conn) {
	fmt.Fprintf(c, "257 \"%s\" is current directory.\n", cl.pwd)
}

func (cl *client) changeDir(c net.Conn) {
	cl.pwd = cl.pwd + "/1001"
	fmt.Fprintf(c, "change dir to %s\n", cl.pwd)
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("listen localhost:8000\n")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)

	}
}
