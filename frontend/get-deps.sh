#! /bin/bash -e +x
# https://gist.github.com/developit/f1a3425ded9b52206c372b6c6b54ddcb
npm i
# scaffold your project
# npx preact-cli@rc create default my-app && cd my-app
# upgrade to Preact X
npm i -D preact-cli@rc && npm i preact@latest preact-router@latest preact-render-to-string@latest
npm rm preact-compat
npm uninstall -S preact-compat
npm i preact-render-to-string@">=5.0.2"
