package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/timestreamerror/capstone/internal/database"
)

func importCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}
	fmt.Println("*2*")
	for _, record := range records {
		fmt.Println(record)
	}
	return records, nil
}

func (apiCfg *APIConfig) handleImportCSV(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	filename = "./csv/" + filename
	log.Println(filename)
	records, err := importCSV(filename)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, record := range records {
		dbParams := database.PutQuoteParams{}
		dbParams.ID, err = uuid.NewUUID()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		dbParams.CreatedAt = time.Now()
		dbParams.UpdatedAt = time.Now()
		dbParams.Quote = record[0]
		if record[1] != "" {
			dbParams.Author.String = record[1]
			dbParams.Author.Valid = true
		} else {
			dbParams.Author.Valid = false
		}
		apiCfg.dbQueries.PutQuote(r.Context(), dbParams)
	}

}
