package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const quotationURI = "/quotation"
const dollarQuotationURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type USDBRL struct {
	Code   string `json:"code"`
	CodeIn string `json:"codein"`
	Name   string `json:"name"`
	Bid    string `json:"bid"`
}

type DollarQuotation struct {
	USDBRL USDBRL `json:"USDBRL"`
}

func main() {
	http.HandleFunc(quotationURI, handleQuotation)
	http.ListenAndServe(":7000", nil)
}

func handleQuotation(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	if request.URL.Path != quotationURI {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := getDollarQuotation(ctx)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	}

	err = saveDataIntoDatabase(ctx, data)
	if err != nil {
		panic(err)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	json.NewEncoder(response).Encode(data)
}

func getDollarQuotation(ctx context.Context) (*USDBRL, error) {
	client := http.Client{}

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*800)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", dollarQuotationURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var dollarQuotation DollarQuotation
	err = json.Unmarshal(data, &dollarQuotation)
	if err != nil {
		return nil, err
	}

	return &dollarQuotation.USDBRL, nil
}

func saveDataIntoDatabase(ctx context.Context, data *USDBRL) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	db, err := sql.Open("sqlite3", "./database.s3db")
	if err != nil {
		return err
	}

	defer db.Close()

	err = dbMigrate(ctx, db)
	if err != nil {
		return err
	}

	err = save(ctx, db, data)
	if err != nil {
		return err
	}

	return nil
}

func dbMigrate(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)

	defer cancel()

	_, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS USDBRL (code text, codein text, name text, bid string)")
	if err != nil {
		return err
	}

	return nil
}

func save(ctx context.Context, db *sql.DB, data *USDBRL) error {
	stmt, err := db.Prepare("INSERT INTO USDBRL(code, codein, name, bid) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, data.Code, data.CodeIn, data.Name, data.Bid)
	if err != nil {
		return err
	}

	return nil
}
