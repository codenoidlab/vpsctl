# VPSctl Git Quick Reference

> Fast lookup for common git commands and workflows

---

## 🚀 Quick Start

```bash
# Setup (one-time)
bash scripts/git-setup.sh

# Daily workflow
git st                           # Status
git co feature/my-feature        # Checkout branch
git add .                        # Stage changes
git ci -m "feat: description"    # Commit
git push origin feature/branch   # Push

# Sync with upstream (for forks)
git sync                         # Equivalent to: fetch + rebase
```

---

## 📋 Branch Quick Commands

```bash
# Create and checkout new branch
git co -b feature/123-name          # From current branch
git co -b feature/123 develop       # Specify source branch

# List branches
git br                              # Local branches
git br -a                           # All branches (local + remote)
git br -r                           # Remote branches only

# Delete branches
git br -d feature/123-name          # Delete local (safe)
git br -D feature/123-name          # Force delete local
git push origin --delete feature    # Delete remote

# Rename branch
git br -m old-name new-name         # Rename locally
git push origin -u new-name         # Push new name
git push origin --delete old-name   # Delete old name

# Switch branches
git co feature/my-feature           # Switch to existing
git co -                            # Switch to previous branch
git co develop                      # Master branch
```

---

## ✍️ Commits

```bash
# Basic commit
git ci -m "feat: Add new feature"

# Amend last commit (no new commit)
git amend                           # Requires: git config --global alias.amend

# Amend with changes
git add .
git ci --amend --no-edit

# Commit specific files
git ci -- file1.go file2.go -m "feat: Update specific files"

# Staged vs unstaged
git add .                           # Stage all
git add file.go                     # Stage specific file
git add -p                          # Stage interactively (patch mode)

# View commits
git logs                            # Pretty log (requires alias)
git log --oneline                   # Simple log
git log -p                          # Log with diffs
git log --grep="search"             # Search commit messages
```

---

## 🔄 Push & Pull

```bash
# Push new branch
git push -u origin feature/123      # -u sets upstream
git push origin feature/123         # After upstream is set

# Push with force (careful!)
git push origin feature/123 --force-with-lease  # Safe force push
git push origin feature/123 --force             # DANGEROUS - don't use

# Pull
git pull origin develop             # Pull from remote
git pull --rebase                   # Rebase instead of merge

# Fetch only (no merge)
git fetch origin                    # Fetch all
git fetch upstream                  # Fetch from upstream (forks)
```

---

## 🔀 Merge & Rebase

```bash
# Merge branches
git co develop
git merge feature/123               # Merge feature into current
git merge --squash feature/123      # Squash commits before merge
git merge --no-ff feature/123       # Create merge commit

# Rebase (replay commits)
git rebase develop                  # Rebase current onto develop
git rebase -i HEAD~3                # Interactive rebase (last 3 commits)

# During conflict
git rebase --continue               # Continue after fixing
git rebase --abort                  # Cancel rebase

# Abort all operations
git merge --abort                   # Cancel merge
git rebase --abort                  # Cancel rebase
```

---

## 🔍 Inspect & Diff

```bash
# Status
git st                              # Short status
git status                          # Long status

# Diff
git diff                            # Unstaged vs staged
git diff --cached                   # Staged vs committed
git diff HEAD                       # Unstaged + staged vs HEAD
git diff develop feature/branch     # Compare branches
git diff develop -- file.go         # Specific file

# Show commit details
git show <commit-hash>              # Show commit and diff
git show <commit-hash>:file.go      # Show file at commit

# Who changed what
git blame file.go                   # Show blame info
git blame -L 10,20 file.go          # Specific lines
```

---

## 🔧 Undo & Recovery

```bash
# Undo uncommitted changes
git restore file.go                 # Discard changes (git 2.23+)
git checkout -- file.go             # Older git versions
git clean -fd                       # Remove untracked files

# Undo commits
git undo                            # Soft reset last commit (requires alias)
git reset --soft HEAD~1             # Same as above
git reset --mixed HEAD~1            # Undo + unstage
git reset --hard HEAD~1             # Undo + discard changes (DANGEROUS)

# Revert (safe undo for public branches)
git revert <commit-hash>            # Create new commit that undoes changes

# Stash temporarily
git stash                           # Save changes
git stash list                      # List stashes
git stash pop                       # Restore and delete
git stash apply                     # Restore (keep stash)
git stash drop                      # Delete stash
```

