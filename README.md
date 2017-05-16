# go mahjong
[![Build Status](https://travis-ci.org/mapleque/go_mahjong.png)](https://travis-ci.org/mapleque/go_mahjong)

## UI
```
----> login.html <----> register.html
         |
         +----> game_list.html <----+-----------+
         |          |               |           |
         +----------+----> ready_view.html      |
         |          |           |               |
         +----------+-----------+----> playing_view.html
```
## API
```
login
    >username
    >password
    <token
info
    >token
    <gameList
    <gameInfo
create
    >token
    >option
join
    >token
    >gameId
next
    >token
leave
    >token
op
    >token
    >op
    >indexs
```
## FSM
```
player create
|playerList = []  |playerList = [...]  |playerChain        |playerChain
|rule             |rule                |rule               |ruleGuard
|game = prepare   |game = prepare      |game = ing         |game = ing
|set = prepare    |set = prepare       |set = prepare      |set = ing
|---------------->|------------------->|------------------>|
 player join       playerCount match    player ready
                                    start game           init set

init ruleGuard
|playerEventStack = []  |playerEventStack = [...]  |                  |playerEventStack = [...]
|                       |                          |curPlayerEvent    |
|---------------------->|------------------------->|----------------->|
  init player event       pop player event           do event or pass
                                                     rebuild palyer event or
                                                     pop next player event

if player event stack empty
|game = ing    |game = ing     |game = ing
|set = end     |set = prepare  |set = ing
|              |winInfo        |
|------------->|-------------->|
  calc result    player ready
  restart set    init set

player event stack
init stack
|{player 0 PULL4}        |{player x BUHUA}    |{player x after pull}   |{player x+1 HU}
|{player 1 PULL4}        |------------------->|                        |{player x+2 HU}         |empty
|{player 2 PULL4}        | op is BUHUA                                 |{player x+3 HU}         |
|{player 3 PULL4}        |                                             |----------------------->|
| ... * 3                |                                             | op is HU
|{player 0 PULL2}        |{player x HU,GANG}  |empty                   |                        |empty
|{player 1 PULL1}        |                    |                        |{player x+1 GANG}       |{player y+1 HU}   |empty
|{player 2 PULL1}        |------------------->|                        |{player x+2 GANG}       |{player y+2 HU}   |
|{player 3 PULL1}        | op is HU                                    |{player x+3 GANG}       |{player y+3 HU}   |
|{player 0 BUHUA}        |                                             |----------------------->|----------------->|
|{player 1 BUHUA}        |                                             | y op is GANG           | op is HU
|{player 2 BUHUA}        |                    |{player x CI}           |                        |
|{player 3 BUHUA}        |                    |{player x after pull}   |{player x+1 PENG}       |{player y CI}
|{player 0 after pull}   |------------------->|                        |{player x+2 PENG}       |{player y after pull}
|                        | op is GANG                                  |{player x+3 PENG}
|                        |                                             |{player x+1 CHI}        |empty
|                        |                                             |----------------------->|{player y PUSH}
|                        |                                             | y op is PENG or CHI    |{player y after push}
|                        |                    |{player x after push}   |
|                        |{player x PUSH}     |                        |{palyer x+1 PULL}
|                        |                    |                        |{player x+1 after pull}
|----------------------->|------------------->|----------------------->|
  after pull               all pass is push     after push

```
