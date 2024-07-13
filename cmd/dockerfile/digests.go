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
	Use:     "digests",
	Short:   "Replace Dockerfile image tags with immutable digests",
	Aliases: []string{"dg"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf dockerfile digests <Dockerfile>` to replace image tags with digests"))
			os.Exit(1)
		}

		dockerfile := args[0]
		file, err := os.ReadFile(dockerfile)
		if err != nil {
			fmt.Printf("Error opening Dockerfile: %v\n", err)
			os.Exit(1)
		}

		var dgMap = map[string]string{}

		dgMap, err = readByte(file)
		if err != nil {
			fmt.Printf("Error in readByte %v\n", err)
			os.Exit(1)
		}

		dgMap, err = getDigest(dgMap)
		if err != nil {
			fmt.Printf("Error retrieving digest %v\n", err)
			os.Exit(1)
		}

		updatedData, err := updateDockerfileWithDigests(file, dgMap)
		if err != nil {
			fmt.Printf("Error updating Dockerfile with digests: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(dockerfile, updatedData, 0644); err != nil {
			fmt.Printf("Error writing updated Dockerfile: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Changes made in")
		for _, data := range updatedData {
			fmt.Println(string(data))
		}
	},
}

func getDigest(dgMap map[string]string) (map[string]string, error) {

	for img := range dgMap {
		dg, err := crane.Digest(img)
		if err != nil {
			return nil, fmt.Errorf("failed to get manifest for %w", err)
		}
		dgMap[img] = dg
	}
	return dgMap, nil
}

func readByte(file []byte) (map[string]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(file)))
	re := regexp.MustCompile(`FROM\s+(\S+):(\S+)`)
	var lines = make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); matches != nil {
			image := matches[1]
			tag := matches[2]
			lines[fmt.Sprintf("%s:%s", image, tag)] = ""
		}
	}
	return lines, nil
}

func updateDockerfileWithDigests(data []byte, digestMap map[string]string) ([]byte, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var updatedLines []string

	for scanner.Scan() {
		line := scanner.Text()
		for tag, digest := range digestMap {
			if strings.Contains(line, tag) {
				line = strings.Replace(line, tag, digest, 1)
			}
		}
		updatedLines = append(updatedLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Dockerfile: %v", err)
	}

	return []byte(strings.Join(updatedLines, "\n")), nil
}
