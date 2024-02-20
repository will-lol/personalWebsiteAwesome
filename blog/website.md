---
title: Building will.forsale
date: 2024-02-20
slug: website
description: >
  6 months ago I started work on my personal website. Now it's completed, I reflect on this process and my learnings.
---
A personal website is a common project for new and experienced programmers alike. The blank slate that is an HTML file leads people to an infinite number of different websites of varying complexity. Having created two of my own before, I knew this website would lose its shine quickly—but that didn't stop me from attempting to build on what I thought would be the ultimate tech stack.

From watching many hours of [*The Primeagen*](https://www.youtube.com/c/theprimeagen), I knew that JavaScript was the devil and Rust is king. And because of my desire to leave the JavaScript ecosystem, I decided to build using Go and the much revered HTMX. Coming from the land of Vercel, TypeScript, and JavaScript frameworks, this was a transition sure to teach me a lot. I laced up my boots, init'ed a Git repo, and made a very hopeful Google search: "golang web framework".

Of course, I was yet to learn that the sanctum of JavaScript does things quite differently to other programming language. In Go, a web framework is more of an "http router" and less of a towering monolith of minified code and build steps that somehow results in an actual hosted URL on the internet. Regardless, I chose the one I saw first ([chi](https://go-chi.io/)), and joined the hoard of other confused JavaScript and Python developers in (r/golang)[https://old.reddit.com/r/golang/comments/ysy2jn/i_just_built_a_new_fullstack_framework_for_golang/] that desperately want to build their website in Go.
