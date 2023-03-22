package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type USDBRL struct {
	Bid string `json:"bid"`
}

func main() {
	timeOut := time.Millisecond * 1000

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	dollar, err := getDollarQuotation(ctx)
	if err != nil {
		panic(err)
	}

	err = saveDollarQuotation(dollar)
	if err != nil {
		panic(err)
	}
}

func getDollarQuotation(ctx context.Context) (*USDBRL, error) {
	client := http.Client{}

	request, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:7000/quotation", nil)
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

	var usdbrl USDBRL
	err = json.Unmarshal(data, &usdbrl)
	if err != nil {
		return nil, err
	}

	return &usdbrl, nil
}

func saveDollarQuotation(usdbrl *USDBRL) error {
	file, err := os.Create("Quotation.txt")
	if err != nil {
		return err
	}

	dollar := fmt.Sprintf("Dollar: %s", usdbrl.Bid)

	size, err := file.Write([]byte(dollar))
	if err != nil {
		return err
	}

	fmt.Printf("File created with size %d\n", size)

	return nil
}
