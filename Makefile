GITHUB_API_TOKEN := ""

changelog:
	github_changelog_generator -t $(GITHUB_API_TOKEN)
	cat CHANGELOG.md | awk 'skip == 0 {print}; $$0 ~ "appix-1.0.1.1" {skip = 1}' | tee CHANGELOG.md | echo "Generated CHANGELOG.md"
