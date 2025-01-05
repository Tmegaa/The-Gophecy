@echo off

:: Teste 1
echo Running simulation with 10 agents and strategies [Random, Patrol, HeatMap]
go run . -agents 10 -strategies1 0 -strategies2 1 -strategies3 2 -time 3

:: Teste 2
echo Running simulation with 10 agents and strategies [Patrol, HeatMap, Center of Mass]
go run . -agents 10 -strategies1 1 -strategies2 2 -strategies3 3 -time 3

:: Teste 3
echo Running simulation with 20 agents and strategies [Random, HeatMap, Center of Mass]
go run . -agents 20 -strategies1 0 -strategies2 2 -strategies3 3 -time 3

:: Teste 4
echo Running simulation with 20 agents and strategies [Patrol, Random, HeatMap]
go run . -agents 20 -strategies1 1 -strategies2 0 -strategies3 2 -time 3

:: Teste 5
echo Running simulation with 30 agents and strategies [HeatMap, Center of Mass, Random]
go run . -agents 30 -strategies1 2 -strategies2 3 -strategies3 0 -time 3

:: Teste 6
echo Running simulation with 30 agents and strategies [Center of Mass, Patrol, Random]
go run . -agents 30 -strategies1 3 -strategies2 1 -strategies3 0 -time 3

:: Teste 7
echo Running simulation with 40 agents and strategies [Random, Patrol, Center of Mass]
go run . -agents 40 -strategies1 0 -strategies2 1 -strategies3 3 -time 3

:: Teste 8
echo Running simulation with 40 agents and strategies [HeatMap, Random, Patrol]
go run . -agents 40 -strategies1 2 -strategies2 0 -strategies3 1 -time 3

:: Teste 9
echo Running simulation with 50 agents and strategies [Center of Mass, HeatMap, Patrol]
go run . -agents 50 -strategies1 3 -strategies2 2 -strategies3 1 -time 3

:: Teste 10
echo Running simulation with 50 agents and strategies [Patrol, Center of Mass, Random]
go run . -agents 50 -strategies1 1 -strategies2 3 -strategies3 0 -time 3

:: Teste 11
echo Running simulation with 60 agents and strategies [Random, HeatMap, Center of Mass]
go run . -agents 60 -strategies1 0 -strategies2 2 -strategies3 3 -time 3

:: Teste 12
echo Running simulation with 60 agents and strategies [HeatMap, Patrol, Center of Mass]
go run . -agents 60 -strategies1 2 -strategies2 1 -strategies3 3 -time 3

pause
