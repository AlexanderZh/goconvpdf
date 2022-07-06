package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

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
	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("fileName")
	if err != nil {
		panic(err)
	}
	fileName := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", fileName[0])
	defer file.Close()
	w.Header().Set("Content-Type", "multipart/form-data")
	w.WriteHeader(http.StatusOK)
	tmp_filename := uuid.NewString()
	f, err := os.OpenFile(tmp_filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	io.Copy(f, file)
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
}
