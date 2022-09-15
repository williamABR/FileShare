package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const uploadPath = "./uploadsFile"

type File struct {
	name  string
	users map[string]bool
}
type Info struct {
	idUser  string
	channel string
}

var userSubscription = make(map[string]map[string]bool)        // Usuarios que estan suscritos a los canales
var channelsFile = make(map[string]map[string]map[string]bool) // Los usuarios a los que se les a compartido los archivos

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// truncated for brevity

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend

	file, fileHeader, err := r.FormFile("myFile")
	r.ParseForm()
	productsSelected := r.Form["channel"][0]
	productId := r.Form["idUser"][0]
	subscriptionUser(productsSelected, productId)

	fmt.Println("--------------------")
	for k, v := range userSubscription {
		fmt.Println(k, "value is", v)
	}
	if err != nil {
		fmt.Printf("ahi va")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./uploadsFile", os.ModePerm)
	if err != nil {
		fmt.Printf("ahi va" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	//dst, err := os.Create(fmt.Sprintf("./uploadsFile/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	name := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileHeader.Filename)
	dst, err := os.Create(fmt.Sprintf("./uploadsFile/%s", name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, foundFile := channelsFile[productsSelected]
	if !foundFile {
		channelsFile[productsSelected] = make(map[string]map[string]bool)
	}
	v2, found2 := channelsFile[productsSelected]
	if found2 {
		v2[name] = make(map[string]bool)
		v3, found3 := v2[name]
		if found3 {
			fmt.Println(productId)
			v3[productId] = true
		}
	}
	fmt.Println(channelsFile)
	fmt.Fprintf(w, name)
}
func Index(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		fmt.Println("id is missing in parameters")
	}
	fmt.Println(id)

	fileBytes, err := ioutil.ReadFile("uploadsFile/" + id)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)

	return
}
func checkHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	r.ParseForm()
	productsSelected := r.Form["sizeFile"]
	fmt.Println(productsSelected)

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		fmt.Println("id is missing in parameters")
	}
	//fmt.Println(id)

	info := map[string]interface{}{}
	for k, v := range r.PostForm {
		fmt.Println(k, "value is", v)
		if err := json.Unmarshal([]byte(k), &info); err != nil {
			panic(err)
		}
	}

	//channel := info["channel"]
	c := info["isFile"].(bool)

	//buscar archivos sin enviar
	filesAux := make(map[string]string)
	for k, v := range userSubscription {
		//fmt.Printf("key[%s] value[%s]\n", k, v)
		_, found := v[id]
		if found {
			for k1, v1 := range channelsFile[k] { //recoriiendo todos los archivos
				v2 := v1[id]
				fmt.Println("archivo:", k1, "-Usuarios:", v1)
				if c {
					filesAux[k1] = k
				} else if !v2 {
					filesAux[k1] = k
					v1[id] = true
				}
			}
		}
	}

	fmt.Println(filesAux)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(filesAux)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

	//fmt.Fprintf(w, filesAux)
	return
}
func subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// truncated for brevity

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)

	info := map[string]interface{}{}
	for k, v := range r.PostForm {
		fmt.Println(k, "value is", v)
		if err := json.Unmarshal([]byte(k), &info); err != nil {
			panic(err)
		}
	}

	//channel := info["channel"]
	c := info["channel"].(string)
	u := info["idUser"].(string)
	subscriptionUser(c, u)

	fmt.Println(userSubscription)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Subscrito con exito")
	return
}
func subscriptionUser(productsSelected string, productId string) {
	v, found := userSubscription[productsSelected]
	if found {
		v[productId] = true
	} else {
		userSubscription[productsSelected] = make(map[string]bool)
		v, found := userSubscription[productsSelected]
		if found {
			v[productId] = true
		}
	}
}

func setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadHandler)
	r.HandleFunc("/subscription", subscriptionHandler)
	r.HandleFunc("/checkUpdate/{id}", checkHandler)
	r.HandleFunc("/file/{id}", Index)

	http.ListenAndServe(":8080", r)
}

func main() {
	fmt.Println(userSubscription)
	setupRoutes()
}
