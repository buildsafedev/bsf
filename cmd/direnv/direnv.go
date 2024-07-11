package direnv

import (
	"errors"
	"fmt"
	"os"
	"regexp"
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
			err = setDirenv(envVar)
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
				return err
			}
			return nil
		}

	} else {
		err = os.WriteFile(".envrc", []byte("use flake bsf/."), 0644)
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
			if err != nil {
				return err
			}
			_, err = file.WriteString("\n.envrc")
			if err != nil {
				return err
			}
		}

		return nil
	} else {
		err = os.WriteFile(".gitignore", []byte(".envrc"), 0644)
		if err != nil {
			return err
		}
	}

	return nil

}

func setDirenv(args string) error {

	err := validateEnvVars(args)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(".envrc", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.WriteString("\nexport " + args)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("", err.Error()))
		return err
	}

	return nil
}

func validateEnvVars(args string) error {
	validKeyValRegex := regexp.MustCompile(`^[\w]+=[^\s]+$`)

	envVars := strings.Split(args, ",")

	for _, envVar := range envVars {
		if !validKeyValRegex.MatchString(envVar) {
			return errors.New("Invalid key-value pair format")
		}

		resp := strings.SplitN(envVar, "=", 2)
		key := resp[0]
		value := resp[1]

		if strings.ContainsAny(key, "= \t\n") {
			return errors.New("Invalid characters in key")
		}

		if strings.ContainsAny(value, "\x00") {
			return errors.New("Invalid characters in value")
		}
	}

	return nil
}
