# Release process

1. Update `CHANGELOG.md`.
2. Ensure CI is green on main branch.
3. Create tag: `git tag v1.0.0 && git push origin v1.0.0`.
4. GitHub Actions `release.yml` builds artifacts for Windows/macOS/Linux and uploads them.
5. Verify generated GitHub Release notes and attached artifacts.
