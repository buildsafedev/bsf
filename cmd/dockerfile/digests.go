package dockerfile

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

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
			fmt.Println(styles.ErrorStyle.Render("error:", "opening dockerfile:", err.Error()))
			os.Exit(1)
		}

		var dgMap = make(map[string]string)

		line, err := readByte(file)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error in parsing Dockerfile contents", err.Error()))
			os.Exit(1)
		}

		dgMap, err = getDigest(line)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("Error retrieving digest", err.Error()))
			os.Exit(1)
		}

		updatedData, err := updateDockerfileWithDigests(file, dgMap)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("Error updating Dockerfile with digests", err.Error()))
			os.Exit(1)
		}

		if err := os.WriteFile(dockerfile, updatedData, 0644); err != nil {
			fmt.Println(styles.ErrorStyle.Render("Error writing updated Dockerfile", err.Error()))
			os.Exit(1)
		}
	},
}

func getDigest(lines []string) (map[string]string, error) {
	var (
		dgMap = make(map[string]string)
		wg    sync.WaitGroup
		mu    sync.Mutex
	)

	for _, line := range lines {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			dg, err := crane.Digest(line)
			if err != nil {
				fmt.Println(styles.WarnStyle.Render("warning:", "skipping ", line, "can't find"))
				return
			}
			mu.Lock()
			dgMap[line] = dg
			mu.Unlock()
		}(line)
	}

	wg.Wait()
	return dgMap, nil
}

func readByte(file []byte) ([]string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(file))
	re := regexp.MustCompile(`FROM\s+(\S+):(\S+)`)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); matches != nil {
			image := matches[1]
			tag := matches[2]
			lines = append(lines, fmt.Sprintf("%s:%s", image, tag))
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
			img := strings.Split(tag, ":")
			line = strings.Replace(line, tag, img[0]+"@"+digest, 1)
		}
		updatedLines = append(updatedLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Dockerfile: %v", err)
	}

	return []byte(strings.Join(updatedLines, "\n")), nil
}
