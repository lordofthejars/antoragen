package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "antoragen",
	Short: "Generator for Antora sites",
}

func main() {

	var projectName string
	var repo, public, startRepo string

	var cmdDoc = &cobra.Command{
		Use:   "doc",
		Short: "Generate Documentation",
		Long:  `generates documentation for a project/service`,
		Run: func(cmd *cobra.Command, args []string) {
			wd, _ := os.Getwd()
			err := generateDoc(wd, projectName)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmdDoc.Flags().StringVarP(&projectName, "project-name", "p", "", "Project Name")
	cmdDoc.MarkFlagRequired("project-name")

	var cmdSite = &cobra.Command{
		Use:   "site",
		Short: "Generate Site",
		Long:  `generates site project`,
		Run: func(cmd *cobra.Command, args []string) {
			wd, _ := os.Getwd()

			if len(startRepo) == 0 {
				lastSlashIndex := strings.LastIndex(repo, "/")
				lastPeriodIndex := strings.LastIndex(repo, ".")

				if lastPeriodIndex < 0 {
					lastPeriodIndex = len(repo)
				}

				startRepo = repo[lastSlashIndex+1 : lastPeriodIndex]

			}

			err := generateSite(wd, projectName, repo, public, startRepo)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmdSite.Flags().StringVarP(&projectName, "project-name", "p", "", "Project Name")
	cmdSite.Flags().StringVarP(&repo, "repo", "r", "", "Repository of Project to generate docs")
	cmdSite.Flags().StringVarP(&public, "public", "u", "", "Public URL where the site is going to be published")
	cmdSite.Flags().StringVarP(&startRepo, "start-repo", "s", "", "Repository name containing the initial page")

	cmdSite.MarkFlagRequired("project-name")
	cmdSite.MarkFlagRequired("repo")
	cmdSite.MarkFlagRequired("public")

	rootCmd.AddCommand(cmdDoc)
	rootCmd.AddCommand(cmdSite)
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf(err.Error())
		os.Exit(1)
	}

}
