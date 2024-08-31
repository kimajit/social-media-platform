package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"socialmedia/models"
	"socialmedia/utils"
)

// SocialMedia holds the in-memory store and mutex.
type SocialMedia struct {
	Posts map[string]*models.Post
	mu    sync.RWMutex
}

// NewSocialMedia initializes the SocialMedia struct.
func NewSocialMedia() *SocialMedia {
	return &SocialMedia{
		Posts: make(map[string]*models.Post),
	}
}

var (
	validate = validator.New()
	logger   = utils.InitializeLogger()
)

// AddPostHandler handles the creation of a new post.
func (sm *SocialMedia) AddPostHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		logger.Error("Invalid request payload", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if err := validate.Struct(post); err != nil {
		logger.Error("Validation error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed")
		return
	}

	// Generate unique ID
	post.ID = uuid.New().String()
	post.CreatedAt = time.Now()

	// Store the post
	sm.mu.Lock()
	sm.Posts[post.ID] = &post
	sm.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// GetPostHandler retrieves a post by ID.
func (sm *SocialMedia) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	sm.mu.RLock()
	post, exists := sm.Posts[postID]
	sm.mu.RUnlock()

	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// AddCommentHandler adds a comment to a specific post.
func (sm *SocialMedia) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	sm.mu.RLock()
	post, exists := sm.Posts[postID]
	sm.mu.RUnlock()

	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		logger.Error("Invalid request payload", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if err := validate.Struct(comment); err != nil {
		logger.Error("Validation error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed")
		return
	}

	// Generate unique ID
	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now()

	// Add comment to the post
	sm.mu.Lock()
	post.Comments = append(post.Comments, comment)
	sm.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// LikePostHandler increments the like count of a post.
func (sm *SocialMedia) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if post, exists := sm.Posts[postID]; exists {
		post.Likes++
		json.NewEncoder(w).Encode(post)
	} else {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
	}
}

// DislikePostHandler increments the dislike count of a post.
func (sm *SocialMedia) DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if post, exists := sm.Posts[postID]; exists {
		post.Dislikes++
		json.NewEncoder(w).Encode(post)
	} else {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
	}
}

// SharePostHandler generates a shareable link for a post.
func (sm *SocialMedia) SharePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	sm.mu.RLock()
	_, exists := sm.Posts[postID]
	sm.mu.RUnlock()

	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	// Generate shareable link
	baseURL := "http://localhost:8000/post/"
	shareableLink := baseURL + postID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"shareable_link": shareableLink})
}
