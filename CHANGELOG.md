# Change Log

## [Unreleased](https://github.com/Travix-International/appix/tree/HEAD)

[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.16...HEAD)

**Merged pull requests:**

- Implement pushing the app package directly to Google Cloud Storage [\#81](https://github.com/Travix-International/appix/pull/81) ([markvincze](https://github.com/markvincze))
- Added certContent and keyContent keys go config helper and to build script [\#80](https://github.com/Travix-International/appix/pull/80) ([nunorfpt](https://github.com/nunorfpt))
- Update readme file [\#79](https://github.com/Travix-International/appix/pull/79) ([nunorfpt](https://github.com/nunorfpt))
- Add a helper script to generate the Firebase config env vars [\#78](https://github.com/Travix-International/appix/pull/78) ([markvincze](https://github.com/markvincze))

## [1.1.16](https://github.com/Travix-International/appix/tree/1.1.16) (2017-06-16)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.15...1.1.16)

**Merged pull requests:**

- Update to the README.md file [\#77](https://github.com/Travix-International/appix/pull/77) ([nunorfpt](https://github.com/nunorfpt))
- Fail safe if app catalog does not return valid poll URI [\#76](https://github.com/Travix-International/appix/pull/76) ([alexmiranda](https://github.com/alexmiranda))
- Introduce the ability to override the dev server URL when pushing an app [\#75](https://github.com/Travix-International/appix/pull/75) ([markvincze](https://github.com/markvincze))

## [1.1.15](https://github.com/Travix-International/appix/tree/1.1.15) (2017-05-11)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.13...1.1.15)

**Merged pull requests:**

- Scaffold supports other templates and read from new repo [\#74](https://github.com/Travix-International/appix/pull/74) ([jvantillo](https://github.com/jvantillo))
- Increased coverage on logger package + minor tweaks [\#72](https://github.com/Travix-International/appix/pull/72) ([jvantillo](https://github.com/jvantillo))
- Fix: vet indicates issue in body being used before we check for error [\#69](https://github.com/Travix-International/appix/pull/69) ([jvantillo](https://github.com/jvantillo))
- Added some suggested Makefile additions. [\#68](https://github.com/Travix-International/appix/pull/68) ([jvantillo](https://github.com/jvantillo))
- Removed executable debug file [\#67](https://github.com/Travix-International/appix/pull/67) ([jvantillo](https://github.com/jvantillo))
- Hotfix remove private data [\#63](https://github.com/Travix-International/appix/pull/63) ([jackTheRipper](https://github.com/jackTheRipper))
- \[TR-13172\] Unit test for appcatalog + auth [\#62](https://github.com/Travix-International/appix/pull/62) ([jackTheRipper](https://github.com/jackTheRipper))

## [1.1.13](https://github.com/Travix-International/appix/tree/1.1.13) (2017-02-20)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.14...1.1.13)

## [1.1.14](https://github.com/Travix-International/appix/tree/1.1.14) (2017-02-20)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.12...1.1.14)

**Merged pull requests:**

- TR-13157: Make the logs great again [\#61](https://github.com/Travix-International/appix/pull/61) ([jackTheRipper](https://github.com/jackTheRipper))

## [1.1.12](https://github.com/Travix-International/appix/tree/1.1.12) (2017-02-20)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.11...1.1.12)

**Merged pull requests:**

- Remove accidentally commited debug println [\#60](https://github.com/Travix-International/appix/pull/60) ([markvincze](https://github.com/markvincze))

## [1.1.11](https://github.com/Travix-International/appix/tree/1.1.11) (2017-02-20)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.10...1.1.11)

**Merged pull requests:**

- Not retry appcatalog push/submit if status code != 500 [\#58](https://github.com/Travix-International/appix/pull/58) ([alexmiranda](https://github.com/alexmiranda))
- TR-13122 Increase timeouts [\#57](https://github.com/Travix-International/appix/pull/57) ([jackTheRipper](https://github.com/jackTheRipper))
- TR-11443 Replace the logger instance creation with IoC [\#56](https://github.com/Travix-International/appix/pull/56) ([markvincze](https://github.com/markvincze))
- Logging for appix [\#52](https://github.com/Travix-International/appix/pull/52) ([jackTheRipper](https://github.com/jackTheRipper))

## [1.1.10](https://github.com/Travix-International/appix/tree/1.1.10) (2017-02-13)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.9...1.1.10)

**Fixed bugs:**

- v1.1.7 breaks AppCatalog push \(assuming via watch\) [\#50](https://github.com/Travix-International/appix/issues/50)

**Merged pull requests:**

- Fix retries for appix [\#55](https://github.com/Travix-International/appix/pull/55) ([jackTheRipper](https://github.com/jackTheRipper))
- Fix problem with reading .appixignore [\#54](https://github.com/Travix-International/appix/pull/54) ([alexmiranda](https://github.com/alexmiranda))

## [1.1.9](https://github.com/Travix-International/appix/tree/1.1.9) (2017-02-09)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.8...1.1.9)

**Merged pull requests:**

- Fix glob patterns [\#53](https://github.com/Travix-International/appix/pull/53) ([alexmiranda](https://github.com/alexmiranda))

## [1.1.8](https://github.com/Travix-International/appix/tree/1.1.8) (2017-02-07)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.7...1.1.8)

**Merged pull requests:**

- Safely ignore any path errors when reading `auth.json` file [\#51](https://github.com/Travix-International/appix/pull/51) ([alexmiranda](https://github.com/alexmiranda))

## [1.1.7](https://github.com/Travix-International/appix/tree/1.1.7) (2017-02-03)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.6...1.1.7)

**Closed issues:**

- appix submit issue [\#45](https://github.com/Travix-International/appix/issues/45)

**Merged pull requests:**

- Retry appcatalog calls [\#49](https://github.com/Travix-International/appix/pull/49) ([alexmiranda](https://github.com/alexmiranda))

## [1.1.6](https://github.com/Travix-International/appix/tree/1.1.6) (2017-02-02)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.5...1.1.6)

**Implemented enhancements:**

- Changes IgnoreFilePath to ignore files and/or folders that contain the ignoredFileName in it [\#47](https://github.com/Travix-International/appix/pull/47) ([mAiNiNfEcTiOn](https://github.com/mAiNiNfEcTiOn))

**Closed issues:**

- When running ./build.sh I get an error [\#46](https://github.com/Travix-International/appix/issues/46)

**Merged pull requests:**

- Fixes clone's destination path on the README.md [\#48](https://github.com/Travix-International/appix/pull/48) ([mAiNiNfEcTiOn](https://github.com/mAiNiNfEcTiOn))
- Update appcatalog routes [\#44](https://github.com/Travix-International/appix/pull/44) ([alexmiranda](https://github.com/alexmiranda))

## [1.1.5](https://github.com/Travix-International/appix/tree/1.1.5) (2016-12-22)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.4...1.1.5)

**Merged pull requests:**

- Format http log [\#42](https://github.com/Travix-International/appix/pull/42) ([alexmiranda](https://github.com/alexmiranda))

## [1.1.4](https://github.com/Travix-International/appix/tree/1.1.4) (2016-12-21)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.3...1.1.4)

**Closed issues:**

- Using `whoami` without authenticating displays an error message [\#23](https://github.com/Travix-International/appix/issues/23)

**Merged pull requests:**

- TR-12984 Simplify package structure [\#40](https://github.com/Travix-International/appix/pull/40) ([markvincze](https://github.com/markvincze))
- Log server response headers [\#39](https://github.com/Travix-International/appix/pull/39) ([alexmiranda](https://github.com/alexmiranda))
- Ignore files changes with `appix watch` [\#38](https://github.com/Travix-International/appix/pull/38) ([alexmiranda](https://github.com/alexmiranda))
- Not print all ignored files inside ignored folders [\#37](https://github.com/Travix-International/appix/pull/37) ([alexmiranda](https://github.com/alexmiranda))
- Fix a typo in the install script [\#32](https://github.com/Travix-International/appix/pull/32) ([markvincze](https://github.com/markvincze))
- Adjust the install scripts to the extended Json response of Github [\#31](https://github.com/Travix-International/appix/pull/31) ([markvincze](https://github.com/markvincze))

## [1.1.3](https://github.com/Travix-International/appix/tree/1.1.3) (2016-12-06)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.2...1.1.3)

**Merged pull requests:**

- Refresh the token only if it's near expiry [\#30](https://github.com/Travix-International/appix/pull/30) ([markvincze](https://github.com/markvincze))

## [1.1.2](https://github.com/Travix-International/appix/tree/1.1.2) (2016-12-05)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.1...1.1.2)

**Merged pull requests:**

- Make authentication optional, print warning if user is not logged in [\#29](https://github.com/Travix-International/appix/pull/29) ([markvincze](https://github.com/markvincze))

## [1.1.1](https://github.com/Travix-International/appix/tree/1.1.1) (2016-12-01)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.1.0...1.1.1)

**Implemented enhancements:**

- whoami tells user how to login [\#28](https://github.com/Travix-International/appix/pull/28) ([alexmiranda](https://github.com/alexmiranda))

**Fixed bugs:**

- Fix submit register method requires user to be logged in [\#26](https://github.com/Travix-International/appix/pull/26) ([alexmiranda](https://github.com/alexmiranda))

**Merged pull requests:**

- Add documentation about releasing for maintainers [\#27](https://github.com/Travix-International/appix/pull/27) ([markvincze](https://github.com/markvincze))

## [1.1.0](https://github.com/Travix-International/appix/tree/1.1.0) (2016-11-30)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.0.7...1.1.0)

**Merged pull requests:**

- Add auth to appix verbs [\#24](https://github.com/Travix-International/appix/pull/24) ([alexmiranda](https://github.com/alexmiranda))
- Clean up and restructuring [\#22](https://github.com/Travix-International/appix/pull/22) ([fahad19](https://github.com/fahad19))
- Changelogs [\#21](https://github.com/Travix-International/appix/pull/21) ([fahad19](https://github.com/fahad19))
- Authentication [\#20](https://github.com/Travix-International/appix/pull/20) ([fahad19](https://github.com/fahad19))

## [1.0.7](https://github.com/Travix-International/appix/tree/1.0.7) (2016-11-16)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.0.6...1.0.7)

**Merged pull requests:**

- Show specific version query url on app submit [\#19](https://github.com/Travix-International/appix/pull/19) ([alexmiranda](https://github.com/alexmiranda))

## [1.0.6](https://github.com/Travix-International/appix/tree/1.0.6) (2016-11-14)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.0.5...1.0.6)

**Merged pull requests:**

- Print the error detail if the polling fails [\#18](https://github.com/Travix-International/appix/pull/18) ([markvincze](https://github.com/markvincze))

## [1.0.5](https://github.com/Travix-International/appix/tree/1.0.5) (2016-11-01)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.0.4...1.0.5)

**Merged pull requests:**

- Use TLS with WebSockets [\#16](https://github.com/Travix-International/appix/pull/16) ([markvincze](https://github.com/markvincze))

## [1.0.4](https://github.com/Travix-International/appix/tree/1.0.4) (2016-10-14)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.0.3...1.0.4)

**Merged pull requests:**

- Implement smarter state handling when watching for file changes [\#15](https://github.com/Travix-International/appix/pull/15) ([markvincze](https://github.com/markvincze))
- Implement Livereload capability with WebSockets [\#14](https://github.com/Travix-International/appix/pull/14) ([markvincze](https://github.com/markvincze))

## [1.0.3](https://github.com/Travix-International/appix/tree/1.0.3) (2016-10-06)
[Full Changelog](https://github.com/Travix-International/appix/compare/1.0.2...1.0.3)

**Merged pull requests:**

- Replace AppVeyor with Travis-CI [\#13](https://github.com/Travix-International/appix/pull/13) ([markvincze](https://github.com/markvincze))

## [1.0.2](https://github.com/Travix-International/appix/tree/1.0.2) (2016-10-05)
[Full Changelog](https://github.com/Travix-International/appix/compare/appix-1.0.1.1...1.0.2)
