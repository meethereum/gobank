package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)


type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

type APIServer struct {
	listenAddr string
	store Storage
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountById))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	idstr := mux.Vars(r)["id"]
	id,err:= strconv.Atoi(idstr)
	if(err!=nil){
		return err
	}
	account,err:=s.store.GetAccountByID(id);
	if(err!=nil){
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

// create account from json recieved from user isliye r.Body is reader

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("hiii");
	createAccountRequest := new(CreateAccountRequest)
	// decode json to an object i.e, unmarshalling
	// from body of request(that is json) , decode the json into an object

	// it returns error so we are checking if it is equal to nil or not
	if err:= json.NewDecoder(r.Body).Decode(createAccountRequest); err!=nil{
		return err
	}

	// now remaining task is to make a query in db and also return all json to user
	// this feature can be used in json or coookies/ session (maybe idk)
	account:= NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName)
	if err:= s.store.CreateAccount(account); err!=nil{
		return err
	}
	return WriteJSON(w, http.StatusOK, account);
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}