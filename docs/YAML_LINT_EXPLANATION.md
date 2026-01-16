# YAML Lint Error Explanation

## The Error You Saw

```
templates/report.yml
  Warning: 4:1 [document-start] missing document start "---"
Error: Process completed with exit code 1.
```

## Why This Happens

### 1. **The Linter Configuration**

Looking at your GitHub Actions workflow (`.github/workflows/ci.yml`), line 33:

```yaml
yamllint -d "{extends: default, rules: {line-length: {max: 200}}}" "$f"
```

This uses `yamllint` with `extends: default`, which means it inherits the **default yamllint rules**.

### 2. **The Default Rules Include `document-start`**

The default yamllint configuration includes a rule called `document-start` which **requires** all YAML files to start with `---`.

From [yamllint documentation](https://yamllint.readthedocs.io/en/stable/rules.html#module-yamllint.rules.document_start):

> **document-start**: Use this rule to control the use of document start marker (---).
> 
> **Default**: `present: true` (document start marker is required)

### 3. **Why Is This a Rule?**

The `---` marker serves several purposes in YAML:

#### **a) Multi-Document Files**
YAML allows multiple documents in a single file, separated by `---`:

```yaml
---
# Document 1
name: first
---
# Document 2
name: second
```

#### **b) Explicit Document Boundary**
It makes it clear where the YAML document starts, which helps:
- **Parsers**: Some tools rely on it for accurate parsing
- **Humans**: Makes the structure obvious when reading
- **Concatenation**: Safely combine multiple YAML files

#### **c) Prevents Ambiguity**
Consider this file:

```yaml
# This is a comment
key: value
```

vs.

```yaml
---
# This is a comment
key: value
```

The second version is unambiguous - it's clearly a YAML document, not just text that happens to look like YAML.

### 4. **It's a Style Rule, Not a Syntax Requirement**

**Important**: `---` is **NOT required** by the YAML specification itself. YAML is perfectly valid without it:

- ✅ Valid YAML: `key: value`
- ✅ Also valid YAML: `---\nkey: value`

However, your **linter's configuration** enforces it as a **style/best practice rule**.

## Why Your CI Uses This Rule

Looking at your workflows, they enforce this for:

1. **Consistency**: All YAML files follow the same format
2. **GitLab CI compatibility**: GitLab CI often uses multiple YAML documents
3. **Template composition**: Your templates are included/combined, so explicit boundaries help
4. **Industry best practice**: Many style guides recommend it

## The Two Different Configurations

You have **two** separate linting setups:

### **GitLab CI** (`.gitlab/ci/lint.yml`)
```yaml
yamllint -d "{extends: default, rules: {line-length: {max: 200}, comments: disable}}" "$file"
```
- Extends `default` (includes `document-start` rule)
- Disables `comments` rule
- Sets `line-length` to 200

### **GitHub Actions** (`.github/workflows/ci.yml`)
```yaml
yamllint -d "{extends: default, rules: {line-length: {max: 200}}}" "$f"
```
- Extends `default` (includes `document-start` rule)
- Sets `line-length` to 200
- Does **NOT** disable comments (difference from GitLab)

Both require `---` because both extend the `default` config.

## How to Customize (If You Want)

If you wanted to make `---` optional instead, you could change the configuration:

```yaml
yamllint -d "{extends: default, rules: {document-start: disable, line-length: {max: 200}}}" "$f"
```

But **I don't recommend this** because:
1. It's a widely-accepted best practice
2. Your templates are designed for composition/inclusion
3. GitLab CI benefits from explicit document markers
4. It makes your templates more professional and consistent

## Real-World Impact

### Without `---`:
```yaml
# templates/security/secrets.yml
secrets-detection:
  stage: source
  script: ...
```

### With `---`:
```yaml
---
# templates/security/secrets.yml
secrets-detection:
  stage: source
  script: ...
```

When these templates are included in a `.gitlab-ci.yml`, the explicit document markers help GitLab CI:
- Correctly parse template boundaries
- Handle includes more reliably
- Provide better error messages when there are issues

## Summary

**Why the error occurred:**
- Your linter uses `yamllint` with `extends: default`
- The default config includes `document-start: {present: true}`
- This **requires** all YAML files to start with `---`

**Why this rule exists:**
- Best practice for YAML files
- Helps with multi-document scenarios
- Makes document boundaries explicit
- Improves parser reliability

**Is `---` required by YAML?**
- **No** - YAML syntax doesn't require it
- **Yes** - Your linter configuration requires it as a style rule

**Should you keep this rule?**
- **Yes** - It's a good practice, especially for GitLab CI templates
- Your templates are now compliant and follow industry best practices

## Further Reading

- [yamllint documentation on document-start](https://yamllint.readthedocs.io/en/stable/rules.html#module-yamllint.rules.document_start)
- [YAML 1.2 Specification - Document Markers](https://yaml.org/spec/1.2.2/#chapter-9-document-stream-productions)
- [GitLab CI YAML Style Guide](https://docs.gitlab.com/ee/development/cicd/templates.html#yaml-style-guide)
