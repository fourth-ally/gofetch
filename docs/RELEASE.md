# Release New Version

Prepare, commit, tag, and publish a new `gofetch-wasm` release. Version example: `1.0.12`

---

# Release Steps

Run the following commands in order:

```bash
# 1. See commits since last release
git log --oneline <last-version-tag>..HEAD

# 2. Update CHANGELOG.md at the ƒtop with:
cat <<EOL >> CHANGELOG.md
## [X.Y.Z] – YYYY-MM-DD

### Added
- ...

### Changed
- ...

### Fixed
- ...

### Removed
- ...
EOL

# 3. Commit all changes
git add .
git commit -m "chore(release): vX.Y.Z"

# 4. Tag the release
git tag -a vX.Y.Z -m "vX.Y.Z"

# 5. Push commit and tag
git push origin master --tags
