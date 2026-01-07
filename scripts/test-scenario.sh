#!/bin/bash

set -e

echo "🧪 Running integration tests..."

# Clean start
rm -rf ~/.twitter-cli/
./twt logout 2>/dev/null || true

echo "✅ Step 1: Create users"
./twt user create alice
./twt user create bob
./twt user create carol

echo "✅ Step 2: Login and post"
./twt login alice
./twt post "Alice's first post"
./twt post "Alice loves Go"

./twt login bob
./twt post "Bob's wisdom"

./twt login carol
./twt post "Carol's insight"

echo "✅ Step 3: Follow users"
./twt login alice
./twt follow bob
./twt follow carol

echo "✅ Step 4: Check feed"
./twt feed

echo "✅ Step 5: Engagement"
./twt profile bob
BOB_POST=$(./twt profile bob | head -1 | awk '{print $1}')
./twt like $BOB_POST
./twt retweet $BOB_POST

echo "✅ Step 6: View stats"
./twt show $BOB_POST

echo "✅ Step 7: Test following/followers"
./twt following
./twt followers

echo "✅ Step 8: Pagination"
./twt feed --limit 2

echo ""
echo "🎉 All integration tests passed!"