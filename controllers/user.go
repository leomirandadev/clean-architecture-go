package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/leomirandadev/clean-architecture-go/entities"
	"github.com/leomirandadev/clean-architecture-go/services"
	"github.com/leomirandadev/clean-architecture-go/utils/hasher"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

type controllers struct {
	userService services.UserService
	log         logger.Logger
	token       token.TokenHash
}

type UserController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
}

func NewUserController(userService services.UserService, log logger.Logger, tokenHasher token.TokenHash) UserController {
	return &controllers{userService: userService, log: log, token: tokenHasher}
}

func (ctr *controllers) Create(w http.ResponseWriter, r *http.Request) {
	var newUser entities.User
	json.NewDecoder(r.Body).Decode(&newUser)

	hasherBcrypt := hasher.NewBcryptHasher()
	passwordHashed, errHash := hasherBcrypt.Generate(newUser.Password)

	if errHash != nil {
		ctr.log.Error("Ctrl.Create: ", "Error generate hash password: ", newUser)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser.Password = passwordHashed
	err := ctr.userService.New(newUser)

	if err != nil {
		ctr.log.Error("Ctrl.Create: ", "Error on create user: ", newUser)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (ctr *controllers) Auth(w http.ResponseWriter, r *http.Request) {
	var newUser entities.User
	json.NewDecoder(r.Body).Decode(&newUser)

	hasherBcrypt := hasher.NewBcryptHasher()
	passwordHashed, errHash := hasherBcrypt.Generate(newUser.Password)

	if errHash != nil {
		ctr.log.Error("Ctrl.Auth: ", "Error generate hash password: ", newUser)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser.Password = passwordHashed
	err := ctr.userService.New(newUser)

	if err != nil {
		ctr.log.Error("Ctrl.Auth: ", "Error on create user: ", newUser)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (ctr *controllers) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idUser, _ := strconv.ParseInt(params["id"], 10, 64)
	user, err := ctr.userService.GetByID(idUser)

	if err != nil {
		ctr.log.Error("Ctrl.GetByid: ", "Error get user by id: ", idUser)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
