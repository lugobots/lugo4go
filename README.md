# MakeItPlay - Football Go Player Client

[![GoDoc](https://godoc.org/github.com/makeitplay/client-player-go?status.svg)](https://godoc.org/github.com/makeitplay/client-player-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/makeitplay/client-player-go)](https://goreportcard.com/report/github.com/makeitplay/client-player-go)

Go Player Client is a [Go](http://golang.org/) implementation of a client player for [MakeItPlay football](http://www.makeitplay.ai/football) game server. 

It **is not a bot** that plays the game, it is only the client for the game server. 

This client implements a brainless player in the game. So, this library implements many methods that does not affect the player
intelligence/behaviour/decisions. It is meant to reduce the developer concerns on communication, protocols, attributes, etc, and 
focus in the player intelligence.

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
0. Now you will need to start your team process. Each team must have 11 process (one for each player).
    
    **Option A**: You may start your team players manually executing the command `./myAwesomeBot -team=home -number=[1-11]`
    eleven times. 
          
    **or**
    
    **Option B**: You can use the script in [the example directory](./example) to do this automatically for you:
    `./play.sh home`
0. And finally you may do the same for the other team. 
    
    **Option A**: You play against your own team repeating the last step, but in the `away` side: `./play.sh away`
    
    **or**
        
    **Option B**: You may play against [**The Dummies** team](https://github.com/makeitplay/the-dummies-go) executing the script `start-team-container.sh`
    available in [the example directory](./example):     
    
    `./start-team-container.sh makeitplay/the-dummies-go away`
     
    **or**
    
    **Option C**: You may play against another team:
    `./start-team-container.sh [container image name] away` 

### Next steps

As you may have noticed, the bot player in the example directory does not play well. 
Now, you may start your bot implementing its behaviour based on the game state after each message got by the function 
`reactToNewState`.  

### Deploying you bots

After developing your bot, you may share it with other developers.

Using this client your code will be compiled into a binary file. So you do not have to share the bot source code.

There is a Dockerfile template in [the example directory](./example) to guide you how to create a container. After
having customized (or not) your Dockerfile, you just need to build the container:

```bash
docker build -t [your username]/[your bot awesome name] .
docker push
```

If you are not familiar with Dockerfile, Docker composer, and so on, consider spending 11 minutes to learn it 
watching [this video](https://www.youtube.com/watch?v=YFl2mCHdv24). It is by far the best and simplest way to learn Docker. 
