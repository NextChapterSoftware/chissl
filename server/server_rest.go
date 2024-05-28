package chserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/NextChapterSoftware/chissl/share/settings"
	"net/http"
	"regexp"
	"strings"
)

type UpdateUserRequest struct {
	Name    string           `json:"username"`
	Pass    string           `json:"password,omitempty"`
	Addrs   []*regexp.Regexp `json:"addresses,omitempty"`
	IsAdmin bool             `json:"is_admin"`
}

// decodeBasicAuthHeader extracts the username and password from auth headers
func (s *Server) decodeBasicAuthHeader(headers http.Header) (username, password string, ok bool) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", "", false
	}
	const basicAuthPrefix = "Basic "
	if !strings.HasPrefix(authHeader, basicAuthPrefix) {
		return "", "", false
	}

	// Decode the base64 encoded username:password
	decoded, err := base64.StdEncoding.DecodeString(authHeader[len(basicAuthPrefix):])
	if err != nil {
		return "", "", false
	}

	// Split the username and password
	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return "", "", false
	}

	return credentials[0], credentials[1], true
}

// BasicAuthMiddleware validates the username and password
func (s *Server) basicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimSpace(s.config.AuthFile) == "" {
			http.Error(w, "No auth file configured on server", http.StatusUnauthorized)
			return
		}

		username, password, ok := s.decodeBasicAuthHeader(r.Header)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		u, found := s.users.Get(username)

		// Validate the credentials (Replace with your validation logic)
		if !found || username != u.Name || password != u.Pass || !u.IsAdmin {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	}
}

func getUsernameFromPath(path string) (string, error) {
	// Define the expected URL pattern
	pattern := `^/user/([^/]+)$`
	re := regexp.MustCompile(pattern)

	// Check if the path matches the pattern
	matches := re.FindStringSubmatch(path)
	if matches == nil {
		return "", fmt.Errorf("invalid URL format: %s", path)
	}
	return matches[1], nil
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	data, err := s.users.ToJSON()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(data))
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	u, found := s.users.Get(username)
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responseJson, err := u.ToJSON()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(responseJson))
}

func (s *Server) handleAddUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to create a User object
	var newUser settings.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the user input
	if err := newUser.ValidateUser(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, found := s.users.Get(newUser.Name)
	if found {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Add the user to the server's user collection
	s.users.Set(newUser.Name, &newUser)
	err := s.users.WriteUsers()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with a status indicating success
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to create a User object
	var targetUser settings.User
	if err := json.NewDecoder(r.Body).Decode(&targetUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	targetUserFromLookup, found := s.users.Get(targetUser.Name)
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get current user making this request
	requestingUser, _, _ := s.decodeBasicAuthHeader(r.Header)

	// Admins cannot revoke admin permission from themselves
	if !targetUser.IsAdmin && targetUser.Name == requestingUser {
		http.Error(w, "Cannot revoke admin from yourself", http.StatusBadRequest)
		return
	}

	if targetUser.Pass == "" {
		targetUser.Pass = targetUserFromLookup.Pass
	}

	if len(targetUser.Addrs) == 0 {
		targetUser.Addrs = targetUserFromLookup.Addrs
	}

	s.users.Set(targetUser.Name, &targetUser)
	err := s.users.WriteUsers()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Implement user update logic here
	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	u, found := s.users.Get(username)
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get current user making this request
	requestingUser, _, _ := s.decodeBasicAuthHeader(r.Header)
	if requestingUser == u.Name {
		http.Error(w, "Cannot delete your own user", http.StatusBadRequest)
		return
	}

	s.users.Del(u.Name)
	err = s.users.WriteUsers()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleAuthfile(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to create a User object
	var users []*settings.User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(users) == 0 {
		http.Error(w, "No users found in file", http.StatusBadRequest)
	}

	// Get current user making this request
	requestingUser, _, _ := s.decodeBasicAuthHeader(r.Header)
	u, _ := s.users.Get(requestingUser)

	requestingUserFromPayload := &settings.User{}
	for _, user := range users {
		err := user.ValidateUser()
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid user setting for %s: %v", user.Name, err), http.StatusBadRequest)
		}
		if user.Name == u.Name {
			requestingUserFromPayload = user
		}
	}
	if requestingUserFromPayload == nil || !requestingUserFromPayload.IsAdmin {
		http.Error(w, "file must include the current requesting user with admin permission", http.StatusBadRequest)
		return
	}

	s.users.Reset(users)
	err := s.users.WriteUsers()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Implement user update logic here
	w.WriteHeader(http.StatusAccepted)
}
