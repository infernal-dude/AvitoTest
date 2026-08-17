package handler

import (
	"avito/internal/domain"
	"avito/internal/service"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type avitoHandler struct {
	service service.AvitoService
}

func NewService(service service.AvitoService) *avitoHandler {
	return &avitoHandler{service: service}
}

func (h *avitoHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Register(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *avitoHandler) Login(w http.ResponseWriter, r *http.Request) {
	var userData struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(userData.Password, userData.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *avitoHandler) DummyHandler(w http.ResponseWriter, r *http.Request) {
	userType := r.URL.Query().Get("user_type")
	if userType == "" {
		http.Error(w, "Client type is needed", http.StatusBadRequest)
		return
	}
	if userType != "moderator" && userType != "client" {
		http.Error(w, "Invalid client type", http.StatusBadRequest)
		return
	}

	token, err := h.service.LoginDummy(userType)
	if err != nil {
		http.Error(w, "error getting token", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *avitoHandler) CreateHouse(w http.ResponseWriter, r *http.Request) {
	userType, ok := r.Context().Value("user_type").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if userType != "moderator" {
		http.Error(w, "access denied", http.StatusUnauthorized)
		return
	}

	var house domain.House
	err := json.NewDecoder(r.Body).Decode(&house)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if house.Address == "" {
		http.Error(w, "address is required", http.StatusBadRequest)
		return
	}
	if house.Year <= 0 {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}

	houseResp, err := h.service.CreateHouse(house)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(houseResp)
}

func (h *avitoHandler) CreateFlat(w http.ResponseWriter, r *http.Request) {
	userType, ok := r.Context().Value("user_type").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if userType != "moderator" && userType != "client" {
		http.Error(w, "invalid user type", http.StatusBadRequest)
		return
	}

	var flat domain.Flat
	err := json.NewDecoder(r.Body).Decode(&flat)
	if err != nil {
		http.Error(w, "invaild data", http.StatusBadRequest)
	}

	if flat.Rooms <= 0 {
		http.Error(w, "rooms number cant be negative", http.StatusBadRequest)
		return
	}
	if flat.Price <= 0 {
		http.Error(w, "price cant be negative", http.StatusBadRequest)
		return
	}
	if flat.FlatNumber <= 0 {
		http.Error(w, "flats number cant be negative", http.StatusBadRequest)
		return
	}

	flatResp, err := h.service.CreateFlat(flat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flatResp)
}

func (h *avitoHandler) ModUpdate(w http.ResponseWriter, r *http.Request) {
	userType, ok := r.Context().Value("user_type").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if userType != "moderator" {
		http.Error(w, "access for moderator only", http.StatusUnauthorized)
		return
	}

	var modUpdate domain.ModUpdate
	err := json.NewDecoder(r.Body).Decode(&modUpdate)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if modUpdate.ID == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if modUpdate.Status == "" {
		http.Error(w, "status required", http.StatusBadRequest)
		return
	}

	flatResp, err := h.service.UpdateForMod(modUpdate.ID, modUpdate.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no flat with such id", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flatResp)
}

func (h *avitoHandler) GetFlatsByHouseId(w http.ResponseWriter, r *http.Request) {
	status, ok := r.Context().Value("user_type").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rawPath := r.URL.Path
	houseIdString := strings.TrimPrefix(rawPath, "/house/")
	houseId, err := strconv.Atoi(houseIdString)
	if err != nil {
		http.Error(w, "internal error getting house id", http.StatusInternalServerError)
		return
	}

	flatsArray, err := h.service.GetFlatsByHouseId(houseId, status)
	if err != nil {
		http.Error(w, "error getting flats", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flatsArray)
}
