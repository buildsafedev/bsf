package dockerfile

import (
	"context"
	"fmt"
	"os"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/spf13/cobra"
	"github.com/stacklok/frizbee/pkg/replacer"
	"github.com/stacklok/frizbee/pkg/utils/config"
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
		file, err := os.Open(dockerfile)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", "opening dockerfile:", err.Error()))
			os.Exit(1)
		}
		defer file.Close()

		r := replacer.NewContainerImagesReplacer(config.DefaultConfig())

		ok, str, err := r.ParseFile(context.TODO(), file)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error in parsing Dockerfile contents", err.Error()))
			os.Exit(1)
		}
		if err := os.WriteFile(dockerfile, []byte(str), 0644); err != nil {
			fmt.Println(styles.ErrorStyle.Render("Error writing updated Dockerfile", err.Error()))
			os.Exit(1)
		}
		if ok {
			fmt.Println(styles.SucessStyle.Render("Dockerfile changed !!!"))
		} else {
			fmt.Println(styles.ErrorStyle.Render("Dockerfile unchanged !!!"))
		}
	},
}

// func getDigest(lines []string) (map[string]string, error) {
// 	var (
// 		dgMap = make(map[string]string)
// 		wg    sync.WaitGroup
// 		mu    sync.Mutex
// 	)

// 	for _, line := range lines {
// 		wg.Add(1)
// 		go func(line string) {
// 			defer wg.Done()
// 			dg, err := crane.Digest(line)
// 			if err != nil {
// 				fmt.Println(styles.WarnStyle.Render("warning:", "skipping ", line, "can't find"))
// 				return
// 			}
// 			mu.Lock()
// 			dgMap[line] = dg
// 			mu.Unlock()
// 		}(line)
// 	}

// 	wg.Wait()
// 	return dgMap, nil
// }

// func updateDockerfileWithDigests(data *os.File, digestMap map[string]string) ([]byte, error) {
// 	if _, err := data.Seek(0, 0); err != nil {
// 		return nil, fmt.Errorf("error seeking to the beginning of the file: %v", err)
// 	}
// 	scanner := bufio.NewScanner(data)
// 	var updatedLines []string

// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		for tag, digest := range digestMap {
// 			img := strings.Split(tag, ":")
// 			line = strings.Replace(line, tag, img[0]+"@"+digest, 1)
// 		}
// 		updatedLines = append(updatedLines, line)
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return nil, fmt.Errorf("error reading Dockerfile: %v", err)
// 	}

// 	return []byte(strings.Join(updatedLines, "\n")), nil
// }
