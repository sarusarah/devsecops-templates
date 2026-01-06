# Quick Start Guide - DevSecOps Template Testing

## ✅ Dagger is Now Working!

Your Dagger module is ready to test GitLab CI/CD pipelines locally.

---

## 🚀 Usage

### **Always run from the `dagger/` directory:**

```bash
cd /home/sarah/Projects/Liip/DevSecOps/dagger
```

---

## 📋 Available Commands

### Run Full Security Scan Suite

```bash
# Test Node.js example
dagger call test --source=../examples/node --language=node

# Test Python example
dagger call test --source=../examples/python --language=python

# Test PHP example
dagger call test --source=../examples/php-symfony --language=php
```

### Run Individual Scans

```bash
# Secrets detection
dagger call secrets-detection --source=../examples/node

# Dependency scanning
dagger call dependency-scanning --source=../examples/node --language=node

# SAST (Semgrep)
dagger call sast-scanning --source=../examples/node

# Container scanning (requires image already built)
dagger call container-scanning \
  --image-name=registry.gitlab.com/myproject/app \
  --image-tag=latest
```

### Build & Test Node.js

```bash
# Build
dagger call build-node \
  --source=../examples/node \
  --node-version=20 \
  --package-manager=pnpm \
  export --path=../build-output

# Test
dagger call test-node \
  --source=../examples/node \
  --node-version=20 \
  --package-manager=pnpm
```

### Validate YAML

```bash
dagger call validate-yaml --yaml-file=../.gitlab-ci.yml
```

---

## 🎯 Common Workflows

### Before Committing

```bash
cd dagger
dagger call test --source=.. --language=node
```

### Testing a Specific Example

```bash
cd dagger

# Node.js/Nuxt
dagger call test --source=../examples/node --language=node

# Python
dagger call test --source=../examples/python --language=python

# PHP Symfony
dagger call test --source=../examples/php-symfony --language=php

# PHP Drupal
dagger call test --source=../examples/php-drupal --language=php
```

### Testing Your Own Project

```bash
cd dagger
dagger call test --source=/path/to/your/project --language=node
```

---

## 🔍 What Gets Tested

When you run `dagger call test`, it runs:

1. **Secrets Detection** - Gitleaks scans for hardcoded secrets
2. **Dependency Scanning** - Checks for vulnerable dependencies (npm audit, pip-audit, composer audit)
3. **SAST** - Semgrep static code analysis for security issues

All scans run in parallel for speed!

---

## ⚡ Performance

- **First run:** ~30-60 seconds (downloading images)
- **Subsequent runs:** ~5-10 seconds (cached!)
- **GitLab CI:** 2-5 minutes for comparison

**10x faster iteration!**

---

## 📊 Success Example

```bash
$ cd dagger
$ dagger call test --source=../examples/node --language=node

✔ connect 0.5s
✔ load module: . 0.8s
✔ devsecops: Devsecops! 0.0s
✔ .test(source: Address.directory: Directory!, language: "node"): Void 5.7s

✅ All security scans passed!
```

---

## ❌ Handling Failures

If a scan finds issues:

```bash
✘ .test(...): Void ERROR
✘ withExec gitleaks detect... ERROR
! Gitleaks found 3 potential secrets
! exit code: 1
```

**Action:**
1. Review the findings
2. Fix or whitelist false positives
3. Run again to verify

---

## 🛠️ Troubleshooting

### "module not found"

**Solution:** Always run from `dagger/` directory:
```bash
cd /home/sarah/Projects/Liip/DevSecOps/dagger
```

### "failed to connect to engine"

**Solution:** Ensure Docker is running:
```bash
docker info
# If not running: sudo systemctl start docker
```

### Semgrep requires metrics

**Solution:** Already fixed! We use `p/security-audit` config instead of `auto`.

---

## 🔗 Alternative: gitlab-ci-local

If you prefer testing with gitlab-ci-local:

```bash
cd /home/sarah/Projects/Liip/DevSecOps
./test-local.sh
```

---

## 📚 Documentation

- **Full Testing Guide:** `../TESTING.md`
- **Dagger Details:** `./README.md`
- **All Fixes:** `../FIXES_APPLIED.md`

---

## 🎉 You're Ready!

Start testing your pipelines locally:

```bash
cd /home/sarah/Projects/Liip/DevSecOps/dagger
dagger call test --source=../examples/node --language=node
```

**Happy testing!** 🚀
