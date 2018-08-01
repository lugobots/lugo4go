# MakeItPlay - Football Go Player Client

[![GoDoc](https://godoc.org/github.com/makeitplay/client-player-go?status.svg)](https://godoc.org/github.com/makeitplay/client-player-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/makeitplay/client-player-go)](https://goreportcard.com/report/github.com/makeitplay/client-player-go)

Go Player Client is a [Go](http://golang.org/) implementation of a client player for [MakeItPlay football](http://www.makeitplay.ai/football) game server. 

It **is not a bot** that plays the game, it is only the client for the game server. 

This client implements a brainless player in the game. So, this library implements many methods that does not affect the player
intelligence/behaviour/decisions. It was meant to reduce the developer concerns on communication, protocols, attributes, etc,
and focusing in the player intelligence.

Using this client, you just need to implement the Artificial Intelligence of your player and some other few methods to support
your strategy (see the project [The Dummies](https://github.com/makeitplay/the-dummies-go) as an example). 
 
### Documentation

* [API Reference](http://godoc.org/github.com/makeitplay/client-player-go)

### Requirements

0. Docker >= 18.03 (https://docs.docker.com/install/)
0. Docker Compose >= 1.21 (https://docs.docker.com/compose/install/)
0. Go Lang >= 1.10 (https://golang.org/doc/install)

### Installation

    go get github.com/makeitplay/client-player-go

### Kick start

0. Copy [the example directory](./example) as a new Golang project
0. Build your bot executing the command below inside the project directory
    ```bash 
    go build -o myAwesomeBot
    ```
0. Run the game server using the command 
    ```bash
    docker run -p 8080:8080  makeitplay/football:1.0.0-alpha
    ```
0. You will need to spin up 22 process (11 as the home team players, and 11 as the away team players). 

    You can do this manually executing the command `./myAwesomeBot -team=[home|away] -number=[1-11]`
      
    **or**
    
    You can use the shell script in the example directory to do this automatically for you.

**Note:** In this example above, both teams will have the same bots. You may download another bot
to play against your team (soon available at [MakeItPlay Docker Hub](https://hub.docker.com/r/makeitplay/))  

### Deploying you bots (soon)

You will be able to create a Docker container with your player bot and share it with another developers.

See the Dockerfile template (not tested yet) in the example directory.

