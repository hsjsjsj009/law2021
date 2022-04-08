package main

import (
	"MiniLab2/compression"
	"MiniLab2/server1/grpc"
	"context"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {

	compressionConn,compressionClient := grpc.InitClient()
	defer compressionConn.Close()

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
		parsedTemplate,_ := template.ParseFiles("./html/main.html")
		err := parsedTemplate.Execute(writer,nil)
		if err != nil {
			_, _ = writer.Write([]byte("error"))
			return
		}
	})

	http.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		uid := fmt.Sprintf("%s",uuid.New())
		file,header,err := request.FormFile("file")
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()+" formfile"))
			return
		}

		//GRPC Request Side

		fileBytes,err := ioutil.ReadAll(file)
		if err != nil {
			_,_ = writer.Write([]byte(err.Error()+" read file"))
			return
		}

		requestData := &compression.CompressionRequest{
			RoutingKey: uid,
			FileName: header.Filename,
			FileBytes: fileBytes,
		}

		resp,err := compressionClient.CompressFile(context.Background(),requestData)
		if err != nil {
			_,_ = writer.Write([]byte(err.Error()+" grpc call"))
			return
		}

		if resp.Success == compression.SuccessStatus_SUCCESS_STATUS_FAIL {
			_,_ = writer.Write([]byte(resp.ErrorMessage+" response grpc"))
			return
		}

		//HTTP/1.1 Request Side

		//body := new(bytes.Buffer)
		//multipartWriter := multipart.NewWriter(body)
		//part,err := multipartWriter.CreateFormFile("file",header.Filename)
		//if err != nil {
		//	_, _ = writer.Write([]byte(err.Error()+" multipart writer"))
		//	return
		//}
		//_,err = io.Copy(part,file)
		//if err != nil {
		//	_, _ = writer.Write([]byte(err.Error()+" multipart writing"))
		//	return
		//}
		//err = multipartWriter.Close()
		//if err != nil {
		//	_, _ = writer.Write([]byte(err.Error()+" multipart close"))
		//	return
		//}

		//req,err := http.NewRequest(http.MethodPost,"http://localhost:8080/",body)
		//if err != nil {
		//	_, _ = writer.Write([]byte(err.Error()+" req"))
		//	return
		//}
		//req.Header.Add("X-ROUTING-KEY",uid)
		//req.Header.Add("Content-Type",multipartWriter.FormDataContentType())
		//resp,err := http.DefaultClient.Do(req)
		//if err != nil {
		//	_, _ = writer.Write([]byte(err.Error()+" call req"))
		//	return
		//}
		//jsonData := map[string]string{}
		//jsonDecoder := json.NewDecoder(resp.Body)
		//err = jsonDecoder.Decode(&jsonData)
		//if err != nil {
		//	_, _ = io.Copy(writer,resp.Body)
		//	return
		//}
		//if v,ok := jsonData["message"];ok {
		//	_, _ = writer.Write([]byte(v))
		//	return
		//}

		dataTemplate := struct {
			Exchange string
			WSLink string
		}{
			Exchange: fmt.Sprintf(`/exchange/%s/%s`,"1806235832",uid),
			WSLink: "ws://${window.location.host}/ws",
		}

		parsedTemplate,_ := template.ParseFiles("./html/progress.html")
		err = parsedTemplate.Execute(writer,dataTemplate)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()+" template"))
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8081",nil))

}
