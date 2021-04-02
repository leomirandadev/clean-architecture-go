package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/leomirandadev/clean-architecture-go/entities"
	"github.com/leomirandadev/clean-architecture-go/services"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

type controllers struct {
	srv   *services.Container
	log   logger.Logger
	token token.TokenHash
}

type UserController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
}

func New(srv *services.Container, log logger.Logger, tokenHasher token.TokenHash) UserController {
	return &controllers{srv: srv, log: log, token: tokenHasher}
}

func (ctr *controllers) Create(w http.ResponseWriter, r *http.Request) {
	var newUser entities.User
	json.NewDecoder(r.Body).Decode(&newUser)

	ctx := r.Context()
	err := ctr.srv.User.Create(ctx, newUser)

	if err != nil {
		ctr.log.Error("Ctrl.Create: ", "Error on create user: ", newUser)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (ctr *controllers) Auth(w http.ResponseWriter, r *http.Request) {
	var userLogin entities.UserAuth
	json.NewDecoder(r.Body).Decode(&userLogin)

	ctx := r.Context()
	userFound, err := ctr.srv.User.GetUserByLogin(ctx, userLogin)

	if err != nil {
		ctr.log.Error("Ctrl.Auth: ", "Error on find a user", userLogin)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := ctr.token.Encrypt(userFound)

	if err != nil {
		ctr.log.Error("Ctrl.Auth: ", "Error on generate token", userLogin)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entities.Response{Token: token})
}

func (ctr *controllers) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idUser, _ := strconv.ParseInt(params["id"], 10, 64)

	ctx := r.Context()
	user, err := ctr.srv.User.GetByID(ctx, idUser)

	if err != nil {
		ctr.log.Error("Ctrl.GetByid: ", "Error get user by id: ", idUser)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
