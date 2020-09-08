[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

# Twitchfix

Twitch.tv has chosen to only [partially comply](https://dev.twitch.tv/docs/authentication/getting-tokens-oidc#oidc-authorization-code-flow) with the OIDC/Oauth RFC. This library's goal is to act as a mutation proxy that takes twitch's response and converts them to library compatible responses.

This is **not** a library for use in other programs, only a program to act as an in-between for other OIDC programs.

## Standing Up The Service

### Docker Container
The docker container is hosted on dockerhub, and can be run with the standard dockerish ways.
```
docker run --name twitchfix -p 8284:8284 streemtech/twitchfix
```
The makefile also has commands to build the docker container from scratch, but the container is literally from the scratch container, so there is little difference from building standalone. Adding a second dockerfile using `from: golang` is one of the goals for the future.

### Standalone
The commands to build and run the standalone executable are all in the makefile. The command `make standalone` will compile a standalone binary called `twitchfix`, which should be able to be executed from the command line. The program can be built with just `go build` if desired.

## Usage
Twitchfix currently only handles the second half of the oidc authorization code flow, the POST commands to https://id.twitch.tv/oauth2/token, which are steps three and four of the documentation.

To Use twitchfix in your library, instead of using the link provided by twitch, use an http request to the `/oauth2/token` endpoint with ip and port of twitchfix. For example, `http://192.168.100.123:8284/oauth2/token`. The client ID and Secret are not required by twitchfix, and should only be required by your third party library. 

While there should not be any issue with doing so, twitchfix should NOT be exposed to the internet, and ideally should be as close as possible to the library/service making calls to twitchfix.

## Future Goals
* Implement HTTPS
	* Loading of Certs
	* Creation of certs if they don't exist
* Implement config file
	* Endpoint
	* Http(s)
	* Port
	* Log Level
* Logging (currently use labstack/echo default logging)
* Dockerfile that compiles internally
