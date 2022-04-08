package util

import "mime/multipart"

type FileData struct {
	FileName string
	File multipart.File
	Size int64
}

func MultipartFileConverter(multipartFiles map[string][]*multipart.FileHeader) (output map[string][]*FileData,err error) {
	output = map[string][]*FileData{}

	for k,v := range multipartFiles {
		var temp []*FileData
		for _,d := range v {
			data := &FileData{
				FileName: d.Filename,
				Size: d.Size,
			}
			data.File,err = d.Open()
			if err != nil {
				return
			}
			temp = append(temp,data)
		}
		output[k] = temp
	}
	return
}
