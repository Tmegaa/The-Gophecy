# Projet IA04 - A24 : La Gophétie

## Description

Dans cette simulation multi-agents, nous explorons l'évolution de leur croyance: seront-ils fidèles au langage Go?

Les agents, des étudiants d'ingénierie informatique au sein d'un campus, sont plus ou moins adhérents à la doctrine du langage Go. Les plus croyants veulent persuader leurs camarades de la supériorité de ce magnifique langage de programmation, alors que les plus sceptiques ont pour mission de dissuader les autres. Dans cette simulation nous allons nous poser une question: **Quelles politiques d'embrigadement fonctionnent le mieux ?**

## Recupérer le projet du repository (git)

```{bash}
go env -w GOPRIVATE=github.com/Tmegaa/*
go install github.com/Tmegaa/The-Gophecy@latest
```

## Tests avec différents cas de figure

> TODO: fill this

## La Gophétie

### 1. L'architecture

- backend: go
- frontend: ebiten

Packages:

- agent: gestion des agents, de l'environnement et des objets
- carte: gestion de la carte
- simulation: gestion de la simulation (l'affichage graphique, les interactions avec l'utilisateur…)
- tile: gestion des jeux de tuiles (soit les éléments sur la carte)
- utils: constantes et fonctions qui sont utiles dans les autres packages
- gophecy: contient le "main"

Une modélisation des éléments de cette simulation: IMAGE

### 2. Les agents

Les agents sont des étudiants en ingénierie informatique et ont donc des fortes opinions vis-à-vis des langages de programmation. Dans cette simulation, on peut considérer que ces croyances sont un peu sectaires. De plus, cette simulation a lieu dans un campus d'université, les agents peuvent donc se déplacer librement, mais ils auront des preferences par rapport à leur façon de bouger.

Dans la boucle de perception, délibération et action de chaque agent, il y a un temps d'attente de 20ms entre chaque boucle.

Tous les agents ont la même fonction de perception où ils reçoivent de l'environnement une liste des agents et des objets qui sont à une certaine distance. Cet aire de perception, qui sera affichée comme un rectangle, va dépendre de l'acuité de l'agent. Il pourra donc délibérer.

Nous verrons par la suite que les sous-types interviennent dans la prise de décision. Les agents vont donc choisir parmi les actions suivantes:

- **Bouger** : L'agent va se déplacer, avec ou sans but. Ses déplacements ont une durée limitée. Tout agent va choisir cette option s'il ne perçoit aucun autre agent ou objet à proximité, mais aussi à la fin des autres actions. C'est "l'action par défaut".
- **Utiliser un ordinateur** : L'agent va pouvoir accéder à un ordinateur.
- **Prier** : Certains agents peuvent prier auprès d'une statue.
- **Discuter** : Deux agents peuvent s'engager dans une conversation avec une durée limitée. Chaque agent a un paramètre "MaxLastTalked" qui indique avec combien de personnes il se rappelle d'avoir discuté, la liste de ses derniers "MaxLastTalked" interlocuteurs est sauvegardée et constamment mise à jour pour éviter qu'un agent parle trop souvent aux mêmes personnes.
- **Attendre** : Il ne va réaliser aucune action pendant une boucle. Il est envisageable par la suite d'implémenter un "temps d'attente", mais pour l'instant cette action n'a d'effet que pendant une seule boucle de perception, délibération et action.

#### 2.1 Les types d'agents

Le degré de croyance dans le langage Go est modélisé chez chaque agent par une variable "Opinion" qui prend comme valeur un float entre 0 et 1, 0 représentant un scepticisme total et 1 une croyance aveugle. En fonction de leur degré de croyance, les agents prendront un de ces 3 types:

Opinion|Type|Description|
:--------------: | :--------------: |------------- |
[0, 0.33[| Sceptique| Ne croit pas dans le langage Go et va essayer des dissuader ses camarades de l'utiliser.|
[0.33, 0.66]| Neutre| Est mitigé et va être influencé par tous les autres agents.|
]0.66, 1]| Croyant| Croit que le langage Go est incroyable et aura pour mission de répandre sa croyance en plus d'essayer de l'augmenter.|

Le type de chaque agent va influencer son comportement, particulièrement dans 4 domaines:

 1. **Leurs interactions avec d'autres agents** : Les conversations entre agents d'un même type ou entre types différents vont avoir des effets différents. Ceci sera détaillé plus tard dans ce rapport.
 2. **Leurs patrons de mouvement** : Chaque type va évoluer dans l'espace de façon différente. Nous verrons ceci plus en détail après.
 3. **Leurs choix de comportement** : les agents croyants et sceptiques pourront avoir un sous-type qui va influencer leurs choix.
 4. **Leurs actions spécifiques** : Nous pouvons voir que les actions ne sont pas réalisés de la même façon par tous les agents:

Action\Type Agent| Sceptique| Neutre| Croyant|
------------- | :--------------: | :--------------: |:--------------: |
Bouger| Type mouvement 1 | Type mouvement 2 |Type mouvement 3 |
Utiliser un ordinateur | Désinstalle Go | Regarde quel langage de programmation est installé | Installe Go|
Prier | (action non réalisable) | Prie auprès d'une statue | Peut prier auprès d'une statue|

#### 2.2 Les sous-types d'agents

En effet...

### 3. L'évolution des croyances

### 4. Les stratégies de mouvement

## Les Gophètes

- Lepretre Thomas
- Perdereau Tom
- Saby Loyola Sophie
- Sporck Trombini Gabriel
