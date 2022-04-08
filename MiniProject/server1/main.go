package main

import (
	"MiniProject/downloader"
	"MiniProject/server1/grpc"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"path/filepath"
)

func main() {

	staticDir,_ := filepath.Abs("static")
	server1Dir,_ := filepath.Abs("server1")
	static := http.FileServer(http.Dir(staticDir))

	downloaderConn,downloaderClient := grpc.InitClient()
	defer downloaderConn.Close()

	proxy := &httputil.ReverseProxy{Director: func(request *http.Request) {
		request.Header.Add("Upgrade",request.Header.Get("Upgrade"))
		request.Header.Add("Connection","Upgrade")
		request.URL.Scheme = "http"
		request.URL.Host = "localhost:15674"
	}}

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		proxy.ServeHTTP(writer,request)
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		parsedTemplate,_ := template.ParseFiles(filepath.Join(server1Dir,"html","main.html"))
		err := parsedTemplate.Execute(writer,nil)
		if err != nil {
			_, _ = writer.Write([]byte("error"))
			return
		}
	})

	http.HandleFunc("/download-request", func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			_,_ = writer.Write([]byte(err.Error()+ " parse form"))
			return
		}
		link := request.PostFormValue("link")
		if link == "" {
			_,_ = writer.Write([]byte("link can't be empty"))
			return
		}

		//GRPC Request Side

		requestData := &downloader.DownloadRequest{Url: link}

		resp,err := downloaderClient.Download(context.Background(),requestData)
		if err != nil {
			_,_ = writer.Write([]byte(err.Error()+" grpc call"))
			return
		}

		id := resp.UniqId

		dataTemplate := struct {
			Exchange string
			WSLink string
		}{
			Exchange: fmt.Sprintf(`/exchange/%s/%s`,"1806235832",id),
			WSLink: "ws://${window.location.host}/ws",
		}

		parsedTemplate,_ := template.ParseFiles(filepath.Join(server1Dir,"html","progress.html"))
		err = parsedTemplate.Execute(writer,dataTemplate)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()+" template"))
			return
		}
	})

	http.Handle("/download/",http.StripPrefix("/download",static))

	log.Printf("Serve in port 8081\n")
	log.Fatal(http.ListenAndServe(":8081",nil))

}
