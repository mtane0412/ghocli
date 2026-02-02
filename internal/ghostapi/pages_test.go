/**
 * pages_test.go
 * Pages API test code
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListPages_GetPageList tests retrieving a list of pages
func TestListPages_GetPageList(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/pages/" {
			t.Errorf("request path = %q; want %q", r.URL.Path, "/ghost/api/admin/pages/")
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
			"pages": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234601",
					"title":      "Test Page 1",
					"slug":       "test-page-1",
					"status":     "published",
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": "2024-01-15T10:00:00.000Z",
				},
				{
					"id":         "64fac5417c4c6b0001234602",
					"title":      "Test Page 2",
					"slug":       "test-page-2",
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

	// Get page list
	response, err := client.ListPages(ListOptions{})
	if err != nil {
		t.Fatalf("failed to get page list: %v", err)
	}

	// Verify response
	if len(response.Pages) != 2 {
		t.Errorf("number of pages = %d; want %d", len(response.Pages), 2)
	}
	if response.Pages[0].Title != "Test Page 1" {
		t.Errorf("page 1 title = %q; want %q", response.Pages[0].Title, "Test Page 1")
	}
	if response.Pages[0].Status != "published" {
		t.Errorf("page 1 status = %q; want %q", response.Pages[0].Status, "published")
	}
	if response.Pages[1].Title != "Test Page 2" {
		t.Errorf("page 2 title = %q; want %q", response.Pages[1].Title, "Test Page 2")
	}
	if response.Pages[1].Status != "draft" {
		t.Errorf("page 2 status = %q; want %q", response.Pages[1].Status, "draft")
	}
}

// TestGetPage_GetByID tests retrieving a page by ID
func TestGetPage_GetByID(t *testing.T) {
	pageID := "64fac5417c4c6b0001234601"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/pages/" + pageID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}

		// Return response
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":         pageID,
					"title":      "Test Page",
					"slug":       "test-page",
					"html":       "<p>Page Body</p>",
					"status":     "published",
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": "2024-01-15T10:00:00.000Z",
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

	// Get page
	page, err := client.GetPage(pageID)
	if err != nil {
		t.Fatalf("failed to get page: %v", err)
	}

	// Verify response
	if page.ID != pageID {
		t.Errorf("ID = %q; want %q", page.ID, pageID)
	}
	if page.Title != "Test Page" {
		t.Errorf("Title = %q; want %q", page.Title, "Test Page")
	}
	if page.HTML != "<p>Page Body</p>" {
		t.Errorf("HTML = %q; want %q", page.HTML, "<p>Page Body</p>")
	}
}

// TestCreatePage_CreatePage tests creating a new page
func TestCreatePage_CreatePage(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/pages/" {
			t.Errorf("request path = %q; want %q", r.URL.Path, "/ghost/api/admin/pages/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		pages := reqBody["pages"].([]interface{})
		page := pages[0].(map[string]interface{})
		if page["title"] != "New Page" {
			t.Errorf("Title = %q; want %q", page["title"], "New Page")
		}

		// Return response
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234603",
					"title":      "New Page",
					"slug":       "new-page",
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

	// Create page
	newPage := &Page{
		Title:  "New Page",
		Status: "draft",
	}
	createdPage, err := client.CreatePage(newPage)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Verify response
	if createdPage.Title != "New Page" {
		t.Errorf("Title = %q; want %q", createdPage.Title, "New Page")
	}
	if createdPage.Status != "draft" {
		t.Errorf("Status = %q; want %q", createdPage.Status, "draft")
	}
	if createdPage.ID == "" {
		t.Error("ID is empty")
	}
}

// TestUpdatePage_UpdatePage tests updating an existing page
func TestUpdatePage_UpdatePage(t *testing.T) {
	pageID := "64fac5417c4c6b0001234601"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/pages/" + pageID + "/"
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
		pages := reqBody["pages"].([]interface{})
		page := pages[0].(map[string]interface{})
		if page["title"] != "Updated Page Title" {
			t.Errorf("Title = %q; want %q", page["title"], "Updated Page Title")
		}

		// Return response
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":         pageID,
					"title":      "Updated Page Title",
					"slug":       "updated-page",
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

	// Update page
	updatePage := &Page{
		Title:  "Updated Page Title",
		Status: "published",
	}
	updatedPage, err := client.UpdatePage(pageID, updatePage)
	if err != nil {
		t.Fatalf("failed to update page: %v", err)
	}

	// Verify response
	if updatedPage.ID != pageID {
		t.Errorf("ID = %q; want %q", updatedPage.ID, pageID)
	}
	if updatedPage.Title != "Updated Page Title" {
		t.Errorf("Title = %q; want %q", updatedPage.Title, "Updated Page Title")
	}
}

// TestDeletePage_DeletePage tests deleting a page
func TestDeletePage_DeletePage(t *testing.T) {
	pageID := "64fac5417c4c6b0001234601"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/pages/" + pageID + "/"
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

	// Delete page
	err = client.DeletePage(pageID)
	if err != nil {
		t.Fatalf("failed to delete page: %v", err)
	}
}

// TestGetPage_ExtendedFieldParsing tests parsing extended fields of a page
func TestGetPage_ExtendedFieldParsing(t *testing.T) {
	pageID := "64fac5417c4c6b0001234601"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return response with extended fields
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":         pageID,
					"title":      "Extended Fields Test Page",
					"slug":       "extended-fields-test",
					"html":       "<p>Page Body</p>",
					"status":     "published",
					"url":        "https://example.com/extended-fields-test/",
					"excerpt":    "This is a page excerpt.",
					"visibility": "public",
					"featured":   true,
					"authors": []map[string]interface{}{
						{
							"id":   "author1",
							"name": "John Doe",
						},
					},
					"tags": []map[string]interface{}{
						{
							"id":   "tag1",
							"name": "Test",
						},
						{
							"id":   "tag2",
							"name": "Sample",
						},
					},
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

	// Get page
	page, err := client.GetPage(pageID)
	if err != nil {
		t.Fatalf("failed to get page: %v", err)
	}

	// Verify basic fields
	if page.ID != pageID {
		t.Errorf("ID = %q; want %q", page.ID, pageID)
	}
	if page.Title != "Extended Fields Test Page" {
		t.Errorf("Title = %q; want %q", page.Title, "Extended Fields Test Page")
	}

	// Verify extended fields
	if page.URL != "https://example.com/extended-fields-test/" {
		t.Errorf("URL = %q; want %q", page.URL, "https://example.com/extended-fields-test/")
	}
	if page.Excerpt != "This is a page excerpt." {
		t.Errorf("Excerpt = %q; want %q", page.Excerpt, "This is a page excerpt.")
	}
	if page.Visibility != "public" {
		t.Errorf("Visibility = %q; want %q", page.Visibility, "public")
	}
	if !page.Featured {
		t.Errorf("Featured = %v; want %v", page.Featured, true)
	}

	// Verify Authors
	if len(page.Authors) != 1 {
		t.Errorf("number of authors = %d; want %d", len(page.Authors), 1)
	}
	if len(page.Authors) > 0 && page.Authors[0].Name != "John Doe" {
		t.Errorf("Authors[0].Name = %q; want %q", page.Authors[0].Name, "John Doe")
	}

	// Verify Tags
	if len(page.Tags) != 2 {
		t.Errorf("number of tags = %d; want %d", len(page.Tags), 2)
	}
	if len(page.Tags) > 0 && page.Tags[0].Name != "Test" {
		t.Errorf("Tags[0].Name = %q; want %q", page.Tags[0].Name, "Test")
	}
	if len(page.Tags) > 1 && page.Tags[1].Name != "Sample" {
		t.Errorf("Tags[1].Name = %q; want %q", page.Tags[1].Name, "Sample")
	}
}

// TestCreatePageWithOptions_WithHTMLSourceParameter tests creating a page with HTML source parameter
func TestCreatePageWithOptions_WithHTMLSourceParameter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/ghost/api/admin/pages/" {
			t.Errorf("request path = %q; want %q", r.URL.Path, "/ghost/api/admin/pages/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Verify query parameters
		if r.URL.Query().Get("source") != "html" {
			t.Errorf("source parameter = %q; want %q", r.URL.Query().Get("source"), "html")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		pages := reqBody["pages"].([]interface{})
		page := pages[0].(map[string]interface{})
		if page["title"] != "HTML Page" {
			t.Errorf("Title = %q; want %q", page["title"], "HTML Page")
		}
		if page["html"] != "<h1>Heading</h1><p>Paragraph</p>" {
			t.Errorf("HTML = %q; want %q", page["html"], "<h1>Heading</h1><p>Paragraph</p>")
		}

		// Return response (converted to Lexical format)
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":      "64fac5417c4c6b0001234604",
					"title":   "HTML Page",
					"slug":    "html-page",
					"status":  "draft",
					"lexical": `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"Heading","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"heading","version":1,"tag":"h1"},{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"Paragraph","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
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

	// Create HTML page
	newPage := &Page{
		Title:  "HTML Page",
		HTML:   "<h1>Heading</h1><p>Paragraph</p>",
		Status: "draft",
	}
	opts := CreateOptions{Source: "html"}
	createdPage, err := client.CreatePageWithOptions(newPage, opts)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Verify response
	if createdPage.Title != "HTML Page" {
		t.Errorf("Title = %q; want %q", createdPage.Title, "HTML Page")
	}
	if createdPage.Status != "draft" {
		t.Errorf("Status = %q; want %q", createdPage.Status, "draft")
	}
	// Verify Lexical format is set
	if createdPage.Lexical == "" {
		t.Error("Lexical is empty")
	}
}

// TestUpdatePageWithOptions_WithHTMLSourceParameter tests updating a page with HTML source parameter
func TestUpdatePageWithOptions_WithHTMLSourceParameter(t *testing.T) {
	pageID := "64fac5417c4c6b0001234601"

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		expectedPath := "/ghost/api/admin/pages/" + pageID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Verify query parameters
		if r.URL.Query().Get("source") != "html" {
			t.Errorf("source parameter = %q; want %q", r.URL.Query().Get("source"), "html")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		pages := reqBody["pages"].([]interface{})
		page := pages[0].(map[string]interface{})
		if page["html"] != "<h1>Updated Heading</h1>" {
			t.Errorf("HTML = %q; want %q", page["html"], "<h1>Updated Heading</h1>")
		}

		// Return response
		response := map[string]interface{}{
			"pages": []map[string]interface{}{
				{
					"id":      pageID,
					"title":   "Updated Page",
					"slug":    "updated-html-page",
					"status":  "published",
					"lexical": `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"Updated Heading","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"heading","version":1,"tag":"h1"}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
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

	// Update page
	updatePage := &Page{
		Title:  "Updated Page",
		HTML:   "<h1>Updated Heading</h1>",
		Status: "published",
	}
	opts := CreateOptions{Source: "html"}
	updatedPage, err := client.UpdatePageWithOptions(pageID, updatePage, opts)
	if err != nil {
		t.Fatalf("failed to update page: %v", err)
	}

	// Verify response
	if updatedPage.ID != pageID {
		t.Errorf("ID = %q; want %q", updatedPage.ID, pageID)
	}
	if updatedPage.Title != "Updated Page" {
		t.Errorf("Title = %q; want %q", updatedPage.Title, "Updated Page")
	}
	// Verify Lexical format is set
	if updatedPage.Lexical == "" {
		t.Error("Lexical is empty")
	}
}
