package user

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go-rest-api/internal/apperror"
	"go-rest-api/internal/handlers"
	"go-rest-api/pkg/logging"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

const (
	userURL  = "/users"
	usersURL = "/users/:uuid"
)

type handler struct {
	logger  *logging.Logger
	service *Service
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetUserById))
	router.HandlerFunc(http.MethodPost, userURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodPut, usersURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodDelete, usersURL, apperror.Middleware(h.DeleteUser))
}

func NewHandler(logger *logging.Logger, service *Service) handlers.Handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	h.logger.Debug("GetList from handler")
	users, err := h.service.FindAll(r.Context())
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(users)
}

func (h *handler) GetUserById(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	idStr := params.ByName("uuid")
	h.logger.Debugf("GetUserById from handler by id: %s", idStr)

	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return apperror.NewBadRequestError
	}

	user, err := h.service.FindOne(r.Context(), oid.Hex())
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(user)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	h.logger.Debug("CreateUser from handler")
	var userDTO CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		return apperror.NewBadRequestError
	}
	user, err := h.service.Create(r.Context(), userDTO)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(user)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")
	h.logger.Debugf("UpdateUser from handler by id: %s", id)

	var userDTO UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		return apperror.NewBadRequestError
	}

	userDTO.ID = id

	if err := h.service.Update(r.Context(), userDTO); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")
	h.logger.Debugf("DeleteUser from handler by id: %s", id)

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return apperror.NewBadRequestError
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
