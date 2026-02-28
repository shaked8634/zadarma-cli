---
name: forgejo
description: "Forgejo operations via curl + $FORGEJO_API_KEY: issues, PRs, releases, CI/actions. Use when: (1) checking PR status or CI, (2) creating/commenting on issues, (3) managing releases, (4) querying Forgejo API. NOT for: local git operations (use git directly), or when $FORGEJO_API_KEY is not set."
metadata:
  {
    "openclaw":
      {
        "emoji": "🦊",
        "requires": { "bins": ["curl"] },
      },
  }
---

# Forgejo Skill

Use curl + `$FORGEJO_API_KEY` to interact with Forgejo repositories, issues, PRs, and CI.

## When to Use

✅ **USE this skill when:**

- Checking PR status, reviews, or merge readiness
- Viewing CI/actions run status
- Creating, closing, or commenting on issues
- Creating or merging pull requests
- Querying Forgejo API for repository data
- Managing releases
- Listing repos or collaborators

## When NOT to Use

❌ **DON'T use this skill when:**

- Local git operations (commit, push, pull, branch) → use `git` directly
- Cloning repositories → use `git clone`
- Complex multi-file diffs → read files directly
- $FORGEJO_API_KEY is not set

## Setup

```bash
# Set API token (one-time)
export FORGEJO_API_KEY="your_api_token"
export FORGEJO_HOST="https://forgejo.example.com"  # e.g., https://forgejo.o-st.dev

# Verify (if curl returns 200)
curl -s -H "Authorization: token $FORGEJO_API_KEY" "$FORGEJO_HOST/api/v1/user"
```

## Common API Patterns

### Pull Requests

```bash
# List PRs for a repo
curl -s -H "Authorization: token $FORGEJO_API_KEY" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/pulls" | jq

# Get specific PR details
curl -s -H "Authorization: token $FORGEJO_API_KEY" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/pulls/55"

# Create PR
curl -X POST -H "Authorization: token $FORGEJO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"title":"feat: add feature","body":"Description","head":"feature-branch","base":"main"}' \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/pulls"

# Merge PR
curl -X POST -H "Authorization: token $FORGEJO_API_KEY" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/pulls/55/merge"
```

### Issues

```bash
# List open issues
curl -s -H "Authorization: token $FORGEJO_API_KEY" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/issues?state=open"

# Create issue
curl -X POST -H "Authorization: token $FORGEJO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"title":"Bug: something broken","body":"Details..."}' \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/issues"

# Close issue
curl -X PATCH -H "Authorization: token $FORGEJO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"state":"closed"}' \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/issues/42"
```

### Releases

```bash
# List releases
curl -s -H "Authorization: token $FORGEJO_API_KEY" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/releases"

# Create release
curl -X POST -H "Authorization: token $FORGEJO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"tag_name":"v1.0.0","name":"Release v1.0.0","body":"Changelog..."}' \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/releases"

# Upload asset to release
curl -X POST -H "Authorization: token $FORGEJO_API_KEY" \
  -F "attachment=@/path/to/file" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/releases/{release_id}/assets"
```

### Actions (CI/CD Runs)

```bash
# List workflow runs
curl -s -H "Authorization: token $FORGEJO_API_KEY" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/actions/runs"

# Get run details
curl -s -H "Authorization: token $FORGEJO_API_KEY" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/actions/runs/{run_id}"

# Re-run a workflow
curl -X POST -H "Authorization: token $FORGEJO_API_KEY" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/actions/runs/{run_id}/rerun"
```

## JSON Output & Filtering

Use `jq` to filter API responses:

```bash
# List all PR numbers and titles
curl -s -H "Authorization: token $FORGEJO_API_KEY" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/pulls" | \
  jq '.[] | "\(.number): \(.title)"'

# Get PR count by state
curl -s -H "Authorization: token $FORGEJO_API_KEY" \
  "$FORGEJO_HOST/api/v1/repos/owner/repo/pulls?state=all" | \
  jq 'group_by(.state) | map({state: .[0].state, count: length})'
```

## Notes

- Always set `$FORGEJO_API_KEY` and `$FORGEJO_HOST` before using curl
- Forgejo API paths follow `/api/v1/repos/owner/repo/{resource}`
- Use `-H "Content-Type: application/json"` for POST/PATCH requests
- Use `-d` for data (JSON format) and `-F` for file uploads
- HTTP status codes: 200=OK, 201=Created, 204=No Content, 400=Bad Request, 401=Unauthorized, 404=Not Found
