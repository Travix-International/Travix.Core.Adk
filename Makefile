.PHONY: test vet lint testonly cover

GITHUB_API_TOKEN := ""

changelog:
	github_changelog_generator -t $(GITHUB_API_TOKEN)
	cat CHANGELOG.md | awk 'skip == 0 {print}; $$0 ~ "appix-1.0.1.1" {skip = 1}' | tee CHANGELOG.md >> /dev/null 2>&1 && echo "Generated CHANGELOG.md"


test: vet testonly

vet:
	go vet `go list ./... | grep -v /vendor/`

lint:
	golint `go list ./... | grep -v /vendor/`

testonly:
	go test -cover `go list ./... | grep -v /vendor/`

cover:
	echo Covering package ./$(PKG)
	go test -coverprofile=cover.out ./$(PKG) && go tool cover -html=cover.out
