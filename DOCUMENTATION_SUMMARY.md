# VPSctl Git & Contribution Documentation Summary

> A complete guide to managing git branches, configuring workflows, and contributing to VPSctl

---

## 📚 Documentation Files Created

### 1. **GIT_WORKFLOW.md** (Comprehensive Guide)
**Purpose:** Complete Git branching strategy and workflow policies for global-level package management

**Contents:**
- 🌳 Branch hierarchy (main, develop, feature/*, bugfix/*, hotfix/*, release/*, docs/*)
- 📝 Branch naming conventions with examples
- 🔄 Step-by-step workflows (feature > bugfix > hotfix > release)
- 👥 External contributor guidance
- 🔐 Protected branches configuration
- 📊 Commit history rules and squashing
- 🚦 CI/CD pipeline requirements
- 📱 Common scenario walkthroughs
- 🔄 Git configuration best practices

**When to use:** Reference for complex workflows, setting up CI/CD, team guidelines

---

### 2. **GIT_CHEAT_SHEET.md** (Quick Reference)
**Purpose:** Fast lookup for common git commands organized by task

**Contents:**
- 🚀 Quick start (setup & daily workflow)
- 📋 Branch commands (create, list, delete, rename)
- ✍️ Commit commands (create, amend, search)
- 🔄 Push/Pull/Fetch operations
- 🔀 Merge and rebase commands
- 🔍 Inspect and diff tools
- 🔧 Undo and recovery options
- 📌 Tags and versioning
- 👥 Fork collaboration
- 🚨 Emergency commands
- ✅ Pre-push checklist
- 📱 Complete workflow examples

**When to use:** Daily development, quick command lookup, troubleshooting

---

### 3. **CONTRIBUTING.md** (Contributor Guidelines)
**Purpose:** How-to guide for external and internal contributors

**Contents:**
- 🚀 Set up development environment
- 📋 Types of contributions (bugs, features, code)
- 💻 Development guidelines (code style, testing)
- 📝 Commit message format with examples
- 🧪 Testing requirements and examples
- 🏗️ How to add new modules
- 📋 Pull request process and checklist
- 🔄 Common workflow examples with steps
- ❓ Getting help resources
- 🎉 Code of conduct

**When to use:** Onboarding contributors, PR requirements, development standards

---

### 4. **GitHub Actions Workflows** (CI/CD)
**Location:** `.github/workflows/`

#### `build.yml` — Build & Test Pipeline
```yaml
✅ Multi-OS testing (Ubuntu, macOS, Windows)
✅ Multi-version Go testing (1.22, 1.23)
✅ Build binaries for all platforms
✅ Run tests with coverage
✅ Security scanning
✅ Auto-release on tags
```

#### `quality.yml` — Code Quality Checks
```yaml
✅ Go vet analysis
✅ Format checking
✅ Cyclomatic complexity limits
✅ Documentation validation
```

---

### 5. **Dependabot Configuration**
**Location:** `.github/dependabot.yml`

```yaml
✅ Daily Go module updates
✅ Weekly GitHub Actions updates
✅ Automatic PR creation with labels
✅ Reviewers and assignees
```

---

### 6. **Setup Script**
**Location:** `scripts/git-setup.sh`

**Automates:**
- Configure git user (name & email)
- Setup useful git aliases (st, co, br, ci, etc.)
- Configure git behavior (pull strategy, push strategy)
- Add upstream remote for forks
- Verify configuration

**Usage:**
```bash
bash scripts/git-setup.sh
```

---

## 🌳 Branch Model Overview

```
┌─────────────────────────────────────────────────────────┐
│                    PRODUCTION                           │
│                   (main branch)                         │
│         stable, tested, release-ready code             │
└─────────────────────────────────────────────────────────┘
                          ↑
                    [Release Tags]
                   v0.1.0, v0.2.0, ...
                          ↑
┌─────────────────────────────────────────────────────────┐
│              PRE-RELEASE / STAGING                       │
│                 (develop branch)                        │
│     integration point, next release features           │
└─────────────────────────────────────────────────────────┘
      ↑           ↑             ↑              ↑
   feature/    bugfix/      hotfix/        release/
   branches    branches     branches       branches
   (from      (from        (from          (from
   develop)   develop)     main)          develop)
```

---

## 📋 Branch Type Quick Reference

| Branch | Source | Target | Merges Back To | Lifecycle | Purpose |
|--------|--------|--------|---|---|---|
| `main` | N/A | N/A | N/A | ♾️ Permanent | **Production** |
| `develop` | `main` | `main` | N/A | ♾️ Permanent | **Integration** |
| `feature/*` | `develop` | `develop` | N/A | 🔄 Temporary | New features |
| `bugfix/*` | `develop` | `develop` | N/A | 🔄 Temporary | Bug fixes |
| `hotfix/*` | `main` | `main` + `develop` | `develop` | 🔄 Temporary | Critical fixes |
| `release/*` | `develop` | `main` | `develop` | 🔄 Temporary | Release prep |
| `docs/*` | `develop` | `develop` | N/A | 🔄 Temporary | Documentation |

---

## ✍️ Naming Conventions

### Standard Format
```
<type>/<issue-number>-<description>

✅ feature/123-docker-support
✅ bugfix/42-cpu-calculation-fix
✅ hotfix/critical-security-patch
✅ release/0.2.0
✅ docs/installation-guide
```

### Rules
- ✅ Lowercase letters only
- ✅ Hyphens separate words
- ✅ Include issue # if available
- ✅ Max 50 characters
- ✅ Use meaningful verbs: add-, fix-, improve-, refactor-, update-

---

## 🔄 Typical Workflow Example

### Contributor: Adding a Feature

```bash
# 1. Setup (first-time)
bash scripts/git-setup.sh

# 2. Create feature branch
git co develop && git pull upstream develop
git co -b feature/456-awesome-feature

# 3. Make changes
# ... edit files ...
go test ./...
golangci-lint run

# 4. Commit
git add .
git commit -m "feat: Add awesome feature"

# 5. Keep in sync
git fetch upstream
git rebase upstream/develop

# 6. Push
git push origin feature/456-awesome-feature

# 7. Create PR on GitHub
#    - Base: develop
#    - Compare: feature/456-awesome-feature
#    - Title: "[FEATURE] Add awesome feature"
```

### Reviewer: Merging PR

```bash
# 1. Review code on GitHub
# 2. Request changes if needed
# 3. Approve when ready
# 4. Merge via GitHub (squash + merge)
# 5. GitHub auto-deletes remote branch
```

### Developer: After Merge

```bash
# 1. Sync and cleanup
git checkout develop
git pull origin develop
git branch -d feature/456-awesome-feature

# 2. Continue with next feature
```

---

## 🚀 Daily Commands Cheat Sheet

```bash
# Check status
git st

# Create feature branch
git co -b feature/123-name develop

# Make changes and commit
git add .
git ci -m "feat: Description"

# Keep in sync
git fetch upstream
git rebase upstream/develop

# Push
git push -u origin feature/branch

# Pre-push checks
go fmt ./...
go test ./...
golangci-lint run
gosec ./...

# View before push
git logs | head -5
git diff origin/develop
```

---

## 🔐 Protected Branches (GitHub Settings)

### `main` Branch
- ✅ Require PR reviews (1+ reviewer)
- ✅ Require CI checks pass
- ✅ Up-to-date before merge
- ✅ Include administrators
- ✅ Allow to dismiss stale reviews

### `develop` Branch
- ✅ Require PR reviews (1+ reviewer)
- ✅ Require CI checks pass
- ✅ Up-to-date before merge
- ⚪ Allow force pushes (optional)

---

## 📊 CI/CD Pipeline (GitHub Actions)

### On Every Push

1. **Build** — Compile code
2. **Test** — Run all tests
3. **Lint** — Code quality check
4. **Coverage** — Code coverage > 70%
5. **Security** — Vulnerability scan

### PR Requirements

- ✅ All checks must pass
- ✅ Minimum 1 code review approval
- ✅ Branch must be up-to-date
- ✅ All conversations resolved

### Automated on Release (v*.*.* tag)

- Builds for: Linux, macOS, Windows
- Creates GitHub release
- Auto-uploads binaries
- Generates release notes

---

## 🆘 Troubleshooting Quick Guide

### "My branch is behind develop"
```bash
git fetch upstream
git rebase upstream/develop
git push origin feature/branch --force-with-lease
```

### "I committed to wrong branch"
```bash
git reset --soft HEAD~1
git stash
git checkout correct-branch
git stash pop
git commit -m "message"
```

### "I need to undo my last commit"
```bash
git undo                    # If alias is set
# or
git reset --soft HEAD~1     # Keep changes
git reset --hard HEAD~1     # Discard changes
```

### "I need to rebase but have conflicts"
```bash
# Fix conflict markers in files
git add fixed-files.go
git rebase --continue
# or abort
git rebase --abort
```

---

## 📞 Questions & Support

- **GitHub Issues:** [github.com/codenoidlab/vpsctl/issues](https://github.com/codenoidlab/vpsctl/issues)
- **Discussions:** [github.com/codenoidlab/vpsctl/discussions](https://github.com/codenoidlab/vpsctl/discussions)
- **Email:** hello@codenoid.com
- **Documentation:** See files below

---

## 📚 Complete Documentation Index

```
.
├── README.md                    ← Project overview & features
├── GIT_WORKFLOW.md             ← Branching strategy (comprehensive)
├── GIT_CHEAT_SHEET.md          ← Quick command reference
├── CONTRIBUTING.md             ← How to contribute
├── .github/
│   ├── workflows/
│   │   ├── build.yml           ← CI/CD: Build & Test
│   │   └── quality.yml         ← CI/CD: Code Quality
│   └── dependabot.yml          ← Automated dependency updates
└── scripts/
    └── git-setup.sh            ← Git environment setup
```

---

## ✅ Implementation Checklist

- [x] GIT_WORKFLOW.md created
- [x] GIT_CHEAT_SHEET.md created
- [x] CONTRIBUTING.md created
- [x] GitHub Actions workflows configured
- [x] Dependabot configured
- [x] Git setup script created
- [x] Branch protection rules documented
- [x] Commit message format defined
- [x] PR templates documented
- [x] Testing requirements defined
- [x] CI/CD pipeline defined
- [x] Release process documented
- [x] Contributing guidelines complete

---

## 🎯 Next Steps for Your Project

1. **Push to GitHub:**
   ```bash
   git add .
   git commit -m "docs: Add comprehensive git workflow documentation"
   git push origin develop
   ```

2. **Configure GitHub:**
   - [ ] Set up branch protection for `main` and `develop`
   - [ ] Configure required status checks
   - [ ] Set up code owners (CODEOWNERS file)
   - [ ] Enable branch auto-delete option

3. **Add CODEOWNERS (optional):**
   ```
   * @codenoidlab
   /internal/modules/dashboard/ @codenoidlab
   /docs/ @codenoidlab
   ```

4. **Configure Branch Rules:**
   - Settings > Branches > Add rule
   - Apply to: `main`, `develop`
   - Enable protection options

5. **Create Issue Template:**
   - Settings > Discussions > Create template
   - Settings > Issues > Create template

---

## 📖 Version & Updates

- **Version:** 1.0
- **Created:** April 2026
- **Last Updated:** April 2026
- **Maintained by:** Codenoid Lab
- **Repository:** github.com/codenoidlab/vpsctl
- **Website:** vpsctl.codenoid.com

---

**All systems ready for professional open-source development!** 🚀
