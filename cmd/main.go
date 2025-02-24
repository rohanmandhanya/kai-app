package main

import (
	"fmt"
	"log"
	"net/http"

	"kai-app/api/controller"
	"kai-app/arch/database"

	"github.com/gorilla/mux"
)

func main() {

	// Initialize the database
	db, err := database.InitializeDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer database.DisconnectDB(db)

	router := mux.NewRouter()
	router.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	router.HandleFunc(("/query"), controller.QueryHandler).Methods("POST")
	router.HandleFunc(("/scan"), controller.ScanHandler).Methods("POST")

	fmt.Println(`
  \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
  \\   *** SERVER IS ALIVE ***  ||  Your API is Ready to Handle All   //
   \\   Ready to tackle requests  ||  With Speed and Precision          //
    \\   Go Beyond the Limits    ||  Powered by Golang and Passion     //
     \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
     
        *** API is now live at: http://localhost:8080 ***
        All systems go! 🚀

  Connect and make your calls, we're waiting for you! 🌟
`)

	fmt.Println("Server running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}
