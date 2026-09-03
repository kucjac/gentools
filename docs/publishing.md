# Publishing the guides

The workflow in `.github/workflows/pages.yml` deploys only the artifact produced
by `scripts/verify-docs.sh`. GitHub Pages is enabled with **Settings → Pages →
Build and deployment → Source** set to **GitHub Actions**. Review the first
deployment URL emitted by the `deploy` job and record it here. Protect the
`github-pages` environment according to the repository's release policy.

## Verified deployment

Verified on 2026-09-03: [https://kucjac.github.io/gentools/](https://kucjac.github.io/gentools/)
from the successful [Pages deployment workflow](https://github.com/kucjac/gentools/actions/runs/33728206349).
