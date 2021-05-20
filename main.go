package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	MYAPIKEY       string //global variable for the client
	store          = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	mutex          sync.Mutex
	db             *sql.DB
	tpl            *template.Template
	protocol       string
	domain         string
	productionFlag bool = false
)

type MerchantFoodInfo struct {
	FoodName string
	Price    string
}

type merchantInformation struct {
	MerchantName string `json: MerchantName`
	Address      string `json: Address`
}

type foodProperties struct {
	FoodName string `json: FoodName`
	Price    string `json: Price`
}

// var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {

	flag.StringVar(&domain, "domain", "localhost", "domain name to request your certificate")
	flag.BoolVar(&productionFlag, "productionFlag", false, "if true, we'll send https and custom domain requests")
	flag.Parse()

	if productionFlag {
		fmt.Println("\nProductin flag activated. Code to be used in Production only")
		protocol = "https://"
	} else {
		protocol = "http://"
	}

	router := mux.NewRouter()
	router.HandleFunc("/", userprofile)
	router.HandleFunc("/additems", additems)
	router.HandleFunc("/userprofile", userprofile)
	router.HandleFunc("/checkrecords", checkRecords)
	router.HandleFunc("/editrecords", editRecords)
	router.HandleFunc("/edititem", editItem)
	router.HandleFunc("/deleteitems", deleteItems)
	router.HandleFunc("/assets/rs3.png", rs3)
	router.HandleFunc("/merchantsetup", merchantsetup)
	router.Handle("/favicon.ico", http.NotFoundHandler()) //NotFoundHandler returns a simple request handler that replies to each request with a “404 page not found” reply.
	// fmt.Println("Client HTTP server Listeing at :8080")
	fmt.Println("domain name: ", domain)
	http.ListenAndServe(":80", router) //launches HTTP server

} // close main functioin

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
} //end function check
