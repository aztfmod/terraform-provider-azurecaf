---
name: resource-diff-report
description: "Compare two versions of resourceDefinition.json (e.g., branch vs main) and produce a structured change summary. Triggers on: PR review, audit, before/after comparison."
---

# Resource Diff Report

## Procedure

### 1. Get the two versions

**Branch vs main:**
```bash
git show origin/main:resourceDefinition.json > /tmp/rd-main.json
cp resourceDefinition.json /tmp/rd-branch.json
```

**Or two arbitrary commits:**
```bash
git show <commit1>:resourceDefinition.json > /tmp/rd-before.json
git show <commit2>:resourceDefinition.json > /tmp/rd-after.json
```

### 2. Extract resource names

```bash
grep '"name":' /tmp/rd-main.json | sed 's/.*"name": "//;s/".*//' | sort > /tmp/rd-names-main.txt
grep '"name":' /tmp/rd-branch.json | sed 's/.*"name": "//;s/".*//' | sort > /tmp/rd-names-branch.txt
```

### 3. Identify changes

```bash
# New resources (in branch but not in main)
comm -13 /tmp/rd-names-main.txt /tmp/rd-names-branch.txt > /tmp/rd-added.txt

# Removed resources (in main but not in branch)
comm -23 /tmp/rd-names-main.txt /tmp/rd-names-branch.txt > /tmp/rd-removed.txt

# Common resources (check for field changes)
comm -12 /tmp/rd-names-main.txt /tmp/rd-names-branch.txt > /tmp/rd-common.txt
```

### 4. Report

```
## Resource Definition Changes

### Added (<count>)
| Resource | Slug | Length | Scope |
|----------|------|--------|-------|
| <name>   | <slug> | <min>-<max> | <scope> |

### Removed (<count>)
| Resource |
|----------|
| <name>   |

### Modified (<count>)
| Resource | Field | Before | After |
|----------|-------|--------|-------|
| <name>   | <field> | <old> | <new> |

### Summary
- Total resources: <before> → <after>
- Net change: +<added> / -<removed> / ~<modified>
```

### 5. Cleanup

```bash
rm -f /tmp/rd-*.json /tmp/rd-*.txt
```
