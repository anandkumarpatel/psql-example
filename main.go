package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/anandkumarpatel/main/databases"
	"github.com/anandkumarpatel/main/routes"
	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Print("DB_URL not defined\n")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := databases.NewPostgres(ctx, dbURL)
	if err != nil {
		fmt.Printf("Error with db %v\n", err)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error with listing %v\n", err)
		os.Exit(1)
	}
	r := mux.NewRouter()

	routes.AddServiceRoutes(r, db)

	http.Handle("/", r)
	fmt.Printf("server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server %v\n", err)
		os.Exit(1)
	}
}
