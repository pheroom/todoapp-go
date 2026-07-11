package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
)

type PayInfo struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

var mtx = sync.Mutex{}
var money = 1000
var payList = make(map[int]PayInfo)

func payHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var payment PayInfo
	if err := json.Unmarshal(httpRequestBody, &payment); err != nil {
		fmt.Println("Error:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//httpRequestBodyStr := string(httpRequestBody)
	//httpRequestBodyInt, err := strconv.Atoi(httpRequestBodyStr)
	//if err != nil {
	//	fmt.Println("Error:", err.Error())
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	fmt.Println("get body:", payment)

	mtx.Lock()
	if money-payment.Amount < 0 {
		fmt.Println("Error: malo money")
		http.Error(w, "malo money", http.StatusBadRequest)
		return
	} else {
		money -= payment.Amount
		id := rand.Int()
		//newPayInfo := PayInfo{Amount: httpRequestBodyInt, Description: "New Payment. id - " + string(id)}
		payList[id] = payment
		fmt.Println("new payment with id:", id, "and with amount:", payment.Amount)
	}
	mtx.Unlock()
}

func getPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	mtx.Lock()
	jsonList, err := json.MarshalIndent(payList, "", "  ")
	if err != nil {
		fmt.Println("Error:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(jsonList)
	mtx.Unlock()
}

func main() {
	urlExample := "postgres://admin:admin@localhost:5432/testapidb"
	//os.Getenv("DATABASE_URL")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	if err := conn.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		panic(err)
	}

	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/get-pay", getPaymentsHandler)
	fmt.Println("Server starting")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		fmt.Println("listen err:", err)
	}
}
