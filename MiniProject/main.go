package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	newFile,_ := os.Create("test.log")
	resp,_ := http.Get("http://kambing.ui.ac.id/iso/debian/10.1.0-live/amd64/iso-hybrid/debian-live-10.1.0-amd64-kde.log")
	buffer := make([]byte,32*1024)
	body := resp.Body
	for {
		n,err := body.Read(buffer)
		if err == io.EOF || n == 0 {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		_, err = newFile.Write(buffer[0:n])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	newFile.Close()
}
