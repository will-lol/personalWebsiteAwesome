find . -not \( -path ./.git -prune \) -not \( -path ./cdk -prune \) -not \( -path ./.go -prune \) -not \( -name "*_templ.go" \) | entr -r sh start.sh
