# Releasing lazycap

This document describes the release process for lazycap. The project uses trunk-based development with automated releases on every push to `main`.

## Table of Contents

- [Overview](#overview)
- [Automatic Releases](#automatic-releases)
- [Version Bumping](#version-bumping)
- [Semantic Versioning](#semantic-versioning)
- [Homebrew Distribution](#homebrew-distribution)
- [Manual Releases](#manual-releases)
- [Troubleshooting](#troubleshooting)

## Overview

lazycap follows a simple trunk-based development workflow:

1. **Development**: Work on `dev` branch or feature branches
2. **Merge**: Merge changes to `main` via pull request
3. **Release**: GitHub Actions automatically creates a release

Every push to `main` that modifies Go source files triggers the release workflow, which:

1. Analyzes commit messages to determine version bump type
2. Creates and pushes a new version tag
3. Builds binaries for all supported platforms using GoReleaser
4. Creates a GitHub Release with assets and changelog
5. Updates the Homebrew formula

## Automatic Releases

The release workflow triggers on pushes to `main` that modify:

- `**.go` (any Go source file)
- `go.mod` or `go.sum`
- `.goreleaser.yaml`
- `.github/workflows/release.yml`

**No manual intervention is required** for standard releases. Just merge your PR and the release happens automatically.

## Version Bumping

The release workflow automatically determines the version bump based on commit messages since the last tag.

### Priority Order

1. **Manual override** (highest priority): Include `[major]`, `[minor]`, or `[patch]` in any commit message
2. **Breaking changes**: `BREAKING CHANGE` in commit body or `!` after type (e.g., `feat!:`)
3. **Features**: Commits starting with `feat:` or `feat(scope):`
4. **Default**: Patch version bump

### Examples

| Commit Message | Version Bump |
|----------------|--------------|
| `fix: correct device detection` | Patch (0.0.X) |
| `docs: update README` | Patch (0.0.X) |
| `feat: add Firebase plugin` | Minor (0.X.0) |
| `feat(ui): new settings panel` | Minor (0.X.0) |
| `fix!: change config format` | Major (X.0.0) |
| `feat!: redesign plugin API` | Major (X.0.0) |
| `refactor: cleanup [minor]` | Minor (0.X.0) |
| `fix: bug fix [major]` | Major (X.0.0) |

### Forcing a Specific Version Bump

To force a specific version bump regardless of commit content, include one of these tags in your commit message:

```bash
# Force major version bump (1.2.3 -> 2.0.0)
git commit -m "refactor: major API changes [major]"

# Force minor version bump (1.2.3 -> 1.3.0)
git commit -m "chore: add new capability [minor]"

# Force patch version bump (1.2.3 -> 1.2.4)
git commit -m "feat: small addition [patch]"
```

## Semantic Versioning

lazycap follows [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** (X.0.0): Incompatible API changes, breaking changes for users
- **MINOR** (0.X.0): New features that are backward compatible
- **PATCH** (0.0.X): Bug fixes and minor improvements

### When to Use Each

| Change Type | Version Bump |
|-------------|--------------|
| Breaking CLI changes | Major |
| Breaking plugin API changes | Major |
| Removed features | Major |
| New commands or flags | Minor |
| New plugins | Minor |
| New keyboard shortcuts | Minor |
| Bug fixes | Patch |
| Performance improvements | Patch |
| Documentation updates | Patch |
| Dependency updates | Patch |

## Homebrew Distribution

lazycap is distributed via Homebrew using a tap in the same repository.

### Installation

```bash
# Add the tap (one-time setup)
brew tap icarus-itcs/lazycap https://github.com/icarus-itcs/lazycap

# Install lazycap
brew install lazycap
```

### Updating

```bash
# Update all taps and upgrade
brew update && brew upgrade lazycap

# Or force reinstall
brew reinstall lazycap
```

### How It Works

The Homebrew formula (`Formula/lazycap.rb`) is automatically updated by the release workflow:

1. Release workflow creates a new tag
2. GoReleaser builds and publishes the release
3. Workflow calculates SHA256 of the source tarball
4. Formula is updated with new URL and SHA256
5. Changes are committed to `main`

**Note**: The formula update commit does not trigger another release because it only modifies the `.rb` file, not Go source files.

### Limitations

Since the tap is in the main repository (not a separate `homebrew-lazycap` repo), users need to explicitly specify the full URL when adding the tap:

```bash
# This works
brew tap icarus-itcs/lazycap https://github.com/icarus-itcs/lazycap

# This does NOT work (expects homebrew-lazycap repo)
brew tap icarus-itcs/lazycap
```

## Manual Releases

In rare cases, you may need to create a release manually.

### Prerequisites

- Go 1.21+
- [GoReleaser](https://goreleaser.com/) installed
- `GITHUB_TOKEN` environment variable set

### Steps

1. **Create and push a tag**:
   ```bash
   git tag -a v1.2.3 -m "Release v1.2.3"
   git push origin v1.2.3
   ```

2. **Run GoReleaser locally** (optional, workflow will handle it):
   ```bash
   goreleaser release --clean
   ```

3. **Update Homebrew formula manually** (if needed):
   ```bash
   # Get SHA256
   curl -sL https://github.com/icarus-itcs/lazycap/archive/refs/tags/v1.2.3.tar.gz | sha256sum

   # Update Formula/lazycap.rb with new url and sha256
   # Commit and push
   ```

### Dry Run

To test the release process without publishing:

```bash
goreleaser release --snapshot --clean
```

## Troubleshooting

### Release Not Triggered

**Symptoms**: Pushed to `main` but no release was created.

**Solutions**:
1. Check if the commit modified Go source files (`**.go`, `go.mod`, `go.sum`)
2. Check GitHub Actions for workflow errors
3. Verify the tag doesn't already exist

### Tag Already Exists

**Symptoms**: "Tag already exists, skipping release" in workflow logs.

**Solution**: This is expected behavior. If you need to re-release the same version:
1. Delete the tag: `git push --delete origin v1.2.3`
2. Delete the GitHub release
3. Push to `main` again

### GoReleaser Errors

**Symptoms**: GoReleaser step fails in workflow.

**Solutions**:
1. Run `goreleaser check` locally to validate config
2. Check `.goreleaser.yaml` syntax
3. Ensure all build targets are valid

### Homebrew Formula Issues

**Symptoms**: `brew install` fails or installs wrong version.

**Solutions**:
1. Update the tap: `brew update`
2. Check `Formula/lazycap.rb` has correct URL and SHA256
3. Try reinstalling: `brew reinstall lazycap`

### Build Failures

**Symptoms**: Builds fail for specific platforms.

**Solutions**:
1. Test locally: `GOOS=linux GOARCH=arm64 go build`
2. Check for CGO dependencies (should be CGO_ENABLED=0)
3. Verify ignore rules in `.goreleaser.yaml`

---

## Release Checklist

For maintainers preparing a release:

- [ ] All tests pass on `dev` branch
- [ ] Code reviewed and approved
- [ ] Commit messages follow conventional commits format
- [ ] Breaking changes documented in commit body with `BREAKING CHANGE:`
- [ ] Version bump type is appropriate for changes
- [ ] PR merged to `main`
- [ ] Release workflow completed successfully
- [ ] GitHub Release created with correct version
- [ ] Homebrew formula updated automatically

## Questions?

If you have questions about the release process, please [open an issue](https://github.com/icarus-itcs/lazycap/issues).
