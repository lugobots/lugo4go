# Lugo4Go - A Lugo Bots Client
[![Build Status](https://travis-ci.org/lugobots/lugo4go.svg?branch=master)](https://travis-ci.org/lugobots/lugo4go)
[![GoDoc](https://godoc.org/github.com/lugobots/lugo4go?status.svg)](https://godoc.org/github.com/lugobots/lugo4go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lugobots/lugo4go)](https://goreportcard.com/report/github.com/lugobots/lugo4go)

Lugo4Go is a [Go](http://golang.org/) implementation of a client player for [Lugo](https://lugobots.dev/) game. 

It **is not a bot** that plays the game, it is only the client for the game server. 

This client implements a brainless player in the game. So, this library implements many methods that does not affect the player
intelligence/behaviour/decisions. It is meant to reduce the developer concerns on communication, protocols, attributes, etc.

Using this client, you just need to implement the Artificial Intelligence of your player and some other few methods to support
your strategy (see the project [The Dummies](https://github.com/lugobots/the-dummies-go) as an example). 
 
### Documentation

(usage examples below)

* [API Reference](http://godoc.org/github.com/lugobots/lugo4go)

### Requirements

* Go Lang >= 1.16 (https://golang.org/doc/install)

### Installation

    git get https://github.com/lugobots/lugo4go.git

### Usage

There are two ways to use **Lugo4Go** client:

### Implementing a Bot Interface (simpler and recommended)

See [example](./examples/bot-interface)

**Lugo4Go** *PlayAsBot* implements a very basic logic to reduce the code boilerplate. This client will wrap most repetitive
code that handles the raw data got by the bot and will identify the player state.

The `Bot` interface requires the methods to handle each player state based on the ball possession.

```go

type Bot struct {
	
}

// OnDisputing is called when no one has the ball possession
func (b *Bot) OnDisputing(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	// the magic code comes here
	return ...
}

// OnDefending is called when an opponent player has the ball possession
func (b *Bot) OnDefending(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	// the magic code comes here
	return ...
}

// OnHolding is called when this bot has the ball possession
func (b *Bot) OnHolding(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	// the magic code comes here
	return ...
}

// OnSupporting is called when a teammate player has the ball possession
func (b *Bot) OnSupporting(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	// the magic code comes here
	return ...
}

// AsGoalkeeper is only called when this bot is the goalkeeper (number 1). This method is called on every turn,
// and the player state is passed at the last parameter.
func (b *Bot) AsGoalkeeper(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot, state lugo4go.PlayerState) error {
	// the magic code comes here
	return ...
}
```

### Implementing the turn handler (a little more work)

See [example](./examples/turn-handler)

If you rather to handle everything on your side, you only need to implement the `TurnHandler` interface.

The `TurnHandler` will receive the turn context and the turn snapshot for each turn.

Your bot will need a `lugo4go.OrderSender` to send the orders back to the Game Server during each turn.

```go

type Bot struct {
    OrderSender lugo4go.OrderSender
}

func (b *Bot) Handle(ctx context.Context, snapshot *proto.GameSnapshot) {
	// the magic code comes here
	resp, err := b.OrderSender.Send(ctx, snapshot.Turn, orders, "")
}

```

## Kick-start

0. Copy one of the examples from [the example directory](./examples) as a new Golang project.

1. Run the game server using the command 
    ```bash
    docker run -p 8080:8080  lugobots/server:v1.1 play --dev-mode
   ```
2. Now you will need to start your team processes. Each team must have 11 process (one for each player).
    
    **Option A**: You may start your team players manually executing the command `go run main.go -team=home -number=[1-11]`
    eleven times. 
          
    **or**
    
    **Option B**: You can use the script in examples directory to do this automatically for you:
    `./play.sh home`

3. And, **if your have not started the away team**, you may do the same for the other team. 
    
    You play against your own team repeating the last step, but in the `away` side: 
    ```
    ./play.sh away
   ```

## Field Library

There are a many things that you will repeatedly need to do on your bot code, e.g. getting your bot position,
creating a move/kick/catch order, finding your teammates positions, etc. 

The Field library brings a collection of functions that will help you get that data from the game snapshot:

Examples ([see all functions on the documentation page](https://pkg.go.dev/github.com/lugobots/lugo4go))

```go

myTeamGoal := field.GetTeamsGoal(proto.Team_HOME)

moveOrder, err := field.MakeOrderMoveMaxSpeed(*me.Position, myTeamGoal)

```

### Mapper and Region

The Field library also provides a quite useful pair: the Mapper interface and Region struct.

#### The Mapper

`Mapper` slices the field in columns and rows, so your bot does not have to care about precise coordinates or the team
side. The mapper will automatically translate the map position to the bot side.

And you may define how many columns/rows your field will be divided into.

```go

// let's create a map 10x5 
fieldMapper, err := field.NewMapper(10, 5, proto.Team_HOME)

aRegion, err := fieldMapper.GetRegion(2, 4)

```

#### The Region

The `Mapper` will slice the field into `Region`s. The Region struct helps your bot to move over the field without caring
about coordinates or team side.

```go

regionInFrontOfMe _:= aRegion.Front()

moveOrder, err_ := field.MakeOrderMoveMaxSpeed(*me.Position, regionInFrontOfMe.Center())

```

## Next steps

As you may have noticed, the bot player in the example directory does not play well. 
Now, you may start your bot by implementing its behaviour.  

## Deploying you bots

After developing your bot, you may share it with other developers.

Using this client your code will be compiled into a binary file. So you do not have to share the bot source code.

There is a Dockerfile template in [the example directory](./examples) to guide you how to create a container.

You may create an account on https://lugobots.dev and upload your bot.
