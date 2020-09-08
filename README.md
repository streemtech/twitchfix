[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

# Twitchfix

Twitch.tv has chosen to only [partially comply](https://dev.twitch.tv/docs/authentication/getting-tokens-oidc#oidc-authorization-code-flow) with the OIDC/Oauth RFC. This library's goal is to act as a mutation proxy that takes twitch's response and converts them to library compatible responses.