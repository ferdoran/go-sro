# go-sro
A custom server/backend implementation of the game Silkroad Online
written in Go.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ferdoran/go-sro)
![Lines of code](https://img.shields.io/tokei/lines/github/ferdoran/go-sro)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/ferdoran/go-sro/Build%20and%20Publish%20Gateway%20Server%20Image?label=Gateway%20Server)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/ferdoran/go-sro/Build%20and%20Publish%20Gateway%20Server%20Image?label=Agent%20Server)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/ferdoran/go-sro/Build%20and%20Publish%20Agent%20Server%20Image?label=DB)

## Architecture

![architecture diagram](http://plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ferdoran/go-sro/main/_docs/architecture.puml)

The backend architecture consists of several components that may interact with each other.
There are the 3 front-/client-facing servers (_Gateway Server_, _Agent Server_ and _Download / Patch Server_)
and 3 backend servers (multiple _Game Servers_, _Shard Server_, _Chat Server_).

All of them handle different kind of aspects to the game:

- **Download/Patch Server**: Provide updates and patches to the clients.
- **Gateway Server**: Perform authentication and transfer to the specific realm or shard.
- **Agent Server**: Proxy server for the client through which all network traffic is sent.
Takes care routing network traffic to the correct servers.
- **Game Server**: Handles core game logic (navigation, AI, combat, ...)
and game objects (players, pets, NPCs, ...).
Usually there are multiple game servers, each handling a different region of the overall map
- **Shard Server**: Handles all region-independent logic (guild, party, market, events, ...)
- **Chat Server**: Handles overall chat messages (except local/region chat).
Could also be handled by **Shard Server**
