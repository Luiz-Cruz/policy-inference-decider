# Coverage check and branch protection

The **Coverage** workflow runs on every pull request targeting `main` or `develop` when the PR is **not in draft**. It runs again on every new push to the PR.

To block merge until coverage is at least 90%:

1. **Settings** → **Branches** → **Branch protection rules**
2. Add or edit a rule for **main** and another for **develop**.
3. Enable **Require status checks to pass before merging**.
4. Add the status check named **Coverage** (or **coverage**).
5. Save.

After that, PRs into `main` or `develop` cannot be merged until the Coverage workflow passes (and any other required checks). Draft PRs do not run the workflow; when the PR is marked "Ready for review", the next push runs it.
