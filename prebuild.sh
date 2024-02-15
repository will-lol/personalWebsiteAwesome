#!/bin/sh
tailwindcss -i ./lib/tailwind.css -o ./assets/css/main.css & templ generate & fd -E tailwind.config.js -E assets -e js -x rsync -Ra {} assets/js/
wait
