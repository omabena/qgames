# qgames parser 

### Instructions
Solution is implemented in Go to take that takes advantage of concurrecy features to allow for concurrent processing of some parts of the flows.

Reports will be run concurrently and the results will be printed to the console.
### How to run
There are two ways to run the project and they are described below:

Please ensure that you have Go 1.22 installed on your machine.
Also, that the .env file is present in the root of the project, LOG_FILE_PATH is set to the path of the log file and is accessible to the program, as default is set to `logs/qgames.log` in the root of the project.

1. If you have Go 1.22 installed on your machine, you can run the following command:
```shell
go run main.go qgames-parser run
```

2. A binary is provided in the root of the project. You can run the following command:
```shell
./qgameslogs qgames-parser run
```

To finish the execution, you can press `Ctrl + C` to stop the execution.

### Output
The output will be printed for three reports into the console and will be in the following format:

After running is expected to see three reports "Games report", "Players ranking report" and deaths grouped by cause or MOD Death

***The order of the responses may vary as they are processed concurrently, and there's not guarantee on which one will finish first***


#### Games report output (complete output is not shown here)

````
***************
***************
Games report
Game: game_1
Total Kills: 0
Players: [Isgalamido]
Kills {
}
Game: game_2
Total Kills: 15
Players: [Isgalamido, Dono da Bola, Mocinha, Zeh]
Kills {
	Isgalamido: 4
}
Game: game_3
Total Kills: 105
Players: [Assasinu Credi, Dono da Bola, Isgalamido, Zeh]
Kills {
	Dono da Bola: 20
	Isgalamido: 27
	Zeh: 22
	Assasinu Credi: 16
}
Game: game_4
Total Kills: 14
Players: [Isgalamido, Zeh, Assasinu Credi, Dono da Bola]
Kills {
	Assasinu Credi: 5
	Isgalamido: 2
	Zeh: 2
}
Game: game_5
Total Kills: 29
Players: [Zeh, Dono da Bola, UnnamedPlayer, Isgalamido, Oootsimo, Maluquinho, Assasinu Credi, Mal, Fasano Again]
Kills {
	Oootsimo: 9
	Isgalamido: 4
	Zeh: 8
	Dono da Bola: 2
	Maluquinho: 1
	Assasinu Credi: 1
}

.....
````

#### Players ranking report output 

````
***************
Players Ranking
Isgalamido: 34
Oootsimo: 11
Maluquinho: 0
Chessus!: 0
Fasano Again: 0
UnnamedPlayer: -1
Mocinha: -2
Chessus: -13
Zeh: -15
Assasinu Credi: -47
Dono da Bola: -86
Mal: -121
`````


#### Deaths grouped by cause or MOD Death (complete output is not shown here)

````
***************
Death Mode report
Game: game_1
Game: game_2
Mod: MOD_FALLING, Count: 2
Mod: MOD_ROCKET, Count: 1
Mod: MOD_TRIGGER_HURT, Count: 9
Mod: MOD_ROCKET_SPLASH, Count: 3
Game: game_3
Mod: MOD_RAILGUN, Count: 8
Mod: MOD_ROCKET_SPLASH, Count: 51
Mod: MOD_MACHINEGUN, Count: 4
Mod: MOD_SHOTGUN, Count: 2
Mod: MOD_TRIGGER_HURT, Count: 9
Mod: MOD_FALLING, Count: 11
Mod: MOD_ROCKET, Count: 20
Game: game_5
Mod: MOD_SHOTGUN, Count: 4
Mod: MOD_ROCKET_SPLASH, Count: 13
Mod: MOD_TRIGGER_HURT, Count: 3
Mod: MOD_FALLING, Count: 1
Mod: MOD_MACHINEGUN, Count: 1
Mod: MOD_ROCKET, Count: 5
Mod: MOD_RAILGUN, Count: 2
Game: game_6
 
.....
````



The script could be modified to also print into a files or any other output format, as the format is not specified in the requirements.



### Running tests
If interested in running tests you can run
```shell
go test -v ./...
```