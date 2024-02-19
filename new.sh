#!/bin/sh
cat > blog/"$1.md" <<- EOM
---
title: Untitled
date: $(date --iso-8601)
slug: ${1}
description: >
	  Add a description here!
---
EOM
$EDITOR blog/"$1.md"

