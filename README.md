# Projet IA04 - A24 : La Gophétie

## 🏫 Description

Dans cette simulation multi-agents, nous explorons l'évolution de leur croyance: seront-ils fidèles au langage Go?

Les agents, des étudiants d'ingénierie informatique au sein d'un campus, sont plus ou moins adhérents à la doctrine du langage Go. Les plus croyants veulent persuader leurs camarades de la supériorité de ce magnifique langage de programmation, alors que les plus sceptiques ont pour mission de dissuader les autres. Dans cette simulation nous allons nous poser une question: **Quelles politiques d'embrigadement fonctionnent le mieux ?**

## 🔗 Recupérer le projet du repository (git)

Pour simplement récupérer le module et pouvoir faire tourner la simulation:

```{bash}
go install github.com/Tmegaa/The-Gophecy@latest
```
Les fichiers se trouvent dans le GOPATH dans le dossier `pkg/mod/gitlab.utc.fr/`
```{bash}
go run .
```

Dans le cas où vous voudriez récupérer tout le projet (notamment les sources dans le dossier "/pdf"):

```{bash}
git clone https://github.com/Tmegaa/The-Gophecy.git
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

#### 2.3. 📈 L'évolution des croyances

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

#### 2.4 🏃 Les stratégies de mouvement

Chaque type d'agent va avoir une stratégie de mouvement différente. cette stratégie pourra être assignée lors du début de la simulation par l'utilisateur, et c'est envisageable de la prédéfinir avec des fichiers de configuration.

Les 4 stratégies de mouvement sont:

- **Random** : cette stratégie est la plus simple car une direction est choisie aléatoirement.
- **Patrol** : l'agent va choisir un point vers lequel se diriger dans la carte. Il va choisir plusieurs point aléatoirement au début, puis il choisira le meilleur en lui assignant un score qui va dépendre de la distance à parcourir pour arriver à ce point, les potentiels obstacles à éviter et un facteur aléatoire. Ce point peut rester constant tout le long de la simulation s'il n'est pas atteint, mais l'agent a aussi la possibilité de changer de point s'il atteint la position ou de prendre une direction aléatoire.
- **HeatMap** : les agents maintiennent un historique des positions qu'ils ont déjà visité. Avec cette stratégie, les agents vont essayer de se diriger vers les zones qu'ils ont personnellement visité le moins afin de parcourir des nouvelles positions le plus possible.
- **Center of Mass** : les agents vont chercher à se déplacer vers le centre de congrégations. Soit, en calculant le centre de masse des agents aux alentours, ces agents vont avoir comme objectif dans leur déplacement un point qui les rapprochera le plus possible au plus grand nombre d'agents possible. Il y a tout de même une petite chance de passer à un mouvement aléatoire pour éviter un regroupement excessif.

Pour l'instant la vitesse des agents indiquée lors de la création n'a pas d'effet dans leur déplacement, pour notre simulation il n'est pas vital que les agents bougent à des vitesses différentes. Une modification à envisager par la suite serait l'implémentation des vitesses.

### 3. ▶️ La simulation

Le backend est (évidemment) réalisé en Go, mais pour l'affichage nous avons utilisé Ebiten.

Lorsqu'on lance la simulation, on donne le nombre d'agents, la durée de la simulation et les stratégies de mouvement par type. La simulation est donc initialisée et l'affichage graphique est ouvert dans une autre fenêtre.

Tout d'abord nous pouvons observer l'affichage (cette simulation comptait 40 agents):

![simu1](/images/simu_all.png "Capture d'écran de la simulation")

A gauche nous pouvons voir les informations pertinentes de la simulation, tel de que temps écoulé, la répartition des agents par type et le langage de programmation installé sur les ordinateurs. A droite nous observons la carte avec les agents, les objets...

![simu2](/images/three_agents.png "Capture d'écran de trois agents")

Chaque type d'agent est affiché avec une image différente: les croyants sont en noir, les sceptiques sont en rouge et les agents neutres sont en blanc. Le carré qui les entoure est leur zone de perception.

![simu3](/images/simu_click_agent.png "Capture d'écran affichage infos agent")

Lorsqu'on clique sur un agent (ici nous pouvons le voir tout à droite, un sceptique rouge entouré d'un carré jaune), nous pouvons lire sur le bandeau de gauche des informations pertinentes sur cet agent telles que son action courante, son historique de discussions, son paramètre personnel de réceptivité...

Pour l'instant, il n'est pas encore possible de faire un scroll dans ce bandeau, il n'est donc pas possible de voir toutes les relations que cet agent a avec le reste.

Lorsque l'agent sélectionné est en discussion avec un autre, nous avons les cette information aussi.

![simu4](/images/discussion_infos.png "Capture d'écran affichage informations sur discussion")

Un click sur un ordinateur nous donne des informations aussi:

![simu5](/images/click_computer.png "Capture d'écran affichage infos ordinateur")

Chaque action, autre que le mouvement, affiche une petite boîte en dessus de l'agent avec le nom de l'action et une barre qui indique le temps restant pour compléter cette action. Si cette action est une discussion, on affiche aussi le type de chaque agent.

![simu6](/images/discussion.png "Capture d'écran affichage action")

Lorsque la simulation finit, nous avons un petit compte-rendu avec la répartition des agents par type finale et l'opinion moyenne de tous les agents par rapport à Go.

Par exemple, ici on a les résultats d'une simulation de 50 agents dont les stratégies de mouvement étaient toutes aléatoires:

![simu7](/images/simu_end.png "Capture d'écran affichage à la fin d'une simulation")

De plus, un graphique détaillant l'opinion globale sur Go en fonction du temps écoulé est sauvegardé. Pour la même simulation nous obtenons:

![simu8](/images/results_example.png "Graphique représentant la croyance moyenne de la population en fonction du temps")

### 4.💡 Idées pour la suite

Tout au long de ce rapport nous avons vu des améliorations possibles pour ce projet. Nous pouvons en explorer d'avantage.

En effet, nos agents ont un booléen qui indique s'ils sont vivants ou pas. Avec cette version de la simulation, il n'est pas possible de mourir, cependant il serait envisageable de rajouter des fonctionnalités en rapport à la santé des agents: un agent fatigué ou affamé pourrait être beaucoup plus influençable qu'un agent en pleine santé! Des paramètres de faim ou d'énergie avec des actions de type "Manger" ou "Dormir" (des fonctions ont été laissées en commentaire pour montrer l'emplacement des fonctions dans notre architecture) seraient donc rajoutées à nos agents. La conclusion si un agent est beaucoup trop affamé ou beaucoup trop fatigué? Notre booléen prendrait la valeur `false`.

De plus, nous explorons ici l'opinion vis à vis de Go, mais l'évolution que nous avions prévu de base pour cette simulation serait l'introduction d'autres sectes! Que ce soit le C++ulte, le Hask Hell
la BASH astrée ou l'HTMLM, il serait très intéressant d'observer la concurrence des différentes croyances au sein d'une même population.

Nous avions pensé à une liste (ou un map, peu importe), au lieu d'une seule valeur modélisant l'opinion d'une personne. Il y aurait à priori plus de types de croyants et des questions à se poser:

- Est-ce qu'on peut être croyant pour une seule secte ou pour plusieurs?
- Dans le cas où un agent devient croyant, que deviennent ses autres opinions?
- Quel est la nouvelle signification du scepticisme?

De plus, étant donné que cette simulation a lieu au sein d'un campus universitaire, nous pourrions rajouter des personnages tels que des professeurs qui instruisent un groupe d'agents.

Finalement, pour l'instant l'utilisateur ne peut pas intervenir dans la simulation. Il donne les paramètres de départ, mais il ne peut pas agir après en dehors de la visualisation d'informations. Nous pourrions envisager d'ajouter des boutons qui permettraient à l'utilisateur de rajouter des agents ou des objets en cours de route.

## 😇  Les Gophètes

- 🌱 Lepretre Thomas
- 🐤 Perdereau Tom
- 🌟 Saby Loyola Sophie
- 👽 Sporck Trombini Gabriel
