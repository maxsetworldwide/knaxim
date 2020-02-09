# Testing the Knaxim Client.
Instructions for setting up a karma instance for testing a Knaxim client.


The requirements are a Knaxim web server connected to a Knaxim client along
with application code, tests, and karma dependencies.  Node and npm are also
required.

The Kaxim server and client are stored in repositories. All relevant code and
tests are located in the Knaxim client repository.  All karma dependencies
are located in the package.json of the Knaxim client repository.  Any 
environment specific settings should be placed in the /local folder of the
Knaxim client repo.


## Installing Karma
Install node.js if it is not installed already. nvm is recommended for node
installation and management.

[Installing nvm](https://github.com/nvm-sh/nvm#install--update-script)
[Installing node with nvm](https://github.com/nvm-sh/nvm#usage)
[NodeJS](https://nodejs.org/en/)

### Install node
nvm install node

### Install karma-cli globally.
npm install -g karma-cli

### Install node dependencies, run with repository root as working directory
npm install

## Test Server Setup
1. Launch a knaximserver in test mode (no ssl)
	options for config karma for knaximserver proxy
	1. custom karma config in local folder
	2. ssh port map knaximserver to localhost:8000
2. Start Karma with `karma start`, or `karma start local/custom.config` if
using a custom config.
3. Connect one or more browsers.  After running Karma it should give you a
url for connecting to the karma server.  Load that url in any browser you
want to test from.

## Run tests.
Tests will automatically run every time a tracked file is modified.  Tracked
files include any tests and code paths listed in the `files` section of the
karma config (karma.conf.js).

## Write Tests.
All tests are located in the project /test folder, and all test files
must end with `spec.js`.
