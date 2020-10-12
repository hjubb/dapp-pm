# dapp-pm
A tool to manage dependencies from a github repo if you're using dapp tools

### Motivation
I wasn't satisfied with using the dapp tools to import dependencies because it relied on git submodules and from what I could tell it had no way of specifying a version of the dependency you wanted as well as enforcing a specific folder structure it expected to see in the repo. Instead of copy-pasting individual solidity files and hoping I didn't make any human errors with versioning or resorting to re-writing interfaces from scratch if I couldn't be bothered to copy all the dependencies to satisfy a compilation I decided to build a tool to do it for me, specifying the repo and the version you want by tag (to make it reproducible or something). It also caches the repos you get so you can build offline in new projects

### Usage
`dapp-pm [github-username]/[github-repo]@[version-tag]`
