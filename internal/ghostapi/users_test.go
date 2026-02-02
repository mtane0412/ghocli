/**
 * users_test.go
 * Tests for Users API
 *
 * Provides tests for Ghost Admin API Users functionality.
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListUsers_GetUserList tests basic user list retrieval
func TestListUsers_GetUserList(t *testing.T) {
	// Create test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected method: GET, actual: %s", r.Method)
		}
		if r.URL.Path != "/ghost/api/admin/users/" {
			t.Errorf("expected path: /ghost/api/admin/users/, actual: %s", r.URL.Path)
		}

		// Return test response
		resp := UserListResponse{
			Users: []User{
				{
					ID:    "user1",
					Name:  "John Smith",
					Slug:  "john-smith",
					Email: "john@example.com",
				},
				{
					ID:    "user2",
					Name:  "Alice Johnson",
					Slug:  "alice-johnson",
					Email: "alice@example.com",
				},
			},
		}
		resp.Meta.Pagination.Page = 1
		resp.Meta.Pagination.Limit = 15
		resp.Meta.Pagination.Pages = 1
		resp.Meta.Pagination.Total = 2

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Create test client
	client, err := NewClient(ts.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Get user list
	resp, err2 := client.ListUsers(UserListOptions{})
	if err2 != nil {
		t.Fatalf("failed to retrieve user list: %v", err2)
	}

	// Verify response
	if len(resp.Users) != 2 {
		t.Errorf("expected number of users: 2, actual: %d", len(resp.Users))
	}

	// Verify first user
	if resp.Users[0].Name != "John Smith" {
		t.Errorf("expected name: John Smith, actual: %s", resp.Users[0].Name)
	}
	if resp.Users[0].Email != "john@example.com" {
		t.Errorf("expected email address: john@example.com, actual: %s", resp.Users[0].Email)
	}
}

// TestListUsers_WithOptions tests user list retrieval with query parameters
func TestListUsers_WithOptions(t *testing.T) {
	// Create test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()
		if query.Get("limit") != "5" {
			t.Errorf("expected limit: 5, actual: %s", query.Get("limit"))
		}
		if query.Get("page") != "2" {
			t.Errorf("expected page: 2, actual: %s", query.Get("page"))
		}
		if query.Get("include") != "roles,count.posts" {
			t.Errorf("expected include: roles,count.posts, actual: %s", query.Get("include"))
		}

		// Return test response
		resp := UserListResponse{
			Users: []User{
				{
					ID:    "user3",
					Name:  "Bob Williams",
					Slug:  "bob-williams",
					Email: "bob@example.com",
					Roles: []Role{
						{ID: "role1", Name: "Author"},
					},
				},
			},
		}
		resp.Meta.Pagination.Page = 2
		resp.Meta.Pagination.Limit = 5
		resp.Meta.Pagination.Pages = 3
		resp.Meta.Pagination.Total = 15

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Create test client
	client, err := NewClient(ts.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Get user list with options
	opts := UserListOptions{
		Limit:   5,
		Page:    2,
		Include: "roles,count.posts",
	}
	resp, err2 := client.ListUsers(opts)
	if err2 != nil {
		t.Fatalf("failed to retrieve user list: %v", err2)
	}

	// Verify response
	if len(resp.Users) != 1 {
		t.Errorf("expected number of users: 1, actual: %d", len(resp.Users))
	}
	if resp.Meta.Pagination.Page != 2 {
		t.Errorf("expected page: 2, actual: %d", resp.Meta.Pagination.Page)
	}
	if len(resp.Users[0].Roles) != 1 {
		t.Errorf("expected number of roles: 1, actual: %d", len(resp.Users[0].Roles))
	}
}

// TestGetUser_GetByID tests retrieving a user by ID
func TestGetUser_GetByID(t *testing.T) {
	// Create test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected method: GET, actual: %s", r.Method)
		}
		if r.URL.Path != "/ghost/api/admin/users/user123/" {
			t.Errorf("expected path: /ghost/api/admin/users/user123/, actual: %s", r.URL.Path)
		}

		// Return test response
		resp := UserResponse{
			Users: []User{
				{
					ID:           "user123",
					Name:         "Charlie Brown",
					Slug:         "charlie-brown",
					Email:        "charlie@example.com",
					Bio:          "Engineer",
					Location:     "Tokyo",
					Website:      "https://example.com",
					ProfileImage: "https://example.com/profile.jpg",
					CreatedAt:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Create test client
	client, err := NewClient(ts.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Get user
	user, err2 := client.GetUser("user123")
	if err2 != nil {
		t.Fatalf("failed to retrieve user: %v", err2)
	}

	// Verify response
	if user.ID != "user123" {
		t.Errorf("expected ID: user123, actual: %s", user.ID)
	}
	if user.Name != "Charlie Brown" {
		t.Errorf("expected name: Charlie Brown, actual: %s", user.Name)
	}
	if user.Email != "charlie@example.com" {
		t.Errorf("expected email address: charlie@example.com, actual: %s", user.Email)
	}
	if user.Bio != "Engineer" {
		t.Errorf("expected bio: Engineer, actual: %s", user.Bio)
	}
}

// TestGetUser_GetBySlug tests retrieving a user by slug
func TestGetUser_GetBySlug(t *testing.T) {
	// Create test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected method: GET, actual: %s", r.Method)
		}
		if r.URL.Path != "/ghost/api/admin/users/slug/charlie-brown/" {
			t.Errorf("expected path: /ghost/api/admin/users/slug/charlie-brown/, actual: %s", r.URL.Path)
		}

		// Return test response
		resp := UserResponse{
			Users: []User{
				{
					ID:    "user123",
					Name:  "Charlie Brown",
					Slug:  "charlie-brown",
					Email: "charlie@example.com",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Create test client
	client, err := NewClient(ts.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Get user by slug
	user, err2 := client.GetUser("slug:charlie-brown")
	if err2 != nil {
		t.Fatalf("failed to retrieve user: %v", err2)
	}

	// Verify response
	if user.Slug != "charlie-brown" {
		t.Errorf("expected slug: charlie-brown, actual: %s", user.Slug)
	}
	if user.Name != "Charlie Brown" {
		t.Errorf("expected name: Charlie Brown, actual: %s", user.Name)
	}
}

// TestUpdateUser_UpdateUser tests updating a user
func TestUpdateUser_UpdateUser(t *testing.T) {
	// Create test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "PUT" {
			t.Errorf("expected method: PUT, actual: %s", r.Method)
		}
		if r.URL.Path != "/ghost/api/admin/users/user123/" {
			t.Errorf("expected path: /ghost/api/admin/users/user123/, actual: %s", r.URL.Path)
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}

		users, ok := reqBody["users"].([]interface{})
		if !ok || len(users) == 0 {
			t.Fatal("request body does not contain users")
		}

		user := users[0].(map[string]interface{})
		if user["name"] != "Updated Name" {
			t.Errorf("expected name: Updated Name, actual: %v", user["name"])
		}

		// Return test response
		resp := UserResponse{
			Users: []User{
				{
					ID:       "user123",
					Name:     "Updated Name",
					Slug:     "updated-slug",
					Email:    "updated@example.com",
					Bio:      "Updated Bio",
					Location: "Osaka",
					Website:  "https://updated.example.com",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Create test client
	client, err := NewClient(ts.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Update user
	updateData := &User{
		Name:     "Updated Name",
		Slug:     "updated-slug",
		Bio:      "Updated Bio",
		Location: "Osaka",
		Website:  "https://updated.example.com",
	}
	user, err2 := client.UpdateUser("user123", updateData)
	if err2 != nil {
		t.Fatalf("failed to update user: %v", err2)
	}

	// Verify response
	if user.Name != "Updated Name" {
		t.Errorf("expected name: Updated Name, actual: %s", user.Name)
	}
	if user.Bio != "Updated Bio" {
		t.Errorf("expected bio: Updated Bio, actual: %s", user.Bio)
	}
	if user.Location != "Osaka" {
		t.Errorf("expected location: Osaka, actual: %s", user.Location)
	}
}

// TestGetUser_UserNotFound tests error handling when user is not found
func TestGetUser_UserNotFound(t *testing.T) {
	// Create test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return empty response
		resp := UserResponse{
			Users: []User{},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Create test client
	client, err := NewClient(ts.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Get user (expect error to be returned)
	_, err2 := client.GetUser("nonexistent")
	if err2 == nil {
		t.Fatal("Expected error but no error was returned")
	}
}

// TestListUsers_APIError tests error handling when API returns an error
func TestListUsers_APIError(t *testing.T) {
	// Create test HTTP server (return error)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":[{"message":"Internal Server Error"}]}`))
	}))
	defer ts.Close()

	// Create test client
	client, err := NewClient(ts.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Get user list (expect error to be returned)
	_, err2 := client.ListUsers(UserListOptions{})
	if err2 == nil {
		t.Fatal("Expected error but no error was returned")
	}
}
