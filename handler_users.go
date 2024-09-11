package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Cee-X/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig)handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}
	user , err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	respondWithJSON(w, 200, user)
}

func (apiCfg *apiConfig) handleGetUserByApiKey(w http.ResponseWriter, r *http.Request, user database.User){

	respondWithJSON(w, 200, databaseUserToUser(user))
}


func (apiCfg *apiConfig) handleGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User){

	post, err :=apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: 10,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to get posts: %v", err))
		return
	}
	respondWithJSON(w, 200, databasePostsToPosts(post))

}