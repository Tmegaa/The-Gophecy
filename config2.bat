@echo off
echo Running simulation with 50 agents and strategies [Center of Mass, HeatMap, Patrol] [Convincente, Independente, Explorador]
go run . -agents 50 -strategies1 3 -strategies2 2 -strategies3 1 -time 3 -strategiesA1 0 -strategiesA2 1 -strategiesA3 2

echo Running simulation with 50 agents and strategies [Patrol, Random, HeatMap] [Independente, Explorador, Convincente]
go run . -agents 50 -strategies1 1 -strategies2 0 -strategies3 2 -time 3 -strategiesA1 1 -strategiesA2 2 -strategiesA3 0

echo Running simulation with 60 agents and strategies [HeatMap, Center of Mass, Random] [Explorador, Convincente, Independente]
go run . -agents 60 -strategies1 2 -strategies2 3 -strategies3 0 -time 3 -strategiesA1 2 -strategiesA2 0 -strategiesA3 1

echo Running simulation with 60 agents and strategies [Random, Patrol, Center of Mass] [Independente, Convincente, Explorador]
go run . -agents 60 -strategies1 0 -strategies2 1 -strategies3 3 -time 3 -strategiesA1 1 -strategiesA2 0 -strategiesA3 2

echo Running simulation with 40 agents and strategies [Center of Mass, HeatMap, Patrol] [Convincente, Explorador, Independente]
go run . -agents 40 -strategies1 3 -strategies2 2 -strategies3 1 -time 3 -strategiesA1 0 -strategiesA2 2 -strategiesA3 1

echo Running simulation with 40 agents and strategies [Patrol, Random, HeatMap] [Explorador, Independente, Convincente]
go run . -agents 40 -strategies1 1 -strategies2 0 -strategies3 2 -time 3 -strategiesA1 2 -strategiesA2 1 -strategiesA3 0

echo Running simulation with 30 agents and strategies [HeatMap, Center of Mass, Random] [Independente, Convincente, Explorador]
go run . -agents 30 -strategies1 2 -strategies2 3 -strategies3 0 -time 3 -strategiesA1 1 -strategiesA2 0 -strategiesA3 2

echo Running simulation with 30 agents and strategies [Random, Patrol, Center of Mass] [Explorador, Convincente, Independente]
go run . -agents 30 -strategies1 0 -strategies2 1 -strategies3 3 -time 3 -strategiesA1 2 -strategiesA2 0 -strategiesA3 1

echo Running simulation with 20 agents and strategies [Center of Mass, HeatMap, Patrol] [Convincente, Independente, Explorador]
go run . -agents 20 -strategies1 3 -strategies2 2 -strategies3 1 -time 3 -strategiesA1 0 -strategiesA2 1 -strategiesA3 2

echo Running simulation with 20 agents and strategies [Patrol, Random, HeatMap] [Independente, Explorador, Convincente]
go run . -agents 20 -strategies1 1 -strategies2 0 -strategies3 2 -time 3 -strategiesA1 1 -strategiesA2 2 -strategiesA3 0

echo Running simulation with 10 agents and strategies [HeatMap, Center of Mass, Random] [Explorador, Convincente, Independente]
go run . -agents 10 -strategies1 2 -strategies2 3 -strategies3 0 -time 3 -strategiesA1 2 -strategiesA2 0 -strategiesA3 1

echo Running simulation with 10 agents and strategies [Random, Patrol, Center of Mass] [Independente, Convincente, Explorador]
go run . -agents 10 -strategies1 0 -strategies2 1 -strategies3 3 -time 3 -strategiesA1 1 -strategiesA2 0 -strategiesA3 2

echo Running simulation with 50 agents and strategies [HeatMap, Patrol, Center of Mass] [Convincente, Explorador, Independente]
go run . -agents 50 -strategies1 2 -strategies2 1 -strategies3 3 -time 3 -strategiesA1 0 -strategiesA2 2 -strategiesA3 1

echo Running simulation with 50 agents and strategies [Random, HeatMap, Patrol] [Independente, Convincente, Explorador]
go run . -agents 50 -strategies1 0 -strategies2 2 -strategies3 1 -time 3 -strategiesA1 1 -strategiesA2 0 -strategiesA3 2

echo Running simulation with 60 agents and strategies [Center of Mass, Random, HeatMap] [Explorador, Convincente, Independente]
go run . -agents 60 -strategies1 3 -strategies2 0 -strategies3 2 -time 3 -strategiesA1 2 -strategiesA2 0 -strategiesA3 1

echo Running simulation with 60 agents and strategies [Patrol, Center of Mass, Random] [Convincente, Explorador, Independente]
go run . -agents 60 -strategies1 1 -strategies2 3 -strategies3 0 -time 3 -strategiesA1 0 -strategiesA2 2 -strategiesA3 1

echo Running simulation with 40 agents and strategies [HeatMap, Patrol, Center of Mass] [Independente, Convincente, Explorador]
go run . -agents 40 -strategies1 2 -strategies2 1 -strategies3 3 -time 3 -strategiesA1 1 -strategiesA2 0 -strategiesA3 2

echo Running simulation with 40 agents and strategies [Random, HeatMap, Patrol] [Explorador, Convincente, Independente]
go run . -agents 40 -strategies1 0 -strategies2 2 -strategies3 1 -time 3 -strategiesA1 2 -strategiesA2 0 -strategiesA3 1

echo Running simulation with 30 agents and strategies [Center of Mass, Random, HeatMap] [Convincente, Explorador, Independente]
go run . -agents 30 -strategies1 3 -strategies2 0 -strategies3 2 -time 3 -strategiesA1 0 -strategiesA2 2 -strategiesA3 1

echo Running simulation with 30 agents and strategies [Patrol, HeatMap, Center of Mass] [Independente, Convincente, Explorador]
go run . -agents 30 -strategies1 1 -strategies2 2 -strategies3 3 -time 3 -strategiesA1 1 -strategiesA2 0 -strategiesA3 2

pause