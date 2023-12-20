#!/bin/sh
npx tailwindcss -i ./tailwind.css -o ./assets/css/main.css
templ generate
go run main.go
