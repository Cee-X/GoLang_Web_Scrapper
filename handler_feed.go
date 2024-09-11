package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/Cee-X/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig)handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		Name string `json:"name"`
		URL string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}
	feed , err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Url : params.URL,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to create feed: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedToFeed(feed))
}


func (apiCfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request){
	feeds , err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to get feeds: %v", err))
		return
	}
	respondWithJSON(w, 200, databaseFeedsToFeeds(feeds))
}