package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Cee-X/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
);
type apiConfig struct {
	DB *database.Queries
}
func main() {
	godotenv.Load(".env");
	portString := os.Getenv("PORT");
	if portString == ""{
		log.Fatal("PORT is not found in the enviroment")
	}

    DbURl := os.Getenv("DB_URL");
	if DbURl == ""{
		log.Fatal("DbURL is not found in the enviroment")
	}
    conn, err := sql.Open("postgres", DbURl)
	if err != nil {
		log.Fatal("Can't connect to the database", err)
	}
	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}
	go startScrapping(db, 10, time.Minute)
	router := chi.NewRouter();

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,	
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter();
	v1Router.Get("/healthz", handlerReadiness); 
	v1Router.Get("/err", handlerErr);
	v1Router.Post("/users", apiCfg.handleCreateUser);
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUserByApiKey))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handleGetPostsForUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.handleGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollows))
	router.Mount("/v1", v1Router);

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Printf("Server is running on port %s", portString)
	err = srv.ListenAndServe();

	if err != nil {
		log.Fatal(err)
	}	
	
}

