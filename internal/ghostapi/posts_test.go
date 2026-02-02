/**
 * posts_test.go
 * Posts API test code
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListPosts_GetPostList tests retrieving a list of posts
func TestListPosts_GetPostList(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/posts/" {
			t.Errorf("request path = %q; want %q", r.URL.Path, "/ghost/api/admin/posts/")
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Verify Authorization header exists
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorization header not set")
		}

		// Return response
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":           "64fac5417c4c6b0001234567",
					"title":        "Test Post 1",
					"slug":         "test-post-1",
					"status":       "published",
					"created_at":   "2024-01-15T10:00:00.000Z",
					"updated_at":   "2024-01-15T10:00:00.000Z",
					"published_at": "2024-01-15T10:00:00.000Z",
				},
				{
					"id":         "64fac5417c4c6b0001234568",
					"title":      "Test Post 2",
					"slug":       "test-post-2",
					"status":     "draft",
					"created_at": "2024-01-16T10:00:00.000Z",
					"updated_at": "2024-01-16T10:00:00.000Z",
				},
			},
			"meta": map[string]interface{}{
				"pagination": map[string]interface{}{
					"page":  1,
					"limit": 15,
					"pages": 1,
					"total": 2,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Get post list
	response, err := client.ListPosts(ListOptions{})
	if err != nil {
		t.Fatalf("failed to get post list: %v", err)
	}

	// Verify response
	if len(response.Posts) != 2 {
		t.Errorf("number of posts = %d; want %d", len(response.Posts), 2)
	}
	if response.Posts[0].Title != "Test Post 1" {
		t.Errorf("post 1 title = %q; want %q", response.Posts[0].Title, "Test Post 1")
	}
	if response.Posts[0].Status != "published" {
		t.Errorf("post 1 status = %q; want %q", response.Posts[0].Status, "published")
	}
	if response.Posts[1].Title != "Test Post 2" {
		t.Errorf("post 2 title = %q; want %q", response.Posts[1].Title, "Test Post 2")
	}
	if response.Posts[1].Status != "draft" {
		t.Errorf("post 2 status = %q; want %q", response.Posts[1].Status, "draft")
	}
}

// TestListPosts_StatusFilter tests filtering posts by status
func TestListPosts_StatusFilter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		status := r.URL.Query().Get("filter")
		if status != "status:draft" {
			t.Errorf("status filter = %q; want %q", status, "status:draft")
		}

		// Return response (draft only)
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234568",
					"title":      "Draft Post",
					"slug":       "draft-post",
					"status":     "draft",
					"created_at": "2024-01-16T10:00:00.000Z",
					"updated_at": "2024-01-16T10:00:00.000Z",
				},
			},
			"meta": map[string]interface{}{
				"pagination": map[string]interface{}{
					"page":  1,
					"limit": 15,
					"pages": 1,
					"total": 1,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Filter by draft status
	response, err := client.ListPosts(ListOptions{Status: "draft"})
	if err != nil {
		t.Fatalf("failed to get post list: %v", err)
	}

	// Verify response
	if len(response.Posts) != 1 {
		t.Errorf("number of posts = %d; want %d", len(response.Posts), 1)
	}
	if response.Posts[0].Status != "draft" {
		t.Errorf("status = %q; want %q", response.Posts[0].Status, "draft")
	}
}

// TestGetPost_GetByID tests retrieving a post by ID
func TestGetPost_GetByID(t *testing.T) {
	postID := "64fac5417c4c6b0001234567"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/posts/" + postID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}

		// Return response
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":           postID,
					"title":        "Test Post",
					"slug":         "test-post",
					"html":         "<p>Body</p>",
					"status":       "published",
					"created_at":   "2024-01-15T10:00:00.000Z",
					"updated_at":   "2024-01-15T10:00:00.000Z",
					"published_at": "2024-01-15T10:00:00.000Z",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Get post
	post, err := client.GetPost(postID)
	if err != nil {
		t.Fatalf("failed to get post: %v", err)
	}

	// Verify response
	if post.ID != postID {
		t.Errorf("ID = %q; want %q", post.ID, postID)
	}
	if post.Title != "Test Post" {
		t.Errorf("Title = %q; want %q", post.Title, "Test Post")
	}
	if post.HTML != "<p>Body</p>" {
		t.Errorf("HTML = %q; want %q", post.HTML, "<p>Body</p>")
	}
}

// TestGetPost_GetBySlug tests retrieving a post by slug
func TestGetPost_GetBySlug(t *testing.T) {
	slug := "test-post"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/posts/slug/" + slug + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}

		// Return response
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":           "64fac5417c4c6b0001234567",
					"title":        "Test Post",
					"slug":         slug,
					"html":         "<p>Body</p>",
					"status":       "published",
					"created_at":   "2024-01-15T10:00:00.000Z",
					"updated_at":   "2024-01-15T10:00:00.000Z",
					"published_at": "2024-01-15T10:00:00.000Z",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Get post
	post, err := client.GetPost(slug)
	if err != nil {
		t.Fatalf("failed to get post: %v", err)
	}

	// Verify response
	if post.Slug != slug {
		t.Errorf("Slug = %q; want %q", post.Slug, slug)
	}
	if post.Title != "Test Post" {
		t.Errorf("Title = %q; want %q", post.Title, "Test Post")
	}
}

// TestCreatePost_CreatePost tests creating a new post
func TestCreatePost_CreatePost(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/posts/" {
			t.Errorf("request path = %q; want %q", r.URL.Path, "/ghost/api/admin/posts/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		posts := reqBody["posts"].([]interface{})
		post := posts[0].(map[string]interface{})
		if post["title"] != "New Post" {
			t.Errorf("Title = %q; want %q", post["title"], "New Post")
		}

		// Return response
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234569",
					"title":      "New Post",
					"slug":       "new-post",
					"status":     "draft",
					"created_at": time.Now().Format(time.RFC3339),
					"updated_at": time.Now().Format(time.RFC3339),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Create post
	newPost := &Post{
		Title:  "New Post",
		Status: "draft",
	}
	createdPost, err := client.CreatePost(newPost)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Verify response
	if createdPost.Title != "New Post" {
		t.Errorf("Title = %q; want %q", createdPost.Title, "New Post")
	}
	if createdPost.Status != "draft" {
		t.Errorf("Status = %q; want %q", createdPost.Status, "draft")
	}
	if createdPost.ID == "" {
		t.Error("ID is empty")
	}
}

// TestUpdatePost_UpdatePost tests updating an existing post
func TestUpdatePost_UpdatePost(t *testing.T) {
	postID := "64fac5417c4c6b0001234567"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/posts/" + postID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		posts := reqBody["posts"].([]interface{})
		post := posts[0].(map[string]interface{})
		if post["title"] != "Updated Title" {
			t.Errorf("Title = %q; want %q", post["title"], "Updated Title")
		}

		// Return response
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":         postID,
					"title":      "Updated Title",
					"slug":       "updated-post",
					"status":     "published",
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": time.Now().Format(time.RFC3339),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Update post
	updatePost := &Post{
		Title:  "Updated Title",
		Status: "published",
	}
	updatedPost, err := client.UpdatePost(postID, updatePost)
	if err != nil {
		t.Fatalf("failed to update post: %v", err)
	}

	// Verify response
	if updatedPost.ID != postID {
		t.Errorf("ID = %q; want %q", updatedPost.ID, postID)
	}
	if updatedPost.Title != "Updated Title" {
		t.Errorf("Title = %q; want %q", updatedPost.Title, "Updated Title")
	}
}

// TestUpdatePost_PreserveUpdatedAt tests updating a post while preserving updated_at timestamp
func TestUpdatePost_PreserveUpdatedAt(t *testing.T) {
	postID := "64fac5417c4c6b0001234567"
	originalUpdatedAt := "2024-01-15T10:00:00Z"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/posts/" + postID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		posts := reqBody["posts"].([]interface{})
		post := posts[0].(map[string]interface{})

		// Verify updated_at is included in request
		if _, ok := post["updated_at"]; !ok {
			t.Error("updated_at not included in request")
		}

		// Verify updated_at matches original value
		if post["updated_at"] != originalUpdatedAt {
			t.Errorf("updated_at = %q; want %q", post["updated_at"], originalUpdatedAt)
		}

		// Return response
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":         postID,
					"title":      "Updated Title",
					"slug":       "updated-post",
					"status":     "published",
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": time.Now().Format(time.RFC3339),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Update post with original updated_at
	parsedTime, _ := time.Parse(time.RFC3339, originalUpdatedAt)
	updatePost := &Post{
		Title:     "Updated Title",
		Status:    "published",
		UpdatedAt: parsedTime,
	}
	_, err = client.UpdatePost(postID, updatePost)
	if err != nil {
		t.Fatalf("failed to update post: %v", err)
	}
}

// TestDeletePost_DeletePost tests deleting a post
func TestDeletePost_DeletePost(t *testing.T) {
	postID := "64fac5417c4c6b0001234567"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/posts/" + postID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "DELETE" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "DELETE")
		}

		// Return response (204 No Content)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Delete post
	err = client.DeletePost(postID)
	if err != nil {
		t.Fatalf("failed to delete post: %v", err)
	}
}

// TestCreatePostWithOptions_WithHTMLSourceParameter tests creating a post with HTML source parameter
func TestCreatePostWithOptions_WithHTMLSourceParameter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/posts/" {
			t.Errorf("request path = %q; want %q", r.URL.Path, "/ghost/api/admin/posts/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Verify query parameters
		sourceParam := r.URL.Query().Get("source")
		if sourceParam != "html" {
			t.Errorf("source parameter = %q; want %q", sourceParam, "html")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		posts := reqBody["posts"].([]interface{})
		post := posts[0].(map[string]interface{})
		if post["title"] != "HTML Post" {
			t.Errorf("Title = %q; want %q", post["title"], "HTML Post")
		}
		if post["html"] != "<p>HTML Content</p>" {
			t.Errorf("HTML = %q; want %q", post["html"], "<p>HTML Content</p>")
		}

		// Return response (assumed to be converted to Lexical format by server)
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234569",
					"title":      "HTML Post",
					"slug":       "html-post",
					"html":       "<p>HTML Content</p>",
					"lexical":    `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"HTML Content","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
					"status":     "draft",
					"created_at": time.Now().Format(time.RFC3339),
					"updated_at": time.Now().Format(time.RFC3339),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Create HTML post with source=html option
	newPost := &Post{
		Title:  "HTML Post",
		HTML:   "<p>HTML Content</p>",
		Status: "draft",
	}
	opts := CreateOptions{
		Source: "html",
	}
	createdPost, err := client.CreatePostWithOptions(newPost, opts)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Verify response
	if createdPost.Title != "HTML Post" {
		t.Errorf("Title = %q; want %q", createdPost.Title, "HTML Post")
	}
	if createdPost.HTML != "<p>HTML Content</p>" {
		t.Errorf("HTML = %q; want %q", createdPost.HTML, "<p>HTML Content</p>")
	}
	// Verify Lexical format is set
	if createdPost.Lexical == "" {
		t.Error("Lexical is empty (should be converted by server)")
	}
}

// TestCreatePostWithOptions_BackwardsCompatibility tests creating a post without source parameter for backwards compatibility
func TestCreatePostWithOptions_BackwardsCompatibility(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/posts/" {
			t.Errorf("request path = %q; want %q", r.URL.Path, "/ghost/api/admin/posts/")
		}

		// Verify query parameter does not exist
		sourceParam := r.URL.Query().Get("source")
		if sourceParam != "" {
			t.Errorf("source parameter should not exist; got %q", sourceParam)
		}

		// Return response
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234569",
					"title":      "Normal Post",
					"slug":       "normal-post",
					"status":     "draft",
					"created_at": time.Now().Format(time.RFC3339),
					"updated_at": time.Now().Format(time.RFC3339),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Create normal post without source parameter
	newPost := &Post{
		Title:  "Normal Post",
		Status: "draft",
	}
	opts := CreateOptions{} // Source is empty string
	createdPost, err := client.CreatePostWithOptions(newPost, opts)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Verify response
	if createdPost.Title != "Normal Post" {
		t.Errorf("Title = %q; want %q", createdPost.Title, "Normal Post")
	}
}

// TestUpdatePostWithOptions_WithHTMLSourceParameter tests updating a post with HTML source parameter
func TestUpdatePostWithOptions_WithHTMLSourceParameter(t *testing.T) {
	postID := "64fac5417c4c6b0001234567"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/posts/" + postID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Verify query parameters
		sourceParam := r.URL.Query().Get("source")
		if sourceParam != "html" {
			t.Errorf("source parameter = %q; want %q", sourceParam, "html")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		posts := reqBody["posts"].([]interface{})
		post := posts[0].(map[string]interface{})
		if post["html"] != "<p>Updated HTML Content</p>" {
			t.Errorf("HTML = %q; want %q", post["html"], "<p>Updated HTML Content</p>")
		}

		// Return response
		response := map[string]interface{}{
			"posts": []map[string]interface{}{
				{
					"id":         postID,
					"title":      "Updated Post",
					"slug":       "updated-post",
					"html":       "<p>Updated HTML Content</p>",
					"lexical":    `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"Updated HTML Content","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
					"status":     "published",
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": time.Now().Format(time.RFC3339),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "keyid", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Update post with source=html option
	updatePost := &Post{
		Title:  "Updated Post",
		HTML:   "<p>Updated HTML Content</p>",
		Status: "published",
	}
	opts := CreateOptions{
		Source: "html",
	}
	updatedPost, err := client.UpdatePostWithOptions(postID, updatePost, opts)
	if err != nil {
		t.Fatalf("failed to update post: %v", err)
	}

	// Verify response
	if updatedPost.ID != postID {
		t.Errorf("ID = %q; want %q", updatedPost.ID, postID)
	}
	if updatedPost.HTML != "<p>Updated HTML Content</p>" {
		t.Errorf("HTML = %q; want %q", updatedPost.HTML, "<p>Updated HTML Content</p>")
	}
	// Verify Lexical format is set
	if updatedPost.Lexical == "" {
		t.Error("Lexical is empty (should be converted by server)")
	}
}
