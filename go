#!/bin/bash
set -e

goal_serve() {
  bundle exec jekyll serve --drafts
}

goal_updateDependencies() {
  bundle update
}

if type -t "goal_$1" &>/dev/null; then
  goal_$1 ${@:2}
else
  echo "usage: $0 <goal>

goal:
    serve              -- start a development server for preview (including drafts)
    updateDependencies -- update github pages and jekyll dependencies
"
  exit 1
fi
