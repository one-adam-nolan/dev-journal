# Dev Journal (DJ)

## Overview

It started with tracking my work in markdown files a few years back. I'd create a new markdown file and start typing away each day. The formatting could have been more consistent, timestamps were only sometimes there, or  I would get too lazy to type it out. 

After a while, I made a few Python scripts to spice things up. But I wrote those with relative paths, meaning I had to execute them in a specific directory. Also, my "source code" was in the journal directory itself. 

While attempting to add a Vim script to execute the Python scripts from any terminal session anywhere, I realized I needed to create a CLI application to solve this problem.

This project solves a few problems:
* I get work with `Go`, which is new to me
* Log my work from anywhere in a terminal, including with Vim scripts
* Work with Dev Containers

## Installing

1. `go build -o dj`
1. `sudo mv dj /usr/local/bin`

> If you use `go install` it will put the binary in the $GOBIN, which may make it so the autocompletion does not work properly.