package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	res, err := http.Get("http://kambing.ui.ac.id/iso/debian/10.1.0-live/amd64/iso-hybrid/debian-live-10.1.0-amd64-cinnamon.iso")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(res.ContentLength)
	//buffer := make([]byte,4*1024*1024)
	//for {
	//	_,err := res.Body.Read(buffer)
	//	if err == io.EOF {
	//		break
	//	}
	//}
}
