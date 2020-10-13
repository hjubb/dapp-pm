# dapp-pm
A tool to manage dependencies from a github repo if you're using dapp tools

![Build](https://github.com/hjubb/dapp-pm/workflows/Build/badge.svg)

### Motivation
I wasn't satisfied with using the dapp tools to import dependencies because it relied on git submodules and from what I could tell it had no way of specifying a version of the dependency you wanted as well as enforcing a specific folder structure it expected to see in the repo. Instead of copy-pasting individual solidity files and hoping I didn't make any human errors with versioning or resorting to re-writing interfaces from scratch if I couldn't be bothered to copy all the dependencies to satisfy a compilation I decided to build a tool to do it for me, specifying the repo and the version you want by tag (to make it reproducible or something). It also caches the repos you get so you can build offline in new projects

### Usage
#### Regular usage
Running `dapp-pm [github-username]/[github-repo]@[version-tag]` will download the repo to `~/.dapp-pm/[path]` and present you with a terminal UI like the following

![Terminal UI](https://imgur.com/LuJYdGK.png)

After selecting the contracts you want to import, it will copy the .sol files and their dependencies to `./src/lib/[path]/[name].sol`, example output of import `ERC20.sol` from `openzeppelin/openzeppelin-contracts@v3.2.0`

![Tree ./src](https://i.imgur.com/PhCd8aG.png)

#### Power users
Running `dapp-pm [github-username]/[github-repo]@[version-tag] [solidity file with extension]` will result in a shortcut skipping the UI, directly importing the file and its dependencies if there exists only one file with that name.
