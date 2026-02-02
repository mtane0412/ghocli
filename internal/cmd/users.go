/**
 * users.go
 * User management commands
 *
 * Provides functionality for retrieving and updating Ghost users (site administrators and contributors).
 * User creation and deletion should be done using the Ghost dashboard's invitation feature.
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/ghocli/internal/fields"
	"github.com/mtane0412/ghocli/internal/ghostapi"
	"github.com/mtane0412/ghocli/internal/outfmt"
)

// UsersCmd is the user management command
type UsersCmd struct {
	List   UsersListCmd   `cmd:"" help:"List users"`
	Get    UsersInfoCmd   `cmd:"" help:"Show user information"`
	Update UsersUpdateCmd `cmd:"" help:"Update a user"`
}

// UsersListCmd is the command to retrieve user list
type UsersListCmd struct {
	Limit   int    `help:"Number of users to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page    int    `help:"Page number" short:"p" default:"1"`
	Include string `help:"Include additional data (e.g., roles,count.posts)" short:"i"`
	Filter  string `help:"Filter query" aliases:"where,w"`
}

// Run executes the list subcommand of the users command
func (c *UsersListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Parse field specification
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.UserFields)
		if err != nil {
			return fmt.Errorf("failed to parse field specification: %w", err)
		}
		selectedFields = parsedFields
	}

	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get user list
	response, err := client.ListUsers(ghostapi.UserListOptions{
		Limit:   c.Limit,
		Page:    c.Page,
		Include: c.Include,
		Filter:  c.Filter,
	})
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter and output if fields are specified
	if len(selectedFields) > 0 {
		// Convert User struct to map[string]interface{}
		var usersData []map[string]interface{}
		for _, user := range response.Users {
			userMap, err := outfmt.StructToMap(user)
			if err != nil {
				return fmt.Errorf("failed to convert user data: %w", err)
			}
			usersData = append(usersData, userMap)
		}

		// Filter fields and output
		return outfmt.FilterFields(formatter, usersData, selectedFields)
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Users)
	}

	// Output in table format
	headers := []string{"ID", "Name", "Slug", "Email", "Created"}
	rows := make([][]string, len(response.Users))
	for i, user := range response.Users {
		rows[i] = []string{
			user.ID,
			user.Name,
			user.Slug,
			user.Email,
			user.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// UsersInfoCmd is the command to show user information
type UsersInfoCmd struct {
	IDOrSlug string `arg:"" help:"User ID or slug (use 'slug:user-slug' format for slug)"`
}

// Run executes the info subcommand of the users command
func (c *UsersInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// Parse field specification
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.UserFields)
		if err != nil {
			return fmt.Errorf("failed to parse field specification: %w", err)
		}
		selectedFields = parsedFields
	}

	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get user
	user, err := client.GetUser(c.IDOrSlug)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter and output if fields are specified
	if len(selectedFields) > 0 {
		// Convert User struct to map[string]interface{}
		userMap, err := outfmt.StructToMap(user)
		if err != nil {
			return fmt.Errorf("failed to convert user data: %w", err)
		}

		// Filter fields and output
		return outfmt.FilterFields(formatter, []map[string]interface{}{userMap}, selectedFields)
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(user)
	}

	// Output in key/value format (no headers)
	rows := [][]string{
		{"id", user.ID},
		{"name", user.Name},
		{"slug", user.Slug},
		{"email", user.Email},
		{"bio", user.Bio},
		{"location", user.Location},
		{"website", user.Website},
		{"profile_image", user.ProfileImage},
		{"cover_image", user.CoverImage},
		{"created", user.CreatedAt.Format("2006-01-02 15:04:05")},
		{"updated", user.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	// Add role information
	if len(user.Roles) > 0 {
		roleNames := ""
		for i, role := range user.Roles {
			if i > 0 {
				roleNames += ", "
			}
			roleNames += role.Name
		}
		rows = append(rows, []string{"roles", roleNames})
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		return err
	}

	return formatter.Flush()
}

// UsersUpdateCmd is the command to update user
type UsersUpdateCmd struct {
	ID       string `arg:"" help:"User ID"`
	Name     string `help:"User name" short:"n"`
	Slug     string `help:"User slug"`
	Bio      string `help:"User bio" short:"b"`
	Location string `help:"User location" short:"l"`
	Website  string `help:"User website" short:"w"`
}

// Run executes the update subcommand of the users command
func (c *UsersUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing user
	existingUser, err := client.GetUser(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Apply updates
	updateUser := &ghostapi.User{
		Name:     existingUser.Name,
		Slug:     existingUser.Slug,
		Email:    existingUser.Email,
		Bio:      existingUser.Bio,
		Location: existingUser.Location,
		Website:  existingUser.Website,
	}

	if c.Name != "" {
		updateUser.Name = c.Name
	}
	if c.Slug != "" {
		updateUser.Slug = c.Slug
	}
	if c.Bio != "" {
		updateUser.Bio = c.Bio
	}
	if c.Location != "" {
		updateUser.Location = c.Location
	}
	if c.Website != "" {
		updateUser.Website = c.Website
	}

	// Update user
	updatedUser, err := client.UpdateUser(c.ID, updateUser)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("updated user: %s (ID: %s)", updatedUser.Name, updatedUser.ID))
	}

	// Also output user information if JSON format
	if root.JSON {
		return formatter.Print(updatedUser)
	}

	return nil
}
