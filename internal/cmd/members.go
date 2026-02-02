/**
 * members.go
 * メンバー管理コマンド
 *
 * Ghostメンバー（購読者）の作成、更新、削除機能を提供します。
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

// MembersCmd はメンバー管理コマンドです
type MembersCmd struct {
	List   MembersListCmd   `cmd:"" help:"List members"`
	Get    MembersInfoCmd   `cmd:"" help:"Show member information"`
	Create MembersCreateCmd `cmd:"" help:"Create a member"`
	Update MembersUpdateCmd `cmd:"" help:"Update a member"`
	Delete MembersDeleteCmd `cmd:"" help:"Delete a member"`

	// Phase 1: ステータス別一覧ショートカット
	Paid MembersPaidCmd `cmd:"" help:"List paid members"`
	Free MembersFreeCmd `cmd:"" help:"List free members"`

	// Phase 3: ラベル操作
	Label   MembersLabelCmd   `cmd:"" help:"Add label to member"`
	Unlabel MembersUnlabelCmd `cmd:"" help:"Remove label from member"`
	Recent  MembersRecentCmd  `cmd:"" help:"List recently created members"`
}

// MembersListCmd is the command to retrieve メンバー list
type MembersListCmd struct {
	Limit  int    `help:"Number of members to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page   int    `help:"Page number" short:"p" default:"1"`
	Filter string `help:"Filter query (e.g., status:paid)" aliases:"where,w"`
	Order  string `help:"Sort order (e.g., created_at DESC)" short:"o"`
}

// Run executes the list subcommand of the members command
func (c *MembersListCmd) Run(ctx context.Context, root *RootFlags) error {
	// Parse field specification
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.MemberFields)
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

	// Get member list
	response, err := client.ListMembers(ghostapi.MemberListOptions{
		Limit:  c.Limit,
		Page:   c.Page,
		Filter: c.Filter,
		Order:  c.Order,
	})
	if err != nil {
		return fmt.Errorf("failed to list members: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter and output if fields are specified
	if len(selectedFields) > 0 {
		// Convert Member struct to map[string]interface{}
		var membersData []map[string]interface{}
		for _, member := range response.Members {
			memberMap, err := outfmt.StructToMap(member)
			if err != nil {
				return fmt.Errorf("failed to convert member data: %w", err)
			}
			membersData = append(membersData, memberMap)
		}

		// Filter fields and output
		return outfmt.FilterFields(formatter, membersData, selectedFields)
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Members)
	}

	// Output in table format
	headers := []string{"ID", "Email", "Name", "Status", "Created"}
	rows := make([][]string, len(response.Members))
	for i, member := range response.Members {
		rows[i] = []string{
			member.ID,
			member.Email,
			member.Name,
			member.Status,
			member.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// MembersInfoCmd is the command to show メンバー information
type MembersInfoCmd struct {
	ID string `arg:"" help:"Member ID"`
}

// Run executes the info subcommand of the members command
func (c *MembersInfoCmd) Run(ctx context.Context, root *RootFlags) error {
	// Parse field specification
	var selectedFields []string
	if root.Fields != "" {
		parsedFields, err := fields.Parse(root.Fields, fields.MemberFields)
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

	// Get member
	member, err := client.GetMember(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Filter and output if fields are specified
	if len(selectedFields) > 0 {
		// Convert Member struct to map[string]interface{}
		memberMap, err := outfmt.StructToMap(member)
		if err != nil {
			return fmt.Errorf("failed to convert member data: %w", err)
		}

		// Filter fields and output
		return outfmt.FilterFields(formatter, []map[string]interface{}{memberMap}, selectedFields)
	}

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(member)
	}

	// Output in key/value format (no headers)
	rows := [][]string{
		{"id", member.ID},
		{"uuid", member.UUID},
		{"email", member.Email},
		{"name", member.Name},
		{"note", member.Note},
		{"status", member.Status},
		{"created", member.CreatedAt.Format("2006-01-02 15:04:05")},
		{"updated", member.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		return err
	}

	return formatter.Flush()
}

// MembersCreateCmd is the command to create メンバー
type MembersCreateCmd struct {
	Email  string   `help:"Member email (required)" short:"e" required:""`
	Name   string   `help:"Member name" short:"n"`
	Note   string   `help:"Member note" short:"t"`
	Labels []string `help:"Member labels" short:"l"`
}

// Run executes the create subcommand of the members command
func (c *MembersCreateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Create new member
	newMember := &ghostapi.Member{
		Email: c.Email,
		Name:  c.Name,
		Note:  c.Note,
	}

	// Add labels
	if len(c.Labels) > 0 {
		labels := make([]ghostapi.Label, len(c.Labels))
		for i, labelName := range c.Labels {
			labels[i] = ghostapi.Label{Name: labelName}
		}
		newMember.Labels = labels
	}

	createdMember, err := client.CreateMember(newMember)
	if err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("created member: %s (ID: %s)", createdMember.Email, createdMember.ID))
	}

	// Also output member information if JSON format
	if root.JSON {
		return formatter.Print(createdMember)
	}

	return nil
}

// MembersUpdateCmd is the command to update メンバー
type MembersUpdateCmd struct {
	ID     string   `arg:"" help:"Member ID"`
	Name   string   `help:"Member name" short:"n"`
	Note   string   `help:"Member note" short:"t"`
	Labels []string `help:"Member labels" short:"l"`
}

// Run executes the update subcommand of the members command
func (c *MembersUpdateCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing member
	existingMember, err := client.GetMember(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	// Apply updates
	updateMember := &ghostapi.Member{
		Email:  existingMember.Email,
		Name:   existingMember.Name,
		Note:   existingMember.Note,
		Labels: existingMember.Labels,
	}

	if c.Name != "" {
		updateMember.Name = c.Name
	}
	if c.Note != "" {
		updateMember.Note = c.Note
	}
	if len(c.Labels) > 0 {
		labels := make([]ghostapi.Label, len(c.Labels))
		for i, labelName := range c.Labels {
			labels[i] = ghostapi.Label{Name: labelName}
		}
		updateMember.Labels = labels
	}

	// Update member
	updatedMember, err := client.UpdateMember(c.ID, updateMember)
	if err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("updated member: %s (ID: %s)", updatedMember.Email, updatedMember.ID))
	}

	// Also output member information if JSON format
	if root.JSON {
		return formatter.Print(updatedMember)
	}

	return nil
}

// MembersDeleteCmd is the command to delete メンバー
type MembersDeleteCmd struct {
	ID string `arg:"" help:"Member ID"`
}

// Run executes the delete subcommand of the members command
func (c *MembersDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get member information to build confirmation message
	member, err := client.GetMember(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	// Confirm destructive operation
	action := fmt.Sprintf("delete member '%s' (ID: %s)", member.Email, c.ID)
	if err := ConfirmDestructive(ctx, root, action); err != nil {
		return err
	}

	// Delete member
	if err := client.DeleteMember(c.ID); err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	formatter.PrintMessage(fmt.Sprintf("deleted member (ID: %s)", c.ID))

	return nil
}

// ========================================
// Phase 1: ステータス別一覧ショートカット
// ========================================

// MembersPaidCmd is the command to retrieve paid member list
type MembersPaidCmd struct {
	Limit int `help:"Number of members to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run executes the paid subcommand of the members command
func (c *MembersPaidCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get paid member list
	response, err := client.ListMembers(ghostapi.MemberListOptions{
		Limit:  c.Limit,
		Page:   c.Page,
		Filter: "status:paid",
	})
	if err != nil {
		return fmt.Errorf("failed to list paid members: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Members)
	}

	// Output in table format
	headers := []string{"ID", "Email", "Name", "Status", "Created"}
	rows := make([][]string, len(response.Members))
	for i, member := range response.Members {
		rows[i] = []string{
			member.ID,
			member.Email,
			member.Name,
			member.Status,
			member.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// MembersFreeCmd is the command to retrieve free member list
type MembersFreeCmd struct {
	Limit int `help:"Number of members to retrieve" short:"l" aliases:"max,n" default:"15"`
	Page  int `help:"Page number" short:"p" default:"1"`
}

// Run executes the free subcommand of the members command
func (c *MembersFreeCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get free member list
	response, err := client.ListMembers(ghostapi.MemberListOptions{
		Limit:  c.Limit,
		Page:   c.Page,
		Filter: "status:free",
	})
	if err != nil {
		return fmt.Errorf("failed to list free members: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Members)
	}

	// Output in table format
	headers := []string{"ID", "Email", "Name", "Status", "Created"}
	rows := make([][]string, len(response.Members))
	for i, member := range response.Members {
		rows[i] = []string{
			member.ID,
			member.Email,
			member.Name,
			member.Status,
			member.CreatedAt.Format("2006-01-02"),
		}
	}

	return formatter.PrintTable(headers, rows)
}

// ========================================
// Phase 3: ラベル操作
// ========================================

// MembersLabelCmd is the command to add label to メンバー
type MembersLabelCmd struct {
	ID    string `arg:"" help:"Member ID"`
	Label string `arg:"" help:"Label name"`
}

// Run executes the label subcommand of the members command
func (c *MembersLabelCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing member
	existingMember, err := client.GetMember(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	// Skip if label name already exists in existing labels
	for _, label := range existingMember.Labels {
		if label.Name == c.Label {
			// Create output formatter
			formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())
			formatter.PrintMessage(fmt.Sprintf("member already has label '%s' holds label (ID: %s)", c.Label, c.ID))
			return nil
		}
	}

	// Add new label to existing labels
	newLabels := append(existingMember.Labels, ghostapi.Label{Name: c.Label})

	// Update member
	updateMember := &ghostapi.Member{
		Email:  existingMember.Email,
		Name:   existingMember.Name,
		Note:   existingMember.Note,
		Labels: newLabels,
	}

	updatedMember, err := client.UpdateMember(c.ID, updateMember)
	if err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("added label to member: %s (ID: %s, Label: %s)", updatedMember.Email, updatedMember.ID, c.Label))
	}

	// Also output member information if JSON format
	if root.JSON {
		return formatter.Print(updatedMember)
	}

	return nil
}

// MembersUnlabelCmd is the command to remove label from メンバー
type MembersUnlabelCmd struct {
	ID    string `arg:"" help:"Member ID"`
	Label string `arg:"" help:"Label name"`
}

// Run executes the unlabel subcommand of the members command
func (c *MembersUnlabelCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get existing member
	existingMember, err := client.GetMember(c.ID)
	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	// Remove specified label from existing labels
	var newLabels []ghostapi.Label
	found := false
	for _, label := range existingMember.Labels {
		if label.Name != c.Label {
			newLabels = append(newLabels, label)
		} else {
			found = true
		}
	}

	// If label not found
	if !found {
		// Create output formatter
		formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())
		formatter.PrintMessage(fmt.Sprintf("member does not have label '%s' does not have label (ID: %s)", c.Label, c.ID))
		return nil
	}

	// Update member
	updateMember := &ghostapi.Member{
		Email:  existingMember.Email,
		Name:   existingMember.Name,
		Note:   existingMember.Note,
		Labels: newLabels,
	}

	updatedMember, err := client.UpdateMember(c.ID, updateMember)
	if err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Show success message
	if !root.JSON {
		formatter.PrintMessage(fmt.Sprintf("removed label from member: %s (ID: %s, Label: %s)", updatedMember.Email, updatedMember.ID, c.Label))
	}

	// Also output member information if JSON format
	if root.JSON {
		return formatter.Print(updatedMember)
	}

	return nil
}

// MembersRecentCmd is the command to retrieve recently registered メンバー list
type MembersRecentCmd struct {
	Limit int `help:"Number of members to retrieve" short:"l" aliases:"max,n" default:"15"`
}

// Run executes the recent subcommand of the members command
func (c *MembersRecentCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get recently registered member list (sorted by created_at DESC)
	response, err := client.ListMembers(ghostapi.MemberListOptions{
		Limit: c.Limit,
		Page:  1,
		Order: "created_at DESC",
	})
	if err != nil {
		return fmt.Errorf("failed to list members: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(response.Members)
	}

	// Output in table format
	headers := []string{"ID", "Email", "Name", "Status", "Created"}
	rows := make([][]string, len(response.Members))
	for i, member := range response.Members {
		rows[i] = []string{
			member.ID,
			member.Email,
			member.Name,
			member.Status,
			member.CreatedAt.Format("2006-01-02 15:04"),
		}
	}

	return formatter.PrintTable(headers, rows)
}
