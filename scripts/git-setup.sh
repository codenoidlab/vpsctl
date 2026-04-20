#!/bin/bash
# VPSctl Git Setup Script
# Run this after cloning to configure your git environment
# Usage: bash scripts/git-setup.sh

set -e

echo "🔧 VPSctl Git Environment Setup"
echo "================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo "❌ Git is not installed. Please install Git first."
    exit 1
fi

echo -e "${BLUE}Step 1: Configure Git User${NC}"
echo "-----------------------------"

# Check if user.name is configured
if git config user.name > /dev/null 2>&1; then
    CURRENT_NAME=$(git config user.name)
    echo "Current user.name: $CURRENT_NAME"
    read -p "Use this name? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter your name: " USER_NAME
        git config --global user.name "$USER_NAME"
    fi
else
    read -p "Enter your name: " USER_NAME
    git config --global user.name "$USER_NAME"
fi

# Check if user.email is configured
if git config user.email > /dev/null 2>&1; then
    CURRENT_EMAIL=$(git config user.email)
    echo "Current user.email: $CURRENT_EMAIL"
    read -p "Use this email? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter your email: " USER_EMAIL
        git config --global user.email "$USER_EMAIL"
    fi
else
    read -p "Enter your email: " USER_EMAIL
    git config --global user.email "$USER_EMAIL"
fi

echo -e "${GREEN}✓ User configured${NC}"
echo ""

echo -e "${BLUE}Step 2: Configure Git Aliases${NC}"
echo "-------------------------------"

# Setup helpful aliases
git config --global --replace-all alias.st status
git config --global --replace-all alias.co checkout
git config --global --replace-all alias.br branch
git config --global --replace-all alias.ci commit
git config --global --replace-all alias.amend "commit --amend --no-edit"
git config --global --replace-all alias.undo "reset --soft HEAD~1"
git config --global --replace-all alias.logs "log --oneline --graph --all"
git config --global --replace-all alias.sync "!git fetch origin && git rebase origin/develop"

echo -e "${GREEN}✓ Aliases configured:${NC}"
echo "  git st               → git status"
echo "  git co <branch>      → git checkout"
echo "  git br               → git branch"
echo "  git ci -m 'message'  → git commit"
echo "  git amend            → git commit --amend"
echo "  git undo             → git reset --soft HEAD~1"
echo "  git logs             → git log (formatted)"
echo "  git sync             → git fetch + rebase"
echo ""

echo -e "${BLUE}Step 3: Configure Git Behavior${NC}"
echo "--------------------------------"

# Set pull behavior
git config --global --replace-all pull.rebase false
echo -e "${GREEN}✓ git pull${NC} will use merge (not rebase)"

# Set push behavior
git config --global --replace-all push.default current
echo -e "${GREEN}✓ git push${NC} will use current branch name"

# Auto-follow tags
git config --global --replace-all push.followTags true
echo -e "${GREEN}✓ Tags will be auto-pushed${NC}"

echo ""

echo -e "${BLUE}Step 4: Configure Merge Conflict Style${NC}"
echo "--------------------------------------"

git config --global --replace-all merge.conflictstyle zdiff3
echo -e "${GREEN}✓ Merge conflicts will use zdiff3 style${NC}"

echo ""

echo -e "${BLUE}Step 5: Setup Project Remote (if in vpsctl directory)${NC}"
echo "-------------------------------------------------------"

if [ -f ".git/config" ]; then
    # Check if we're in vpsctl repo
    REMOTE_URL=$(git config --get remote.origin.url || echo "")
    
    if [[ $REMOTE_URL == *"vpsctl"* ]]; then
        # Check if upstream exists
        if ! git remote get-url upstream > /dev/null 2>&1; then
            echo "Setting up upstream remote..."
            
            # Determine if this is a fork
            if [[ $REMOTE_URL == *"codenoidlab"* ]]; then
                echo "This is the main repository (no upstream needed)"
            else
                echo "This appears to be a fork."
                git remote add upstream https://github.com/codenoidlab/vpsctl.git
                echo -e "${GREEN}✓ Upstream remote added${NC}"
            fi
        else
            echo "Upstream already configured: $(git remote get-url upstream)"
        fi
    fi
else
    echo -e "${YELLOW}⚠ Not in a git repository root${NC}"
fi

echo ""

echo -e "${BLUE}Step 6: Verify Configuration${NC}"
echo "-----------------------------"

echo -e "${GREEN}Global Git Config:${NC}"
git config --global --list | grep -E "user\.|alias\.|push\.|pull\.|merge\." | head -15

echo ""
echo -e "${GREEN}Project Git Remotes:${NC}"
if [ -f ".git/config" ]; then
    git remote -v
else
    echo "  (not in a git repository)"
fi

echo ""
echo -e "${GREEN}✅ Setup Complete!${NC}"
echo ""
echo "Next steps:"
echo "  1. Verify remote URLs:"
echo "     git remote -v"
echo "  2. Try the aliases:"
echo "     git st          # Should show status"
echo "  3. Create a feature branch:"
echo "     git co -b feature/my-feature"
echo ""
echo "For more info, see: GIT_WORKFLOW.md"
