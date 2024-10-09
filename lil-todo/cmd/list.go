/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"lil-todo/pkg/dbUtils"
	"lil-todo/pkg/format"

	_ "modernc.org/sqlite"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var listCmd = &cobra.Command{
	Use:     "list all todos",
	Short:   "list all todos",
	Aliases: []string{"l"},
	Args:    cobra.ExactArgs(0),
	Run:     listCommand,
}

func listCommand(cmd *cobra.Command, args []string) {
	db, err := dbUtils.ConnectToDb()
	if err != nil {
		cmd.PrintErrf("Error connecting to database: %v\n", err)
		return
	}
	defer db.Close()
	todos, err := dbUtils.ListAllTodos(db)
	if err != nil {
		cmd.PrintErr("There was an error showing all todos", err)
		return
	}
	format.PrintAllTodos(todos)
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
