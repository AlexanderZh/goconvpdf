package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	//ContentType string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(mux.CORSMethodMiddleware(router))

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			//Headers("Content-Type", route.ContentType).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API for document conversion")
}

func Convert(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		log.Println("Error:" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tmp_filename := uuid.NewString()
	f, err := os.OpenFile(tmp_filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		_, err = io.Copy(f, part)
		if err != nil {
			break
		}
	}

	fileName := r.Header.Get("fileName")
	if err != nil {
		panic(err)
	}
	fmt.Printf("File name %s\n", fileName)
	w.Header().Set("Content-Type", "multipart/form-data")
	w.WriteHeader(http.StatusOK)
	arg0 := "lowriter"
	arg1 := "--invisible" //This command is optional, it will help to disable the splash screen of LibreOffice.
	arg2 := "--convert-to"
	arg3 := "pdf:writer_pdf_Export"
	path := tmp_filename
	nout, err := exec.Command(arg0, arg1, arg2, arg3, path).Output()
	if err != nil {
		log.Println("Error:" + err.Error())
	} else {
		log.Println("Success:" + string(nout))
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+tmp_filename)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	f_converted, _ := os.OpenFile(tmp_filename+".pdf", os.O_RDONLY, 0666)

	io.Copy(w, f_converted)
	defer f_converted.Close()
	defer os.Remove(tmp_filename)
	defer os.Remove(tmp_filename + ".pdf")

}

func calcBytes(b *[]byte, n uint64, c chan uint64) {
	var i uint64 = 0
	var sum uint64 = 0
	for i = 0; i < n; i++ {
		sum += (uint64)((*b)[i])
	}
	c <- sum
}

func testPartProcessing(w http.ResponseWriter, r *http.Request) {
	partReader, err := r.MultipartReader()
	if err != nil {
		log.Println("Error:" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	buf := make([]byte, 1024*1024)
	var sum uint64 = 0
	c := make(chan uint64)

	for {
		part, err := partReader.NextPart()
		if err == io.EOF {
			break
		}
		var n int
		chunks := 0
		//map
		for {
			n, err = part.Read(buf)
			if err == io.EOF {
				break
			}
			go calcBytes(&buf, (uint64)(n), c)
			chunks += 1
			if chunks%128 == 0 {
				//reduce
				for i := 0; i < chunks; i++ {
					sum += <-c
				}
				chunks = 0
			}
		}
		for i := 0; i < chunks; i++ {
			sum += <-c
		}
		chunks = 0
		go calcBytes(&buf, (uint64)(n), c)
		sum += <-c
		fmt.Println(sum)
	}
}

var routes = Routes{
	//root does nothing
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Convert",
		"POST",
		"/",
		Convert,
	},
	Route{
		"TestMultipart",
		"POST",
		"/test",
		testPartProcessing,
	},
}
