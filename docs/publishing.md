# Publishing the guides

The workflow in `.github/workflows/pages.yml` deploys only the artifact produced
by `scripts/verify-docs.sh`. GitHub Pages is enabled with **Settings → Pages →
Build and deployment → Source** set to **GitHub Actions**. Review the first
deployment URL emitted by the `deploy` job and record it here. Protect the
`github-pages` environment according to the repository's release policy.
