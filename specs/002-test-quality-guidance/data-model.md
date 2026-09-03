# Data Model: Test Quality Guidance

## Behavior inventory entry

Canonical file: `docs/testing/behavior-inventory.json`.

| Field | Meaning | Validation |
|-------|---------|------------|
| `id` | Stable behavior identifier | Unique; immutable once reviewed |
| `symbol` | Exported declaration or named high-risk internal behavior | Resolves to a current source declaration |
| `owner` | Package responsible for the behavior | Matches a repository package |
| `risk` | `high`, `medium`, or `low` | High-risk entries cannot remain unclassified |
| `test_level` | `unit`, `integration`, `compatibility`, `example`, or `excluded` | Matches the behavior's appropriate evidence type |
| `status` | `verified`, `gap`, or `excluded` | `excluded` requires rationale; `gap` requires priority |
| `evidence` | Test, scenario, or example path | Referenced path exists |
| `rationale` | Why the classification is appropriate | Required for excluded and integration-only behavior |
| `priority` | Remediation order | Required for gaps |

State transitions: a new symbol begins as `gap`; review changes it to
`verified` only after evidence exists, or `excluded` with a rationale. Source
removal retires the entry in the same change. An evidence-path failure returns
the entry to `gap`.

## Integration scenario

| Field | Meaning | Validation |
|-------|---------|------------|
| `name` | Stable scenario name | One each for inspect, generated, cross-package, and invalid |
| `source_input` | Tracked fixture source | Included in candidate snapshot |
| `module_boundary` | Consumer module and source-under-test relationship | Candidate snapshot is explicit; release smoke has no replacement |
| `expected_result` | Observable output or assertion | Positive scenario passes only when it matches |
| `failure_expectation` | Diagnostic or non-zero outcome | Invalid scenario must not report success |
| `mutation_check` | Deliberate altered expectation | The scenario fails when altered |

## Runnable example

| Field | Meaning | Validation |
|-------|---------|------------|
| `use_case` | User goal demonstrated | Referenced from the Pages guide |
| `input` | Source package to inspect | Present in a clean checkout |
| `command` | Reproducible execution | Runs with selected toolchain |
| `output` | Deterministic generated source | Compared with committed expectation |
| `diagnostic_path` | Troubleshooting route | Linked from the guide |

## Publication artifact

The static site has a source directory, validated artifact directory, commit
identity, and deployment URL. Only the validated artifact may be deployed; the
Pages job is the sole holder of publication permissions.
