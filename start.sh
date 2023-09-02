#!/bin/sh
npx tailwindcss -i ./tailwind.css -o ./assets/css/main.css
go run main.go
