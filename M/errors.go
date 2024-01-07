package M

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	ErrNotFound = errors.New("models: resourses not found ")
	ErrEmailTaken =errors.New("models: email address is already in used ")
)

type FileError struct{
	Issue string
}

func(f FileError) Error()string{
	return fmt.Sprintf("invalid file : %v",f.Issue)
}

func checkContentType(r io.ReadSeeker,allwoedType []string)error{
	testBytes:=make([]byte,512)
	_,err:=r.Read(testBytes)
	if err !=nil{
		return fmt.Errorf("checking contentType :%w",err)
	}
	_,err=r.Seek(0,0)
	if err !=nil{
		return fmt.Errorf("checking contentType :%w",err)
	}
	t:=http.DetectContentType(testBytes)
	for _,aT:= range allwoedType{
		if t==strings.ToLower(aT){
			return nil
		}
	}
	return FileError{Issue:fmt.Sprintf("invalid content Type :%s",t)}
}

func checkExtension(filename string,alloweTypes []string)error{
		if !hasExtension(filename,alloweTypes){
			return FileError{Issue: fmt.Sprintf("invalid file extension :%s",filepath.Ext(filename))}
		}
		return nil
	}
