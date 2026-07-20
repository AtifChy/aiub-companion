#!/bin/sh

# Get list of staged files
STAGED_FILES=$(git diff --cached --name-only)

# If no files are staged, exit early
if [ -z "$STAGED_FILES" ]; then
    exit 0
fi

# Flag to track failure status
FAIL=0

# Filter staged files
GO_FILES=$(echo "$STAGED_FILES" | grep -E '\.go$' || true)
FRONTEND_FILES=$(echo "$STAGED_FILES" | grep -E '^frontend/src/' || true)

# 1. Run Go checks if Go files are staged
if [ -n "$GO_FILES" ]; then
    echo "----------------------------------------"
    echo "🔍 Checking Go files..."
    echo "----------------------------------------"

    # Format check (checks if any staged files need formatting)
    echo "Running 'go fmt' check..."
    # We pass the staged files to gofmt -l to see if any need formatting
    UNFORMATTED_GO=$(gofmt -l $GO_FILES || true)
    if [ -n "$UNFORMATTED_GO" ]; then
        echo "❌ The following Go files are not formatted correctly:"
        echo "$UNFORMATTED_GO"
        echo "Please run 'go fmt ./...' to format them before committing."
        FAIL=1
    else
        echo "✅ Go files are formatted correctly."
    fi

    # Run go vet
    echo "Running 'go vet'..."
    if ! go vet ./...; then
        echo "❌ 'go vet' failed. Fix errors and try again."
        FAIL=1
    else
        echo "✅ 'go vet' passed."
    fi
fi

# 2. Run Frontend checks if Frontend files are staged
if [ -n "$FRONTEND_FILES" ]; then
    echo "----------------------------------------"
    echo "🔍 Checking Frontend files..."
    echo "----------------------------------------"

    # Check if pnpm is installed
    if ! command -v pnpm >/dev/null 2>&1; then
        echo "⚠️  pnpm is not installed. Skipping frontend checks."
    else
        # Run oxfmt check
        echo "Running frontend formatter check (oxfmt)..."
        if ! pnpm --dir frontend run fmt:check; then
            echo "❌ Frontend formatting check failed. Run 'pnpm --dir frontend run fmt' to format."
            FAIL=1
        else
            echo "✅ Frontend formatting check passed."
        fi

        # Run oxlint check
        echo "Running frontend linter (oxlint)..."
        if ! pnpm --dir frontend run lint; then
            echo "❌ Frontend linting failed. Fix lint errors and try again."
            FAIL=1
        else
            echo "✅ Frontend linting passed."
        fi

        # Run Typecheck
        # echo "Running frontend typecheck..."
        # if ! pnpm --dir frontend run typecheck; then
        #     echo "❌ Frontend typecheck failed. Fix typescript errors and try again."
        #     FAIL=1
        # else
        #     echo "✅ Frontend typecheck passed."
        # fi
    fi
fi

echo "----------------------------------------"
if [ $FAIL -ne 0 ]; then
    echo "❌ Git pre-commit hook failed. Please resolve the issues above."
    exit 1
fi

echo "✅ Git pre-commit hook passed successfully!"
exit 0
