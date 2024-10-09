/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"lil-todo/pkg/dbUtils"
	"lil-todo/pkg/format"
	"strconv"

	_ "modernc.org/sqlite"

	"github.com/spf13/cobra"
)

// toggleCmd represents the add command
var toggleCmd = &cobra.Command{
	Use:     "toggle check <todo id>",
	Short:   "toggle checkmark todo",
	Aliases: []string{"t"},
	Args:    cobra.ExactArgs(1),
	Run:     toggleCommand,
}

func toggleCommand(cmd *cobra.Command, args []string) {
	todoIDStr := args[0]
	// Convert string to integer
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		cmd.PrintErrf("Invalid todo ID: %s\n", todoIDStr)
		return
	}

	db, err := dbUtils.ConnectToDb()
	if err != nil {
		cmd.PrintErrf("Error connecting to database: %v\n", err)
		return
	}
	defer db.Close()
	_, err = dbUtils.ToggleTodo(db, todoID)
	if err != nil {
		cmd.PrintErr("There was an error saving toggle to todo list: ", err)
		return
	}
	todos, err := dbUtils.ListAllTodos(db)
	if err != nil {
		cmd.PrintErr("There was an error showing todo list", err)
	}
	format.PrintAllTodos(todos)
}

func init() {
	rootCmd.AddCommand(toggleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toggleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// toggleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
