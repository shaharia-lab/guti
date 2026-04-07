#!/bin/bash
(
  echo "=== Changes since v1.5.0 ===" && \
  echo -e "\n=== Commit Messages ===" && \
  git log v1.5.0..HEAD --pretty=format:"* %s%n  %b" --no-merges && \
  echo -e "\n=== File Changes ===" && \
  git diff --stat v1.5.0..HEAD && \
  echo -e "\n=== Detailed Changes ===" && \
  git log v1.5.0..HEAD --patch --no-merges -- '*.go' && \
  echo -e "\n=== New Files ===" && \
  git diff --name-only --diff-filter=A v1.5.0..HEAD && \
  echo -e "\n=== Deleted Files ===" && \
  git diff --name-only --diff-filter=D v1.5.0..HEAD && \
  echo -e "\n=== Dependencies Changes ===" && \
  git diff v1.5.0..HEAD go.mod go.sum
) > changes.txt