# Contract: Quality Artifacts

## `testinventory assess`

Development-only command contract:

```text
go run ./cmd/testinventory assess \
  --inventory docs/testing/behavior-inventory.json \
  --report docs/testing/behavior-inventory.md \
  --coverage /tmp/gentools.coverprofile
```

The command must:

1. Discover exported declarations in supported library packages.
2. Fail if an exported declaration has no inventory entry or an entry resolves
   to no declaration.
3. Fail if cited evidence does not exist.
4. Render a deterministic report grouping verified work, exclusions, and gaps by
   risk and test level.
5. Return non-zero for unclassified public behavior or high-risk gaps; report
   lower-risk accepted gaps without hiding them.

It must not alter source code, fetch dependencies, publish artifacts, or claim
that coverage proves a behavior is tested.

## Consumer scenario runner

The integration runner accepts a scenario name or `all`. It snapshots tracked
candidate source into a temporary directory, runs the named consumer module
against that snapshot, prints the scenario identity and source boundary, and
removes the temporary data on completion. It must fail closed if the snapshot,
fixture inputs, or expected diagnostic is missing.

The release smoke runner accepts an explicit released module version. It must
run without a local replacement and label its output as published-module
evidence. It is not a pull-request substitute.

## Documentation validator

The documentation validator must run every linked example, compare its expected
output, verify local links and source links in the static site, and prepare the
deployable directory. A failed example or link prevents the Pages artifact from
being uploaded or deployed.
