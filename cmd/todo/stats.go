package main

import (
	"github.com/WHITE13452/toDoList/internal/cli"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "显示统计信息",
	Long:  "显示待办事项的统计信息，包括总数、完成数、完成率等。",
	Run: func(cmd *cobra.Command, args []string) {
		stats, err := store.GetStatistics()
		if err != nil {
			cli.PrintError("获取统计信息失败: %v", err)
			return
		}

		cli.PrintStatistics(stats)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
