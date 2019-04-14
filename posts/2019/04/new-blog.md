title: Updated Site Build
date: 2019-04-14
excerpt: 
image: 
---
<small>2019-04-14</small>
# Updated Site Build

I've recently undertaken rebuilding my site from scratch. Previously I used the opportunity to
gain some experience building a simple blog using Python and Django. Working with those tools
was a good experience in general, but didn't quite satisfy the way I wanted the thing to operate.

So, I took the opportunity to learn once again. This time I've built the site using Go. I relied
almost solely on native libraries, only pulling on one external library to handle Markdown parsing
for any individual posts.

Overall it's pretty simple and probably sloppy. I'll work on cleaning up the code and probably make
it public again soon, *sans* post content. As the deployment scripts are pretty centred around my own
server setup, it would probably benefit from a modification to open that up to variable paths.

The end plan is to keep the codebase simple and malleable, while accepting the minimal amount of responsibility
for catering to others' deployment and security needs. I just want to provide a single Go file with some helpful
project structure and possibly a Makefile to get others started.

To be continued.
