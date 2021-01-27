GITHUB_API_TOKEN := ""
VERSION :=""

run-tests:
	go test -cover -v

cover:
	go test -coverprofile=cover.tmp && go tool cover -html=cover.tmp

release:
	git checkout master
	git pull origin master
	git tag -a v$(VERSION) -m $(VERSION)
	git push origin master --follow-tags

changelog:
	git checkout master
	git pull origin master
	github_changelog_generator -t $(GITHUB_API_TOKEN)

push-changelog:
	git checkout master
	git pull origin master
	git add CHANGELOG.md
	git commit -m 'changelog updated.'
	git push origin master
