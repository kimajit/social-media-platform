package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	handlers "socialmedia/handler"
)

func main() {
	// Initialize SocialMedia instance
	sm := handlers.NewSocialMedia()

	// Initialize router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/posts", sm.AddPostHandler).Methods("POST")
	router.HandleFunc("/posts/{id}", sm.GetPostHandler).Methods("GET")
	router.HandleFunc("/posts/{id}/comments", sm.AddCommentHandler).Methods("POST")
	router.HandleFunc("/posts/{id}/like", sm.LikePostHandler).Methods("POST")
	router.HandleFunc("/posts/{id}/dislike", sm.DislikePostHandler).Methods("POST")
	router.HandleFunc("/posts/{id}/share", sm.SharePostHandler).Methods("GET")

	// Start server
	log.Println("Server is running on port 8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
