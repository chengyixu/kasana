package main

import (
	"fmt"
	"os"

	"github.com/minervacap2022/klik-asana-cli/internal/api"
	"github.com/minervacap2022/klik-asana-cli/internal/auth"
	"github.com/minervacap2022/klik-asana-cli/internal/output"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kasana",
	Short: "Asana CLI for KLIK platform",
	Long:  "Command-line interface for Asana REST API. Auth via ASANA_PAT env var.",
}

func getClient() *api.Client {
	token, err := auth.GetToken()
	if err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
	return api.NewClient(token)
}

// --- task commands ---

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks in a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		project, _ := cmd.Flags().GetString("project")
		limit, _ := cmd.Flags().GetInt("limit")

		path := fmt.Sprintf("/tasks?project=%s&limit=%d&opt_fields=name,completed,assignee.name,due_on,notes", project, limit)
		result, err := client.Get(path)
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

var taskCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		project, _ := cmd.Flags().GetString("project")
		name, _ := cmd.Flags().GetString("name")
		notes, _ := cmd.Flags().GetString("notes")
		assignee, _ := cmd.Flags().GetString("assignee")

		payload := map[string]interface{}{
			"name":     name,
			"projects": []string{project},
		}
		if notes != "" {
			payload["notes"] = notes
		}
		if assignee != "" {
			payload["assignee"] = assignee
		}

		result, err := client.Post("/tasks", payload)
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

var taskUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		id, _ := cmd.Flags().GetString("id")
		name, _ := cmd.Flags().GetString("name")
		completed, _ := cmd.Flags().GetString("completed")

		payload := map[string]interface{}{}
		if name != "" {
			payload["name"] = name
		}
		if completed != "" {
			payload["completed"] = completed == "true"
		}

		result, err := client.Put(fmt.Sprintf("/tasks/%s", id), payload)
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

var taskViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		id, _ := cmd.Flags().GetString("id")

		result, err := client.Get(fmt.Sprintf("/tasks/%s?opt_fields=name,completed,assignee.name,due_on,notes,memberships.project.name,tags.name", id))
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

// --- project commands ---

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects in a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		workspace, _ := cmd.Flags().GetString("workspace")
		limit, _ := cmd.Flags().GetInt("limit")

		path := fmt.Sprintf("/projects?workspace=%s&limit=%d&opt_fields=name,owner.name,current_status,due_date", workspace, limit)
		result, err := client.Get(path)
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

// --- section commands ---

var sectionCmd = &cobra.Command{
	Use:   "section",
	Short: "Manage sections",
}

var sectionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sections in a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		project, _ := cmd.Flags().GetString("project")

		result, err := client.Get(fmt.Sprintf("/projects/%s/sections", project))
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

// --- comment commands ---

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage comments",
}

var commentAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a comment to a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		task, _ := cmd.Flags().GetString("task")
		text, _ := cmd.Flags().GetString("text")

		payload := map[string]interface{}{
			"text": text,
		}

		result, err := client.Post(fmt.Sprintf("/tasks/%s/stories", task), payload)
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

// --- user commands ---

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users in a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		workspace, _ := cmd.Flags().GetString("workspace")

		result, err := client.Get(fmt.Sprintf("/workspaces/%s/users?opt_fields=name,email", workspace))
		if err != nil {
			return err
		}
		output.RawJSON(result)
		return nil
	},
}

func init() {
	// task
	taskListCmd.Flags().String("project", "", "Project GID")
	taskListCmd.Flags().Int("limit", 25, "Max tasks")
	taskListCmd.MarkFlagRequired("project")

	taskCreateCmd.Flags().String("project", "", "Project GID")
	taskCreateCmd.Flags().String("name", "", "Task name")
	taskCreateCmd.Flags().String("notes", "", "Task notes/description")
	taskCreateCmd.Flags().String("assignee", "", "Assignee GID")
	taskCreateCmd.MarkFlagRequired("project")
	taskCreateCmd.MarkFlagRequired("name")

	taskUpdateCmd.Flags().String("id", "", "Task GID")
	taskUpdateCmd.Flags().String("name", "", "New name")
	taskUpdateCmd.Flags().String("completed", "", "true/false")
	taskUpdateCmd.MarkFlagRequired("id")

	taskViewCmd.Flags().String("id", "", "Task GID")
	taskViewCmd.MarkFlagRequired("id")

	taskCmd.AddCommand(taskListCmd, taskCreateCmd, taskUpdateCmd, taskViewCmd)

	// project
	projectListCmd.Flags().String("workspace", "", "Workspace GID")
	projectListCmd.Flags().Int("limit", 25, "Max projects")
	projectListCmd.MarkFlagRequired("workspace")
	projectCmd.AddCommand(projectListCmd)

	// section
	sectionListCmd.Flags().String("project", "", "Project GID")
	sectionListCmd.MarkFlagRequired("project")
	sectionCmd.AddCommand(sectionListCmd)

	// comment
	commentAddCmd.Flags().String("task", "", "Task GID")
	commentAddCmd.Flags().String("text", "", "Comment text")
	commentAddCmd.MarkFlagRequired("task")
	commentAddCmd.MarkFlagRequired("text")
	commentCmd.AddCommand(commentAddCmd)

	// user
	userListCmd.Flags().String("workspace", "", "Workspace GID")
	userListCmd.MarkFlagRequired("workspace")
	userCmd.AddCommand(userListCmd)

	rootCmd.AddCommand(taskCmd, projectCmd, sectionCmd, commentCmd, userCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
