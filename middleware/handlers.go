package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/souravdev-eng/go-postgres-stock/models"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to postgres...")
	return db
}

func HandelError(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func CreateNewStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	HandelError("Unable to decode the request body %v", err)

	insertId := insertStock(stock)

	res := response{
		ID:      insertId,
		Message: "Stock created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	HandelError("Unable to convert the string into int. %v", err)
	stock, err := getStock(int64(id))

	HandelError("Unable to get stock. %v", err)
	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()

	HandelError("Unable to get all stocks. %v", err)
	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	HandelError("Unable to convert the string into int. %v", err)

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)

	HandelError("Unable to decode the request body. %v", err)
	updatedRows := updateStock(int64(id), stock)
	msg := fmt.Sprintf("Stock updated successfully. Total rows/records affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	HandelError("Unable to convert the string into int. %v", err)

	deleteRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock deleted successfully. Total rows/records deleted %v", deleteRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func insertStock(stock models.Stock) int64 {
	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	HandelError("Unable to execute the query %v", err)

	fmt.Printf("Inserted a single record %v", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`

	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("Now rows were returned!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the row %v", err)
	}
	return stock, err
}

func getAllStocks() ([]models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)

	HandelError("Unable to execute query %v", err)
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock

		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		HandelError("Unable to scan the row %v", err)
		stocks = append(stocks, stock)
	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`
	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	HandelError("Unable to execute the query %v", err)

	rowsAffected, err := res.RowsAffected()
	HandelError("Error while checking the affected rows %v", err)

	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id)
	HandelError("Unable to execute the query %v", err)

	rowsAffected, err := res.RowsAffected()
	HandelError("Error while checking the affected rows %v", err)

	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}
