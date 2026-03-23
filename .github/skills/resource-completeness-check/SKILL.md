---
name: resource-completeness-check
description: "Compare resourceDefinition.json against completness/existing_tf_resources.txt to report coverage gaps. Shows which Terraform azurerm resources are supported and which are missing. Triggers on: audit, completeness review."
---

# Resource Completeness Check

## Procedure

### 1. Load reference list

```bash
sort completness/existing_tf_resources.txt > /tmp/rc-reference.txt
wc -l /tmp/rc-reference.txt
```

### 2. Load supported resources

```bash
grep '"name":' resourceDefinition.json | sed 's/.*"name": "//;s/".*//' | sort > /tmp/rc-supported.txt
wc -l /tmp/rc-supported.txt
```

### 3. Compare

```bash
# Missing (in reference but not supported)
comm -23 /tmp/rc-reference.txt /tmp/rc-supported.txt > /tmp/rc-missing.txt

# Extra (supported but not in reference — may be valid out-of-doc resources)
comm -13 /tmp/rc-reference.txt /tmp/rc-supported.txt > /tmp/rc-extra.txt
```

### 4. Report

```
## Resource Completeness Report

### Summary
- Reference resources (existing_tf_resources.txt): <count>
- Supported in provider: <count>
- Missing: <missing_count>
- Coverage: <percentage>%

### Missing Resources (<count>)
| # | Resource |
|---|----------|
| 1 | `azurerm_<name>` |
| 2 | `azurerm_<name>` |
...

### Extra Resources (<count>)
Resources in provider but not in reference list (may be out-of-doc):
| # | Resource |
|---|----------|
| 1 | `azurerm_<name>` |
...
```

### 5. Cleanup

```bash
rm -f /tmp/rc-*.txt
```
