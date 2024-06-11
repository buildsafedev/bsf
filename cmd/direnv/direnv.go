package direnv

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/spf13/cobra"
)

var Direnv = &cobra.Command{
	Use:   "direnv",
	Short: "direnv initializes the direnv environment for the project",
	Run: func(cmd *cobra.Command, args []string) {
		err := generateEnvrc()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}
		err = fetchGitignore()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		if envVar != "" {
			err = setDIrenv(envVar)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
				os.Exit(1)
			}
		}
	},
}

var (
	envVar string
)

func init() {
	Direnv.Flags().StringVarP(&envVar, "env", "e", "", "set environment variable [key=value]")
}

func generateEnvrc() error {

	if _, err := os.Stat(".envrc"); err == nil {

		read, err := os.ReadFile(".envrc")
		if err != nil {
			return err
		}

		if !strings.Contains(string(read), "use flake bsf/.") {

			file, err := os.OpenFile(".envrc", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}

			_, err = file.WriteString("\nuse flake bsf/.")
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("", err.Error()))
				return err
			}
			return nil
		}
		fmt.Println(styles.HelpStyle.Render(" ✅ .envrc already exists"))

	} else {
		file, err := os.Create(".envrc")
		if err != nil {
			return err
		}

		defer file.Close()

		_, err = file.WriteString("use flake bsf/.")
		if err != nil {
			return err
		}

		fmt.Println(styles.HelpStyle.Render(" ✅ .envrc generated"))

		cmd := exec.Command("direnv", "allow")
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchGitignore() error {

	if _, err := os.Stat(".gitignore"); err == nil {

		read, err := os.ReadFile(".gitignore")
		if err != nil {
			return err
		}

		if !strings.Contains(string(read), ".envrc") {

			file, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0644)
			_, err = file.WriteString("\n.envrc")
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("", err.Error()))
				return err
			}
		}
		fmt.Println(styles.HelpStyle.Render(" ✅ .gitignore already exists"))

		return nil
	} else {
		_, err := os.Create(".gitignore")
		if err != nil {
			return err
		}

		file, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		_, err = file.WriteString("\n.envrc")
		if err != nil {
			return err
		}

		fmt.Println(styles.HelpStyle.Render(" ✅ .gitignore generated"))
	}

	return nil

}

func setDIrenv(args string) error {

	if !strings.Contains(args, "=") {
		return fmt.Errorf("Hint: use --env key=value")
	} else {
		file, err := os.OpenFile(".envrc", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		_, err = file.WriteString("\nexport " + args)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("", err.Error()))
			return err
		}

		cmd := exec.Command("direnv", "allow")
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
