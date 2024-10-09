/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"lil-todo/pkg/dbUtils"
	"lil-todo/pkg/format"

	_ "modernc.org/sqlite"

	"strings"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add <todo text>",
	Short:   "Add a new todo",
	Aliases: []string{"a"},
	Args:    cobra.MinimumNArgs(1),
	Run:     addCommand,
}

func addCommand(cmd *cobra.Command, args []string) {
	todoText := strings.Join(args, " ")

	db, err := dbUtils.ConnectToDb()
	if err != nil {
		cmd.PrintErrf("Error connecting to database: %v\n", err)
		return
	}
	defer db.Close()
	err = dbUtils.AddTodo(db, todoText)
	if err != nil {
		cmd.PrintErr("There was an error adding to todo list", err)
		return
	}
	todos, err := dbUtils.ListAllTodos(db)
	if err != nil {
		cmd.PrintErr("There was an error showing todo list", err)
	}
	format.PrintAllTodos(todos)
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