---

## 📌 Tags

```bash
# Create tags
git tag v0.2.0                      # Lightweight tag
git tag -a v0.2.0 -m "Release"     # Annotated tag

# List tags
git tag                             # All tags
git tag -l "v0.*"                   # Pattern match
git show v0.2.0                     # Show tag details

# Push tags
git push origin v0.2.0              # Push specific tag
git push origin --tags              # Push all tags

# Delete tags
git tag -d v0.2.0                   # Delete local
git push origin --delete v0.2.0     # Delete remote
```

---

## 👥 Collaboration (Forks)

```bash
# Setup fork
git remote add upstream https://github.com/codenoidlab/vpsctl.git
git remote -v                       # Verify remotes

# Keep fork in sync
git fetch upstream                  # Fetch latest from main repo
git co develop
git rebase upstream/develop         # Rebase on upstream
git push origin develop             # Push to your fork

# Create PR
git push origin feature/123         # Push feature
# Then create PR on GitHub

# Sync after PR merge
git checkout develop
git pull upstream develop           # Pull main repo develop
git push origin develop             # Update your fork
```

---

## 🚨 Emergency Commands

```bash
# Force push (ONLY if you're sure)
git push origin feature/branch --force-with-lease

# Cherry-pick a commit
git cherry-pick <commit-hash>       # Apply single commit

# Find a commit
git log --grep="search text" --oneline

# Find lost commits
git reflog                          # Show all ref changes
git reflog | head                   # Recent changes

# Create patch
git format-patch origin/develop -o patches/

# Apply patch
git apply patches/0001-fix.patch
```

---

## ✅ Pre-Push Checklist

```bash
# Before pushing any branch:
git st                              # Check status
go fmt ./...                        # Format code
go test ./...                       # Run tests
golangci-lint run                   # Lint
gosec ./...                         # Security check
git logs | head -5                  # Review commits
git diff origin/develop             # Review changes

# Then push
git push origin feature/branch
```

---

## 📱 Common Workflows

### Starting a Feature
```bash
git co develop && git pull upstream develop
git co -b feature/456-awesome
# ... make changes ...
git add . && git ci -m "feat: awesome feature"
git push -u origin feature/456-awesome
# Create PR on GitHub
```

### Syncing Fork with Upstream
```bash
git fetch upstream
git co develop
git rebase upstream/develop
git push origin develop
```

### Fixing Last Commit
```bash
# Forgot to add a file?
git add forgotten-file.go
git amend

# Wrong message?
git ci --amend -m "correct message"

# Push the fix
git push origin feature/branch --force-with-lease
```

### Squashing Commits
```bash
git rebase -i HEAD~3                # Last 3 commits
# Mark: keep (pick), squash (s), reword (r)
git push origin feature/branch --force-with-lease
```

### Creating a Release
```bash
git co develop && git pull upstream develop
git co -b release/v0.3.0
# Update version numbers
git add . && git ci -m "chore: v0.3.0"
git push -u origin release/v0.3.0
# Create PR to main
# After merge:
git co main && git pull upstream main
git tag -a v0.3.0 -m "Release v0.3.0"
git push upstream main v0.3.0
git co develop && git merge main && git push upstream develop
```

---

## 🆘 Need Help?

```bash
# Git help
git help <command>                  # e.g., git help push
man git                             # Full manual

# See what will happen (dry run)
git push --dry-run                  # Simulate push
git merge --no-commit --no-ff       # Test merge

# Undo almost anything
git reflog                          # Find lost commits/changes
```

---

## 🎓 Learning Resources

- [Pro Git Book](https://git-scm.com/book)
- [GitHub Help](https://docs.github.com)
- [VPSctl GIT_WORKFLOW.md](GIT_WORKFLOW.md)
- [Conventional Commits](https://www.conventionalcommits.org/)

---

## 📝 Cheat Sheet Format Quick Reference

```
git <verb> <noun> [options]
git <verb> <target> <action>
```

**Common patterns:**
- Branch: `git branch <name>`, `git checkout <branch>`
- Remote: `git push <remote> <branch>`, `git pull <remote> <branch>`
- Commit: `git commit -m "message"`, `git add <files>`
- History: `git log`, `git diff`, `git show`
- Undo: `git reset`, `git revert`, `git restore`

---

**Keep this handy! 📌**
