package dockerfile

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/spf13/cobra"
)

var DGCmd = &cobra.Command{
	Use:     "dockerfile digests",
	Short:   "Replace Dockerfile image tags with immutable digests",
	Aliases: []string{"dg"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf dockerfile digests <Dockerfile>` to replace image tags with digests"))
			os.Exit(1)
		}

		dockerfile := args[1]
		file, err := os.Open(dockerfile)
		if err != nil {
			fmt.Printf("Error opening Dockerfile: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		re := regexp.MustCompile(`FROM\s+(\S+):(\S+)`)
		var updatedLines []string

		for scanner.Scan() {
			line := scanner.Text()
			if matches := re.FindStringSubmatch(line); matches != nil {
				image := matches[1]
				tag := matches[2]
				digest, err := getDigest(image, tag)
				if err != nil {
					fmt.Printf("Error retrieving digest for %s:%s: %v\n", image, tag, err)
					os.Exit(1)
				}
				line = strings.Replace(line, fmt.Sprintf("%s:%s", image, tag), fmt.Sprintf("%s@%s", image, digest), 1)
			}
			updatedLines = append(updatedLines, line)
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading Dockerfile: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(dockerfile, []byte(strings.Join(updatedLines, "\n")), 0644); err != nil {
			fmt.Printf("Error writing updated Dockerfile: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Dockerfile updated with image digests successfully.")
	},
}

func getDigest(image, tag string) (string, error) {
	digest, err := crane.Digest(fmt.Sprintf("%s:%s", image, tag))
	if err != nil {
		return "", fmt.Errorf("failed to get manifest for %s:%s: %w", image, tag, err)
	}

	return digest, nil
}
