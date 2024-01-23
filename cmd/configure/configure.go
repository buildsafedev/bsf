package configure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/config"
	"github.com/spf13/cobra"
)

// ConfigureCmd represents the configure command
var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configures global settings for bsf",
	Long: `
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// todo : let user configure settings
		_, err := PreCheckConf()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

	},
}

// PreCheckConf checks if the ~/.bsf.json file exists and if not creates it.
func PreCheckConf() (*config.Config, error) {
	// check home directory for ~/.bsf.json file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("Unable to get home directory:", err.Error()))
		return nil, err
	}

	bsfFilePath := filepath.Join(homeDir, ".bsf.json")
	if _, err := os.Stat(bsfFilePath); err != nil {
		file, err := os.Create(bsfFilePath)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("Unable to create file:", err.Error()))
			return nil, err
		}
		defer file.Close()

		conf := config.Config{
			BuildSafeAPI:    "api.buildsafe.dev:443",
			BuildSafeAPITLS: true,
		}

		jsonBytes, err := json.MarshalIndent(conf, "", "  ")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("Unable to marshal json:", err.Error()))
			return nil, err
		}
		_, err = file.Write(jsonBytes)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("Unable to write to file:", err.Error()))
			return nil, err
		}

		return &conf, nil
	}

	fb, err := os.ReadFile(bsfFilePath)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("Unable to read file:", err.Error()))
		return nil, err
	}

	conf := &config.Config{}
	err = json.Unmarshal(fb, conf)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("Unable to unmarshal json:", err.Error()))
		return nil, err
	}
	return conf, nil
}
