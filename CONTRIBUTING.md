## Contributing

[fork]: https://github.com/apecloud/kb-cloud-mcp-server/fork
[pr]: https://github.com/apecloud/kb-cloud-mcp-server/compare
[style]: https://github.com/apecloud/kb-cloud-mcp-server/blob/main/.golangci.yaml

Hi there! We're thrilled that you'd like to contribute to this project. Your help is essential for keeping it great.

Contributions to this project are [released](https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license) to the public under the [project's open source license](LICENSE).

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## Prerequisites for running and testing code

These are one time installations required to be able to test your changes locally as part of the pull request (PR) submission process.

1. Install Go [through download](https://go.dev/doc/install) | [through Homebrew](https://formulae.brew.sh/formula/go)
2. [Install golangci-lint](https://golangci-lint.run/welcome/install/#local-installation)
3. [Install Docker](https://docs.docker.com/get-docker/) (optional, for building and testing Docker images)

## Submitting a pull request

1. [Fork][fork] and clone the repository
2. Make sure the tests pass on your machine: `go test -v ./...`
3. Make sure linter passes on your machine: `golangci-lint run`
4. Create a new branch: `git checkout -b my-branch-name`
5. Make your change, add tests, and make sure the tests and linter still pass
6. Push to your fork and [submit a pull request][pr]
7. Pat yourself on the back and wait for your pull request to be reviewed and merged.

Here are a few things you can do that will increase the likelihood of your pull request being accepted:

- Follow the [style guide][style].
- Write tests.
- Keep your change as focused as possible. If there are multiple changes you would like to make that are not dependent upon each other, consider submitting them as separate pull requests.
- Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).
- Follow the PR title format: `<type>: <subject>` (e.g., `feat: add support for backup listing`)
  - Types: feat, fix, docs, style, refactor, test, chore

## Development workflow

1. Set up your development environment:
   ```bash
   # Clone your fork
   git clone https://github.com/<your-username>/kb-cloud-mcp-server.git
   cd kb-cloud-mcp-server

   # Add upstream remote
   git remote add upstream https://github.com/apecloud/kb-cloud-mcp-server.git

   # Install dependencies
   go mod download
   ```

2. Create a new branch for your changes:
   ```bash
   git checkout -b feat/my-feature
   ```

3. Make your changes and test them:
   ```bash
   # Run tests
   go test -v ./...

   # Run linter
   golangci-lint run

   # Build the binary
   go build -o kb-cloud-mcp-server ./cmd/server
   ```

4. Commit your changes:
   ```bash
   git add .
   git commit -m "feat: add my feature"
   ```

5. Push to your fork and create a pull request

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [GitHub Help](https://help.github.com)
- [KubeBlocks Documentation](https://kubeblocks.io/docs/)