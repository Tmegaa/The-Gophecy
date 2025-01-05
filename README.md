# README : Configurations et Stratégies pour les Simulations d'Agents

Ce document décrit les configurations utilisées dans diverses simulations impliquant des agents dotés de différents types et stratégies de mouvement. Les configurations sont organisées par scénarios, avec une introduction aux types d'agents et aux stratégies appliquées.

---

## Comment Utiliser les Fichiers `.bat`

Les fichiers `.bat` sont des scripts qui exécutent des simulations avec différentes configurations. Pour exécuter un fichier `.bat`, il suffit run dans le terminal le nom du fichier, par exemple :

`config1.bat`

## Types d'Agents

1. **Convaincant** :

   - **Caractéristiques** : Forte capacité à influencer d'autres agents et à changer rapidement d'opinion.
   - **Rôle** : Élément dynamique modifiant les schémas émergents dans l'environnement.

2. **Indépendant** :

   - **Caractéristiques** : Concentration sur des objectifs individuels et résistance à l'influence extérieure.
   - **Rôle** : Maintient la stabilité et explore de manière autonome.

3. **Explorateur** :
   - **Caractéristiques** : Priorise l'exploration et la découverte de nouvelles zones.
   - **Rôle** : Élargit la couverture de l'environnement avec des interactions sociales limitées.

---

## Stratégies de Mouvement

1. **Aléatoire (Random)** :

   - Mouvement imprévisible pour maximiser l'exploration.

2. **Patrouille (Patrol)** :

   - Parcours prédéfinis assurant une couverture cohérente.

3. **Carte de Chaleur (HeatMap)** :

   - Ciblage des zones présentant une activité ou un intérêt élevé.

4. **Centre de Masse (Center of Mass)** :
   - Convergence vers des zones densément peuplées par d'autres agents.

---

## Scénarios de Simulation (Première Configuration)

### Tests avec Différentes Stratégies de Mouvement

1. **Test 1** : Simulation avec une seule statue.

   - **Objectif** : Étudier l'exploration et les interactions des agents autour d’un seul point d'intérêt fixe.
   - **Résultats Attendus** :
     - _Carte de Chaleur (HeatMap)_ : Cartographie des environs de la statue.
     - _Aléatoire (Random)_ : Introduit de la variabilité dans le mouvement.
     - _Patrouille (Patrol)_ : Assure un schéma de surveillance autour de la statue.

2. **Test 2** : Simulation avec trois statues.

   - **Objectif** : Évaluer comment les agents répartissent leur attention entre plusieurs points d'intérêt.
   - **Résultats Attendus** :
     - _Centre de Masse (Center of Mass)_ : Concentre les agents autour des statues.
     - _Patrouille (Patrol)_ : Surveille les zones proches des trois statues.

3. **Test 3** : Simulation avec trois statues et 15 PC.

   - **Objectif** : Observer l'équilibre entre l'exploration des statues et l'interaction avec les PC interactifs.
   - **Résultats Attendus** :
     - _Carte de Chaleur (HeatMap)_ et _Aléatoire (Random)_ : Essentiels pour localiser les PC.
     - _Patrouille (Patrol)_ et _Centre de Masse (Center of Mass)_ : Priorisent les zones proches des statues.

4. **Test 4** : Introduction de Nouveaux Types d'Agents.
   - **Objectif** : Analyser les interactions entre les agents **Convaincant**, **Indépendant** et **Explorateur**.
   - **Résultats Attendus** :
     - La dynamique apportée par les agents **Convaincant** modifie les schémas de comportement.
     - Les différentes stratégies influencent l'efficacité collective.

---

## Configurations des Simulations (Fichiers `.bat`)

### Configuration 1 (`config1.bat`)

Simulations avec différents nombres d'agents (10, 20, 30, 40, 50, 60) et combinaisons de stratégies [Random, Patrol, HeatMap, Center of Mass].

### Configuration 2 (`config2.bat`)

Ajout des nouveaux types d'agents (**Convaincant**, **Indépendant**, **Explorateur**) avec les combinaisons suivantes :

- Les stratégies de mouvement et les types d'agents sont mappés pour analyser comment différents paramètres interagissent dans l'environnement.

