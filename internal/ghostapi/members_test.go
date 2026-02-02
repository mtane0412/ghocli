/**
 * members_test.go
 * Test code for Members API
 */

package ghostapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListMembers_RetrieveMemberList retrieves a list of members
func TestListMembers_RetrieveMemberList(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/ghost/api/admin/members/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/members/")
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Verify that Authorization header exists
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorization header is not set")
		}

		// Return response
		response := map[string]interface{}{
			"members": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234567",
					"uuid":       "abc123-def456-ghi789",
					"email":      "john.smith@example.com",
					"name":       "John Smith",
					"note":       "Test member",
					"status":     "free",
					"created_at": "2024-01-15T10:00:00.000Z",
					"updated_at": "2024-01-15T10:00:00.000Z",
				},
				{
					"id":         "64fac5417c4c6b0001234568",
					"uuid":       "xyz987-uvw654-rst321",
					"email":      "alice.johnson@example.com",
					"name":       "Alice Johnson",
					"note":       "",
					"status":     "paid",
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
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Client creation error: %v", err)
	}

	// Retrieve member list
	resp, err := client.ListMembers(MemberListOptions{})
	if err != nil {
		t.Fatalf("Member list retrieval error: %v", err)
	}

	// Verify response
	if len(resp.Members) != 2 {
		t.Errorf("Member count = %d; want 2", len(resp.Members))
	}

	// Verify first member
	firstMember := resp.Members[0]
	if firstMember.Email != "john.smith@example.com" {
		t.Errorf("Email = %q; want %q", firstMember.Email, "john.smith@example.com")
	}
	if firstMember.Name != "John Smith" {
		t.Errorf("Name = %q; want %q", firstMember.Name, "John Smith")
	}
	if firstMember.Status != "free" {
		t.Errorf("Status = %q; want %q", firstMember.Status, "free")
	}
}

// TestListMembers_FilterParameter tests the filter parameter
func TestListMembers_FilterParameter(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		if !r.URL.Query().Has("filter") {
			t.Error("filter parameter is not set")
		}
		if r.URL.Query().Get("filter") != "status:paid" {
			t.Errorf("filter parameter = %q; want %q", r.URL.Query().Get("filter"), "status:paid")
		}

		// Return response
		response := map[string]interface{}{
			"members": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234568",
					"email":      "alice.johnson@example.com",
					"name":       "Alice Johnson",
					"status":     "paid",
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
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Client creation error: %v", err)
	}

	// Retrieve member list (filtered by status:paid)
	resp, err := client.ListMembers(MemberListOptions{
		Filter: "status:paid",
	})
	if err != nil {
		t.Fatalf("Member list retrieval error: %v", err)
	}

	// Verify response
	if len(resp.Members) != 1 {
		t.Errorf("Member count = %d; want 1", len(resp.Members))
	}
	if resp.Members[0].Status != "paid" {
		t.Errorf("Status = %q; want %q", resp.Members[0].Status, "paid")
	}
}

// TestGetMember_RetrieveMemberByID retrieves a member by ID
func TestGetMember_RetrieveMemberByID(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/ghost/api/admin/members/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "GET" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "GET")
		}

		// Return response
		response := map[string]interface{}{
			"members": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234567",
					"uuid":       "abc123-def456-ghi789",
					"email":      "john.smith@example.com",
					"name":       "John Smith",
					"note":       "Test member",
					"status":     "free",
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
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Client creation error: %v", err)
	}

	// Retrieve member
	member, err := client.GetMember("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("Member retrieval error: %v", err)
	}

	// Verify response
	if member.Email != "john.smith@example.com" {
		t.Errorf("Email = %q; want %q", member.Email, "john.smith@example.com")
	}
	if member.ID != "64fac5417c4c6b0001234567" {
		t.Errorf("ID = %q; want %q", member.ID, "64fac5417c4c6b0001234567")
	}
}

// TestCreateMember_CreateMember creates a member
func TestCreateMember_CreateMember(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/ghost/api/admin/members/" {
			t.Errorf("Request path = %q; want %q", r.URL.Path, "/ghost/api/admin/members/")
		}
		if r.Method != "POST" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "POST")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Request body parse error: %v", err)
		}

		members, ok := reqBody["members"].([]interface{})
		if !ok || len(members) == 0 {
			t.Error("members array does not exist in request body")
		}

		// Return response
		createdAt := time.Now().Format(time.RFC3339)
		response := map[string]interface{}{
			"members": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234999",
					"uuid":       "new-uuid-123",
					"email":      "new.member@example.com",
					"name":       "New Member",
					"note":       "Newly created member",
					"status":     "free",
					"created_at": createdAt,
					"updated_at": createdAt,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Client creation error: %v", err)
	}

	// Create member
	newMember := &Member{
		Email: "new.member@example.com",
		Name:  "New Member",
		Note:  "Newly created member",
	}

	createdMember, err := client.CreateMember(newMember)
	if err != nil {
		t.Fatalf("Member creation error: %v", err)
	}

	// Verify response
	if createdMember.Email != "new.member@example.com" {
		t.Errorf("Email = %q; want %q", createdMember.Email, "new.member@example.com")
	}
	if createdMember.ID == "" {
		t.Error("ID is empty")
	}
}

// TestUpdateMember_UpdateMember updates a member
func TestUpdateMember_UpdateMember(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/ghost/api/admin/members/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "PUT" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "PUT")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Request body parse error: %v", err)
		}

		// Return response
		response := map[string]interface{}{
			"members": []map[string]interface{}{
				{
					"id":         "64fac5417c4c6b0001234567",
					"uuid":       "abc123-def456-ghi789",
					"email":      "john.smith@example.com",
					"name":       "Updated Name",
					"note":       "Updated member",
					"status":     "free",
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
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Client creation error: %v", err)
	}

	// Update member
	updateMember := &Member{
		Name: "Updated Name",
		Note: "Updated member",
	}

	updatedMember, err := client.UpdateMember("64fac5417c4c6b0001234567", updateMember)
	if err != nil {
		t.Fatalf("Member update error: %v", err)
	}

	// Verify response
	if updatedMember.Name != "Updated Name" {
		t.Errorf("Name = %q; want %q", updatedMember.Name, "Updated Name")
	}
}

// TestDeleteMember_DeleteMember deletes a member
func TestDeleteMember_DeleteMember(t *testing.T) {
	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/ghost/api/admin/members/64fac5417c4c6b0001234567/"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q; want %q", r.URL.Path, expectedPath)
		}
		if r.Method != "DELETE" {
			t.Errorf("HTTP method = %q; want %q", r.Method, "DELETE")
		}

		// Return 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL, "test-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Client creation error: %v", err)
	}

	// Delete member
	err = client.DeleteMember("64fac5417c4c6b0001234567")
	if err != nil {
		t.Fatalf("Member deletion error: %v", err)
	}
}
