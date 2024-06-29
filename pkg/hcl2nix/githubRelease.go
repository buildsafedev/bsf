package hcl2nix

import "fmt"

// ReadGitHubReleaseParams reads the github release params from the config
func ReadGitHubReleaseParams(conf *Config, app string) (*GitHubRelease, error) {
	if conf.GitHubReleases == nil {
		return nil, fmt.Errorf("no githubrelease block found")
	}

	for _, gh := range conf.GitHubReleases {
		if gh.App == app {
			if gh.Dir == "" {
				gh.Dir = "bsf-result/"
			}
			return &gh, nil
		}
	}

	return nil, fmt.Errorf("no githubrelease block found for app %s", app)
}
