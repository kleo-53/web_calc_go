package application

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kleo-53/web_calc_go/pkg/calculation"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() error {
	for {
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}

		text = strings.TrimSpace(text)

		if text == "exit" {
			log.Println("application was closed successfully")
			return nil
		}
		result, err := calculation.Calc(text)
		if err != nil {
			log.Println(text, "calculation was failed with error:", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error: " + err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := calculation.Calc(request.Expression)
	if err != nil {
		errorMessage := "error: " + err.Error()
		if errors.Is(err, calculation.ErrDivideByZero) || errors.Is(err, calculation.ErrInvalidExpression) || errors.Is(err, calculation.ErrUnknownSymbol) {
			http.Error(w, errorMessage, http.StatusUnprocessableEntity)
		} else {
			http.Error(w, errorMessage, http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "result: %f", result)
	}
}

func (a *Application) RunServer() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/calculate", CalcHandler).Methods("POST")
	return http.ListenAndServe(":"+a.config.Addr, router)
}
