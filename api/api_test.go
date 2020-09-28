package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/tkrouty/avitojob-trainee-task/db"
	"github.com/tkrouty/avitojob-trainee-task/models"
	"github.com/tkrouty/avitojob-trainee-task/test"
)

var (
	tester             = test.InitTester()
	TestDBWrapper      db.DBWrapper
	TestFinanceManager FinanceManager
	TestFinanceAPI     FinanceAPI
	TestServer         *httptest.Server
)

type BalanceResponse struct {
	UserID   string  `json:"user_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type HistoryResponse struct {
	UserID             string               `json:"user_id"`
	TransactionHistory []models.Transaction `json:"transaction_history"`
}

func BeforeTest() {
	tester.ConnectToDB()
	tester.FlushTables()
	TestDBWrapper = db.DBWrapper{
		Conn:   tester.DB,
		Logger: tester.Logger,
	}
	TestFinanceManager = FinanceManager{
		DB:    TestDBWrapper,
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
	TestFinanceAPI = FinanceAPI{Manager: TestFinanceManager}
	TestServer = httptest.NewServer(setupTestRouter())
}

func AfterTest() {
	tester.CloseDBConnection()
}

func setupTestRouter() *gin.Engine {
	r := gin.Default()

	r.HandleMethodNotAllowed = true
	// list of all routes
	r.POST("/edit_balance/:UserID", TestFinanceAPI.EditBalance)
	r.POST("/transfer", TestFinanceAPI.Transfer)
	r.GET("/show_balance/:UserID", TestFinanceAPI.ShowBalance)
	r.GET("/show_history/:UserID", TestFinanceAPI.ShowHistory)

	return r
}

func TestShowBalance(t *testing.T) {
	BeforeTest()
	// testing that initially there is no info about user
	resp, err := http.Get(fmt.Sprintf("%s/show_balance/1", TestServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

	// adding a new user
	transaction := models.Transaction{TargetID: "1", Sum: 0.48}
	TestFinanceManager.makeTransaction(&transaction)

	// checking that now server gives info about the user
	resp, err = http.Get(fmt.Sprintf("%s/show_balance/1", TestServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	var r BalanceResponse

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	tester.Logger.Printf("got response %v\n", r)

	if r.UserID != "1" {
		t.Fatalf("Expected user_id=1, got %v", r.UserID)
	}
	if r.Balance != 0.48 {
		t.Fatalf("Expected balance=0.48, got %v", r.Balance)
	}
	if r.Currency != "RUB" {
		t.Fatalf("Expected default currency (RUB), got %v", r.Currency)
	}

	resp, err = http.Get(fmt.Sprintf("%s/show_balance/1?currency=USD", TestServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	tester.Logger.Printf("got response %v\n", r)

	exchangeRate, err := getExchangeRatebyHTTP("USD")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if r.UserID != "1" {
		t.Fatalf("Expected user_id=1, got %v", r.UserID)
	}
	if r.Balance != 0.48*exchangeRate {
		t.Fatalf("Expected balance=%v, got %v", 0.48*exchangeRate, r.Balance)
	}
	if r.Currency != "USD" {
		t.Fatalf("Expected default currency (RUB), got %v", r.Currency)
	}

	AfterTest()
}

func TestEditBalance(t *testing.T) {
	BeforeTest()
	// making a transaction
	transaction := models.Transaction{TargetID: "1", Sum: 0.48}
	TestFinanceManager.makeTransaction(&transaction)

	resp, err := postRequest("/edit_balance/1", map[string]interface{}{"sum": "mama"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

	resp, err = postRequest("/edit_balance/1", map[string]interface{}{"sum": 42})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	u := models.User{UserID: "1"}
	balance, err := TestFinanceManager.getBalance(&u, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if balance != 42.48 {
		t.Fatalf("Expected balance=42.48, got %v", balance)
	}

	resp, err = postRequest("/edit_balance/1", map[string]interface{}{"sum": -42})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	balance, err = TestFinanceManager.getBalance(&u, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if balance != 0.48 {
		t.Fatalf("Expected balance=0.48, got %v", balance)
	}

	resp, err = postRequest("/edit_balance/1", map[string]interface{}{"sum": -42})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

	AfterTest()
}

func TestTransfer(t *testing.T) {
	BeforeTest()

	resp, err := postRequest("/transfer", map[string]interface{}{
		"sum":       "volodya",
		"source_id": "1",
		"target_id": "2",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

	transaction := models.Transaction{TargetID: "1", Sum: 0.48}
	TestFinanceManager.makeTransaction(&transaction)

	// checking that we cannot use negative transfer sum
	resp, err = postRequest("/transfer", map[string]interface{}{
		"sum":       -0.42,
		"source_id": "1",
		"target_id": "2",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

	resp, err = postRequest("/transfer", map[string]interface{}{
		"sum":       0.42,
		"source_id": "1",
		"target_id": "2",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	source_user := models.User{UserID: "1"}
	target_user := models.User{UserID: "2"}

	source_balance, err := TestFinanceManager.getBalance(&source_user, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	target_balance, err := TestFinanceManager.getBalance(&target_user, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if source_balance != 0.06 {
		t.Fatalf("Expected source_balance=0.06, got %v", source_balance)
	}
	if target_balance != 0.42 {
		t.Fatalf("Expected target_balance=0.42, got %v", source_balance)
	}

	resp, err = postRequest("/transfer", map[string]interface{}{
		"sum":       0.42,
		"source_id": "1",
		"target_id": "2",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

	AfterTest()
}

func TestShowHistory(t *testing.T) {
	BeforeTest()

	resp, err := http.Get(fmt.Sprintf("%s/show_history/1", TestServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

	firstTransaction := models.Transaction{TargetID: "1", Sum: 0.4}
	TestFinanceManager.makeTransaction(&firstTransaction)

	secondTransaction := models.Transaction{TargetID: "1", Sum: 0.3}
	TestFinanceManager.makeTransaction(&secondTransaction)

	thirdTransaction := models.Transaction{SourceID: "1", TargetID: "2", Sum: 0.6}
	TestFinanceManager.makeTransaction(&thirdTransaction)

	resp, err = http.Get(fmt.Sprintf("%s/show_history/1", TestServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
	var r HistoryResponse

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	tester.Logger.Printf("got response %v\n", r)

	if r.TransactionHistory[0].Sum != thirdTransaction.Sum {
		t.Fatalf("Expected first_transaction sum=0.6, got %v", r.TransactionHistory[0].Sum)
	}
	if r.TransactionHistory[1].Sum != firstTransaction.Sum {
		t.Fatalf("Expected first_transaction sum=0.4, got %v", r.TransactionHistory[0].Sum)
	}
	if r.TransactionHistory[2].Sum != secondTransaction.Sum {
		t.Fatalf("Expected first_transaction sum=0.3, got %v", r.TransactionHistory[0].Sum)
	}

	AfterTest()

}

func postRequest(route string, payload map[string]interface{}) (*http.Response, error) {
	jsonValue, _ := json.Marshal(payload)
	resp, err := http.Post(TestServer.URL+route,
		"application/json", bytes.NewBuffer(jsonValue))
	return resp, err
}
