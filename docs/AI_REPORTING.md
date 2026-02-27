# AI-Powered Pipeline Reporting

Automated CI/CD pipeline analysis using AI (Google Gemini or OpenAI) with Slack notifications. Replaces reading thousands of lines of raw logs with concise, actionable summaries.

---

## Table of Contents
- [Quick Start](#quick-start)
- [How It Works](#how-it-works)
- [Configuration](#configuration)
- [GitLab CI Setup](#gitlab-ci-setup)
- [GitHub Actions Setup](#github-actions-setup)
- [Slack Message Format](#slack-message-format)
- [Variables Reference](#variables-reference)
- [Vertex AI Upgrade Path](#vertex-ai-upgrade-path)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

### 1. Get an AI API Key

Choose your preferred AI provider:

**Google Gemini (default):**
1. Go to [Google AI Studio](https://aistudio.google.com/apikey)
2. Create a new API key
3. Add it as a CI/CD secret named `AI_API_KEY`

**OpenAI:**
1. Go to [OpenAI Platform](https://platform.openai.com/api-keys)
2. Create a new API key
3. Add it as a CI/CD secret named `AI_API_KEY`
4. Set the `AI_PROVIDER` variable to `"openai"`

### 2. Set Up Slack Webhook (Optional)

1. Go to [Slack API: Incoming Webhooks](https://api.slack.com/messaging/webhooks)
2. Create a new webhook for your target channel
3. Add the webhook URL as a CI/CD secret named `SLACK_WEBHOOK_URL`

### 3. Enable in Your Pipeline

**GitLab CI:**

```yaml
# .gitlab-ci.yml
include:
  - project: platform/devsecops-template
    file:
      - /templates/gitlab/base.yml
      - /templates/gitlab/ai-report.yml
      # ... your other templates

variables:
  ENABLE_AI_REPORT: "true"
  # AI_PROVIDER: "openai"       # Optional: defaults to "gemini"
```

**GitHub Actions:**

```yaml
# .github/workflows/ci.yml
jobs:
  # ... your build/test/security jobs ...

  ai-report:
    needs: [build, test, sast, dependency-scan]  # all jobs to analyze
    if: always()
    uses: ./.github/workflows/ai-report.yml  # or your template path
    secrets:
      ai_api_key: ${{ secrets.AI_API_KEY }}
      slack_webhook_url: ${{ secrets.SLACK_WEBHOOK_URL }}
    # with:
    #   ai_provider: "openai"   # Optional: defaults to "gemini"
```

---

## How It Works

### Pipeline Flow

```
source → build → test → package → [deploy] → verify → report → ai-analysis → ai-summary
                                                                     │              │
                                                              AI analyzes       Consolidated
                                                              each report       summary → Slack
```

### Two-Stage Process

1. **ai-analysis** — Collects all pipeline artifact reports (security scans, test results, build output) and sends each to the configured AI provider for individual analysis. Each report gets a structured summary with status, severity, key findings, and recommended actions.

2. **ai-summary** — Aggregates all individual analyses and asks the AI provider for a consolidated pipeline summary. Posts the result to Slack as a single, color-coded message with critical issues, warnings, and what passed cleanly.

### What Gets Analyzed

| Report | Source Job | Description |
|--------|-----------|-------------|
| `secrets-report.json` | secrets-detection (Trivy) | Secret/credential leaks |
| `gitleaks-report.json` | secrets-detection (Gitleaks) | Secret/credential leaks |
| `dependency-scan.json` | dependency-scanning | Vulnerable dependencies |
| `sast-report.json` | sast (Trivy) | Static analysis findings |
| `semgrep.json` | sast (Semgrep) | Static analysis findings |
| `iac-report.json` | iac-security (Trivy) | Infrastructure misconfigurations |
| `polaris.json` | iac-security (Polaris) | Kubernetes best practices |
| `trivy.json` | container-security-scan | Container image vulnerabilities |
| `zap/zap.json` | dast-zap | Dynamic security testing |
| `summary.md` | reporting | Existing aggregated security report |

### What Is Excluded

Deployment stages are **not** analyzed:
- `deploy-staging` (stage)
- `deploy-production` (deploy)
- `post-deployment-monitoring` (operate)

---

## Configuration

### Secrets (CI/CD Settings)

| Secret | Required | Description |
|--------|----------|-------------|
| `AI_API_KEY` | Yes | API key for the configured AI provider (Google AI Studio or OpenAI) |
| `SLACK_WEBHOOK_URL` | No | Slack incoming webhook URL. If not set, reports are saved as artifacts only |

### Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ENABLE_AI_REPORT` | `"false"` | Feature toggle — set to `"true"` to enable |
| `AI_PROVIDER` | `"gemini"` | AI provider to use: `"gemini"` or `"openai"` |
| `AI_MODEL` | Auto per provider | Model override. Defaults: `gemini-2.0-flash` (Gemini), `gpt-4.1-mini` (OpenAI) |
| `AI_API_URL` | Auto per provider | API endpoint override. Defaults are set automatically per provider |

---

## GitLab CI Setup

### Basic Setup

```yaml
include:
  - project: platform/devsecops-template
    file:
      - /templates/gitlab/base.yml
      - /templates/gitlab/build.yml
      - /templates/gitlab/test.yml
      - /templates/gitlab/security/secrets.yml
      - /templates/gitlab/security/dependency.yml
      - /templates/gitlab/security/sast.yml
      - /templates/gitlab/report.yml
      - /templates/gitlab/ai-report.yml          # Add this

variables:
  LANGUAGE: "node"
  ENABLE_SECRETS: "true"
  ENABLE_DEPENDENCY_SCAN: "true"
  ENABLE_SAST: "true"
  ENABLE_AI_REPORT: "true"                       # Enable AI reporting
  # AI_PROVIDER: "openai"                        # Optional: defaults to "gemini"
  # AI_MODEL: "gpt-4.1"                          # Optional: override the model
```

Then add `AI_API_KEY` and optionally `SLACK_WEBHOOK_URL` as CI/CD secrets in GitLab Settings > CI/CD > Variables.

### Monorepo Setup

Works the same — AI reporting analyzes all artifacts from all sub-project jobs:

```yaml
variables:
  ENABLE_AI_REPORT: "true"

build:frontend:
  extends: .build:node
  variables:
    PROJECT_PATH: frontend

build:backend:
  extends: .build:python
  variables:
    PROJECT_PATH: backend

# AI reporting automatically picks up all reports from both sub-projects
```

### Artifacts

| File | Stage | Retention | Description |
|------|-------|-----------|-------------|
| `ai-reports/*.txt` | ai-analysis | 30 days | Individual per-report AI analyses |
| `ai-reports/status.json` | ai-analysis | 30 days | Analysis metadata |
| `ai-summary.md` | ai-summary | 30 days | Consolidated AI summary |

---

## GitHub Actions Setup

### Basic Setup

```yaml
name: CI Pipeline

on:
  push:
    branches: [main]
  pull_request:

jobs:
  build:
    uses: ./.github/workflows/build-node.yml

  test:
    needs: [build]
    uses: ./.github/workflows/test-node.yml

  sast:
    uses: ./.github/workflows/security/sast.yml

  dependency-scan:
    uses: ./.github/workflows/security/dependency.yml

  ai-report:
    needs: [build, test, sast, dependency-scan]
    if: always()                                   # Run even if jobs fail
    uses: ./.github/workflows/ai-report.yml        # Or your template path
    secrets:
      ai_api_key: ${{ secrets.AI_API_KEY }}
      slack_webhook_url: ${{ secrets.SLACK_WEBHOOK_URL }}
    # with:
    #   ai_provider: "openai"
    #   ai_model: "gpt-4.1"
```

### Inputs

| Input | Default | Description |
|-------|---------|-------------|
| `ai_provider` | `"gemini"` | AI provider: `"gemini"` or `"openai"` |
| `ai_model` | Auto per provider | Model override (defaults: `gemini-2.0-flash` for Gemini, `gpt-4.1-mini` for OpenAI) |

---

## Slack Message Format

The Slack notification is a single, color-coded message:

- **Green** — All stages passed, no issues
- **Yellow** — Warnings found, review recommended
- **Red** — Critical issues requiring immediate attention
- **Gray** — AI analysis unavailable (fallback mode)

Example message structure:

```
🚨 Pipeline Summary: myorg/myproject
Branch: main | Commit: abc1234

2 critical issues found, 1 warning

🚨 Critical Issues
- CVE-2024-XXXX: Critical vulnerability in lodash 4.17.20
- Hardcoded API key detected in src/config.js

⚠️ Warnings
- Deprecated API usage in build output

✅ Passed
- Unit tests: 142/142 passed (98% coverage)
- SAST: No issues found
- Container scan: Clean

💡 Recommendation: Fix critical vulnerabilities before merging

[View Pipeline]
```

---

## Variables Reference

### GitLab CI

| Variable | Default | Scope | Description |
|----------|---------|-------|-------------|
| `ENABLE_AI_REPORT` | `"false"` | `base.yml` | Enable AI reporting |
| `AI_PROVIDER` | `"gemini"` | `ai-report.yml` | AI provider: `"gemini"` or `"openai"` |
| `AI_MODEL` | Auto per provider | `ai-report.yml` | Model override (defaults: `gemini-2.0-flash` / `gpt-4.1-mini`) |
| `AI_API_URL` | Auto per provider | `ai-report.yml` | API endpoint override |
| `AI_API_KEY` | — | CI/CD secret | API key for the configured AI provider |
| `SLACK_WEBHOOK_URL` | — | CI/CD secret | Slack webhook URL |

### GitHub Actions

| Input/Secret | Default | Type | Description |
|-------------|---------|------|-------------|
| `ai_provider` | `"gemini"` | Input | AI provider: `"gemini"` or `"openai"` |
| `ai_model` | Auto per provider | Input | Model override (defaults: `gemini-2.0-flash` / `gpt-4.1-mini`) |
| `ai_api_key` | — | Secret (required) | API key for the configured AI provider |
| `slack_webhook_url` | — | Secret (optional) | Slack webhook URL |

---

## Vertex AI Upgrade Path

> **Note:** Vertex AI is specific to Google Gemini. This upgrade path applies only when using `AI_PROVIDER: "gemini"` (the default).

The default Gemini configuration uses Google AI Studio with a simple API key. For organizations requiring Swiss data residency (data stays in Zurich), you can upgrade to Vertex AI:

### Requirements
- Google Cloud account with billing enabled
- Google Cloud project with Vertex AI API enabled
- Service account with `Vertex AI User` role

### Migration

The Gemini model and prompts remain identical — only the API endpoint and authentication change:

| | Google AI Studio (default) | Vertex AI |
|---|---|---|
| **Endpoint** | `generativelanguage.googleapis.com` | `europe-west6-aiplatform.googleapis.com` |
| **Auth** | API key header | OAuth2 bearer token |
| **Data residency** | No guarantee | Zurich (europe-west6) |
| **Cost** | Same | Same |

To switch, override `AI_API_URL` and adjust authentication in the template.

---

## Troubleshooting

### AI analysis skipped — "AI_API_KEY not set"

The API key is not configured as a CI/CD secret.

**Fix:** Add `AI_API_KEY` in GitLab Settings > CI/CD > Variables (or GitHub repo Settings > Secrets).

### Gemini API returns 403

The API key is invalid or the Gemini API is not enabled.

**Fix:**
1. Verify the key at [Google AI Studio](https://aistudio.google.com/apikey)
2. Ensure the Generative Language API is enabled for your project

### OpenAI API returns 401

The API key is invalid or has been revoked.

**Fix:**
1. Verify the key at [OpenAI Platform](https://platform.openai.com/api-keys)
2. Ensure the key has not expired and has sufficient permissions

### OpenAI API returns 429 (Rate Limited / Quota Exceeded)

Too many requests or you have exceeded your usage quota.

**Fix:**
1. Check your [OpenAI usage dashboard](https://platform.openai.com/usage)
2. Ensure your account has billing configured and sufficient credits
3. The template retries automatically (up to 2 times with backoff)

### Gemini API returns 429 (Rate Limited)

Too many requests in a short period. The template retries automatically (up to 2 times with backoff).

**Fix:** If persistent, check your [API quotas](https://ai.google.dev/gemini-api/docs/rate-limits). Paid tier increases limits significantly.

### Slack notification failed

**Common causes:**
- Webhook URL is incorrect or expired
- Slack app was removed from the channel
- Message payload too large

**Fix:**
1. Test the webhook: `curl -X POST -H "Content-Type: application/json" -d '{"text":"test"}' YOUR_WEBHOOK_URL`
2. Regenerate the webhook if expired
3. Check Slack app permissions

### Reports are empty or missing

The AI analysis only processes reports that exist as artifacts from previous jobs. If a security scan is disabled or didn't produce output, it won't be analyzed.

**Fix:** Ensure the relevant scans are enabled (`ENABLE_SECRETS`, `ENABLE_SAST`, etc.) and producing artifacts.

### Large reports are truncated

Reports larger than 500KB are truncated before being sent to the AI provider. The AI is informed about the truncation in the prompt.

This is expected behavior to stay within API limits. The most relevant findings are typically in the first portion of the report.

---

## Cost Estimates

### Google Gemini (gemini-2.0-flash)

| Scenario | Estimated Cost |
|----------|---------------|
| Single pipeline run (5 reports) | ~$0.005 - $0.01 |
| 100 pipelines/day | ~$0.50 - $1.00/day |
| 1000 pipelines/day | ~$5 - $10/day |

Based on Gemini 2.0 Flash pricing: $0.10/1M input tokens, $0.40/1M output tokens.

### OpenAI (gpt-4.1-mini)

| Scenario | Estimated Cost |
|----------|---------------|
| Single pipeline run (5 reports) | ~$0.01 - $0.02 |
| 100 pipelines/day | ~$1.00 - $2.00/day |
| 1000 pipelines/day | ~$10 - $20/day |

Based on GPT-4.1 Mini pricing: $0.40/1M input tokens, $1.60/1M output tokens.
