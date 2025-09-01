package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Employee struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	InStock     bool    `json:"in_stock"`
	Description string  `json:"description"`
}

var db *gorm.DB

func createTableIfNotExists() {
	err := db.AutoMigrate(&Employee{}, &Product{})
	if err != nil {
		log.Fatalf("Erro ao migrar tabelas: %v", err)
	}
}

func addEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if emp.Name == "" && emp.Address == "" {
		http.Error(w, "Nome ou endereço obrigatório", http.StatusBadRequest)
		return
	}
	if err := db.Create(&emp).Error; err != nil {
		http.Error(w, "Erro ao inserir", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

func listEmployees(w http.ResponseWriter, r *http.Request) {
	var emps []Employee
	if err := db.Find(&emps).Error; err != nil {
		http.Error(w, "Erro ao buscar", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(emps)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if p.Name == "" {
		http.Error(w, "Nome do produto obrigatório", http.StatusBadRequest)
		return
	}
	if err := db.Create(&p).Error; err != nil {
		http.Error(w, "Erro ao inserir produto", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	var prods []Product
	if err := db.Find(&prods).Error; err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(prods)
}

func main() {
	var err error
	// DB não existe mais por isso a conexão ta exposta
	urlBb := "postgres://postgres:postgres123@bd-pond.cib8iwbg0b6f.us-east-1.rds.amazonaws.com:5432/sample"
	db, _ = gorm.Open(postgres.Open(urlBb), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	createTableIfNotExists()

	http.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listEmployees(w, r)
		case http.MethodPost:
			addEmployee(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	handler := corsMiddleware(http.DefaultServeMux)
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listProducts(w, r)
		case http.MethodPost:
			addProduct(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	log.Println("Servidor rodando em :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
