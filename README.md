# Projet IA04 - A24 : La Gophétie

## 🏫 Description

Dans cette simulation multi-agents, nous explorons l'évolution de leur croyance: seront-ils fidèles au langage Go?

Les agents, des étudiants d'ingénierie informatique au sein d'un campus, sont plus ou moins adhérents à la doctrine du langage Go. Les plus croyants veulent persuader leurs camarades de la supériorité de ce magnifique langage de programmation, alors que les plus sceptiques ont pour mission de dissuader les autres. Dans cette simulation nous allons nous poser une question: **Quelles politiques d'embrigadement fonctionnent le mieux ?**

## 🔗 Recupérer le projet du repository (git)

```{bash}
go env -w GOPRIVATE=github.com/Tmegaa/*
go install github.com/Tmegaa/The-Gophecy@latest
```

## 🔬  Tests avec différents cas de figure

> TODO: fill this

## 💻 La Gophétie

### 1. 📐 L'architecture

- backend: **go**
- frontend: **ebiten**

Packages:

- **agent**: gestion des agents, de l'environnement et des objets
- **carte**: gestion de la carte
- **simulation**: gestion de la simulation (l'affichage graphique, les interactions avec l'utilisateur…)
- **tile**: gestion des jeux de tuiles (soit les éléments sur la carte)
- **utils**: constantes et fonctions qui sont utiles dans les autres packages
- **gophecy**: contient le "main"

Une modélisation des éléments de cette simulation:

![UML](/pdf/UML_Classe.png "UML des classes")

### 2.🚶 Les agents

Les agents sont des étudiants en ingénierie informatique et ont donc des fortes opinions vis-à-vis des langages de programmation. Dans cette simulation, on peut considérer que ces croyances sont un peu sectaires... De plus, cette simulation a lieu dans un campus d'université, les agents peuvent donc se déplacer librement, mais ils auront des preferences par rapport à leur façon de bouger.

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

Après une modification de l'opinion d'un agent, on vérifie son type et on le met à jour si besoin: les types ne sont donc pas statiques tout au long de la simulation, ils peuvent évoluer.

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

Aucun sous-type n'est possible pour les agents neutres. Cependant les croyants et les sceptiques ont la possibilité d'être des pirates ou des convertisseurs. Ces sous-types rentrent en jeu dans le cas où un agent pourrait percevoir à la fois un ou plusieurs agents proches en plus d'un ou plusieurs objets à proximité.

Si le choix est présenté, un pirate va choisir d'interagir avec un ordinateur plutôt qu'engager une conversation avec un autre agent. Pour les convertisseurs c'est l'inverse.

Les croyants vont avoir une plus grande tendance à être des convertisseurs alors que les sceptiques auront plus tendance à être des pirates.

La probabilité d'avoir un sous-type est de 70%.

### 3. 📈 L'évolution des croyances

Il y a trois actions qui font évoluer les croyances des agents: prier, discuter et utiliser un ordinateur.

Lors de l'utilisation d'un ordinateur, les sceptiques diminuent leur opinion de Go (et le désinstallent si installé), contrairement aux croyants qui l'augmentent (en installant Go). Les agents neutres vont voir leur opinion diminuer ou augmenter en fonction de si Go est installé ou pas.

La prière n'est disponible que pour les agents croyants et neutres: elle fait augmenter la croyance en Go, d'autant plus pour les agents neutres (qui décident d'agir en fonction de leur foi). Les agents neutres vont cependant avoir moins de probabilités de choisir la prière.

Enfin, la façon la plus intéressante de faire évoluer les croyances des agents est la discussion: dans le cas où un agent croyant et un sceptique décident de parler, ils ne font qu'amplifier leur opinion de base. En effet, le croyant voit son opinion augmenter et le sceptique voit la sienne diminuer. C'est une modélisation de deux personnes têtues qui ne vont pas pouvoir écouter des arguments qu'ils jugent presque "extrémistes" de l'autre.

D'un autre côté, les discussions entre un agent neutre et tout autre type d'agent vont voir intervenir bien plus de paramètres: nous voyons entrer en jeu les relations entre les agents, un certain degré de charisme qui donne un certain poids aux conversations...

> Nous avons basé la modélisation sur plusieurs articles, que l'on peut trouver dans le dossier "/pdf" de ce projet. De plus, le document [Résumé et Analyse : Modèle d’Endoctrinement par équations Différentielles](./pdf/Indoctrination_equation%20(1).pdf) détaille toutes les équations.

Tout d'abord, on modélise les relations entre les agents. Un agent peut avoir une des relations suivantes avec un autre agent:

- Ennemi
- Amis
- Famille
- Pas de lien direct / Inconnu
  
Cette relation va avoir un effet sur le calcul des poids absolus. Pour chaque agent, nous allons attribuer le poids qu'il donne à l'opinion d'un autre agent. Il va être beaucoup plus confiant d'un ami que d'un inconnu par exemple. Ces poids absolus sont normalisés. Un agent va avoir une certaine confiance envers lui-même, un poids absolu qu'il donne à ses propres opinions, qui se traduit par la valeur référencée par son propre ID dans son dictionnaire de poids absolus.

Pour les poids relatifs, ce paramètre de confiance en soi rentre en jeu. En effet, un agent A va avoir une certaine confiance générale sur sa propre opinion (poids absolu), une certaine confiance de sa propre opinion en parlant avec un agent B (poids relatif 1) et une certaine confiance dans l'opinion de l'agent B tout en prenant en compte non seulement leur relation mais aussi sa propre confiance (poids relatif 2).

$$
\displaystyle Rel_{A\to A /B}=\frac{Abs_{A\to A}}{Abs_{A\to A}+Abs_{A\to B}} \quad Rel_{A\to B/B}=\frac{Abs_{A\to B}}{Abs_{A\to A}+Abs_{A\to B}}
$$

Chaque agent a en plus un paramètre personnel qui symbolise sa réceptivité.

Lors d'une conversation, nous avons modélisé la mise à jour des opinions des agents A et B de la façon suivante (cf. [source](./pdf/Indoctrination_equation%20(1).pdf) pour plus de détails):

$$
\displaystyle NewO_{A} = Rel_{A\to A /B} * K_{A} * OldO_{A} * (1-OldO_{A}) + Rel_{A\to B /B} * OldO_{B}
$$
$$
\displaystyle NewO_{B} = Rel_{B\to A /A} * OldO_{A} + Rel_{B\to B /A} * OldO_{B} * K_{B} * OldO_{B} * (1-OldO_{B})
$$

- K est le paramètre personnel
- NewO est la nouvelle opinion
- OldO est l'opinion courante
- Rel est le poids relatif que donne le premier agent à l'opinion du deuxième en connaissant l'interlocuteur.

Nous avions prévu de rajouter un paramètre de Charisme qui serait l'influence perçue d'un agent A sur un agent B, mais ceci n'as pas été implémenté.

### 4. 🏃 Les stratégies de mouvement

## 😇  Les Gophètes

- Lepretre Thomas
- Perdereau Tom
- Saby Loyola Sophie
- Sporck Trombini Gabriel
