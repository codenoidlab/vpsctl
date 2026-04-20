#!/bin/bash
# VPSctl GitHub Push Commands
# Copy and run these commands to push your repository to GitHub

# Make sure you've created an empty repository at:
# https://github.com/codenoidlab/vpsctl

echo "🚀 Pushing VPSctl to GitHub"
echo "  Repository: github.com/codenoidlab/vpsctl"
echo ""

# Add remote
echo "Step 1: Add remote"
git remote add origin https://github.com/codenoidlab/vpsctl.git

# Verify remote
echo "Step 2: Verify remote"
git remote -v

# Push main branch
echo "Step 3: Push main branch"
git push -u origin main

# Push develop branch
echo "Step 4: Push develop branch"
git push -u origin develop

# Push tags
echo "Step 5: Push tags (v0.1.0)"
git push origin --tags

echo ""
echo "✅ Repository pushed to GitHub!"
echo ""
echo "Next steps:"
echo "1. Visit: https://github.com/codenoidlab/vpsctl"
echo "2. Go to Settings → Branches"
echo "3. Set 'develop' as the default branch"
echo "4. Add branch protection rules:"
echo "   - Protect main branch"
echo "   - Protect develop branch"
echo ""
echo "5. Your repository is now ready for contributions!"
