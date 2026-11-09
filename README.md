# orders-api (team-alpha)

Application **source** repo. Contains code + Dockerfile + CI/CD only — **no Kubernetes
manifests** (clean config/source separation). `ci.yml` validates every PR and push (build,
vet, scans); `cd.yml` builds the image once on `main`/`releases/**`, pushes it to the nonprod
ACR, and writes the new tag into `team-alpha-gitops` (the GitOps repo), which Argo CD then
reconciles.

The app repo never talks to a cluster — it only moves Git state.

## Deployment target: `team-alpha-gitops`

CI builds a single cloud-agnostic image (SemVer tag + immutable digest) that runs unchanged
on every cluster. `team-alpha-gitops` holds all per-cluster targeting and config, under
`apps/orders-api/envs/<env>/<cloud>/<region>/`. CD pins only the sbx cluster leaf
(`SBX_CLUSTER` in `cd.yml`); the `promote` workflow moves the pin up from there. Keep the
code cloud-agnostic: read endpoints and identities from environment variables supplied by
the chart, not per-cloud branches in the source.

## CI checks (`.github/workflows/ci.yml`)

| Job | What it does |
| --- | --- |
| `build-test` | `go vet` + `go build` (Go version comes from `go.mod`) |
| `secret-scan` | gitleaks over the full git history |
| `dependency-scan` | `govulncheck`: vulnerable modules or stdlib that the code actually calls |
| `image-scan` | builds the image (not pushed), then Trivy fails on fixed HIGH/CRITICAL vulnerabilities and Dockerfile misconfigurations |

The scanners are official container images pinned by digest; bump the version and digest together.
To make them gate merges, mark the four jobs as required checks in branch protection (`cd.yml`
does not wait for `ci.yml`). To silence a reviewed false positive, use a `.gitleaksignore` entry
(gitleaks) or a `.trivyignore` file (Trivy) with a comment explaining why.

