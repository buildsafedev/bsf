### Run goreleaser
`goreleaser release --snapshot --clean`

### Upload binaries to github

### Get Asset ID
`curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/buildsafedev/bsfrelease/releases/tags/v0.0.x``

### Update docs with asset id