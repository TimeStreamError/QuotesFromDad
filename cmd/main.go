package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/timestreamerror/capstone/internal/database"

	_ "modernc.org/sqlite"
)

type APIConfig struct {
	DB                 *sql.DB
	dbQueries          *database.Queries
	tokenSecret        string
	loginAuthorization string
}

func main() {
	// load .env
	godotenv.Load()
	dbURL := os.Getenv("DB_FILE")
	if dbURL == "" {
		log.Fatal("Error: could not read env file")
	}

	// set up the apiConfig to pass around
	apiCfg := APIConfig{}
	db, err := sql.Open("sqlite", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	apiCfg.DB = db
	apiCfg.dbQueries = database.New(db)
	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatal("Error: could not read enc file")
	}
	loginAuth := os.Getenv("LOGIN_AUTHORIZATION")
	if loginAuth == "" {
		log.Fatal("Error: could not read enc file")
	}
	apiCfg.tokenSecret = tokenSecret
	apiCfg.loginAuthorization = loginAuth

	// basic http server
	r := chi.NewRouter()
	httpServer := http.Server{}

	// the various path handlers
	r.Handle("/*", http.FileServer(http.Dir(".")))
	r.Post("/api/quotes", apiCfg.handlePutQuote)                 // put quote requires login
	r.Post("/api/import/{filename}", apiCfg.handleImportCSV)     // import requires login
	r.Get("/api/quotes", apiCfg.handleGetAllQuotes)              // all quotes requires login
	r.Get("/random", apiCfg.handleGetRandom)                     // no login required, will use UUID id
	r.Post("/admin/login", apiCfg.handlerLogin)                  // no login required, but user needs to be in db
	r.Post("/admin/signup", apiCfg.handleSignup)                 // make a new user, should require master pin
	r.Get("/api/refresh", apiCfg.handleRefreshTokenCheck)        // called when the client gets a 401 status
	r.Get("/api/quotes/{quote-id}", apiCfg.handleEditQuote)      // to populate edit quote screen
	r.Put("/api/quotes/{quote-id}", apiCfg.handlePutEditedQuote) // updates the modified quote
	r.Get("/api/quotes/reset", apiCfg.handleQuotesReset)         // deletes all quotes from db
	r.Post("/api/tags", apiCfg.handleAddTag)
	// boot up the server
	httpServer.Addr = ":8080"
	httpServer.Handler = r
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
