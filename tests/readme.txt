les datasets trouvabels dans le dossier test servent à montrer
le comportement des agent en fonctions de différents paramètres

il sont tous constitué de 33 agents de chaque type, car il semble
qu'il n'y ait pas de configuration de type qui emprèche la convergence de l'opinion vers une certaine valeur
pour des paramètres donnés.
cela permet juste de ralentir ou d'accélérer le temps avant d'atteindre la zone de convergence.

Cette zone de convergence semble de manière générale favoriser les agents sceptiques,
Meme une simulation avec uniquement des croyants au départ se terminera avec une majorté de sceptiques

Aussi les simulations de test ont été effectuée sur une durée de 1 minute par défaut car ce délai
semble amplement suffisant pour observer une convergence de l'opinion. Il se peut qu'un test plus long soit 
effectué pour certains datasets si aucune convergence n'est observée, cela sera précisé dans la description du dataset.

les graphique de l'évolution de l'opinion pour chaque dataset peuvent etre retouvé dans l'onglet results
sous la forme "nomdudataset-paramètresdedéplacements.png"

ex : agents_0-123 : dataset agents_0 avec les paramètres de déplacement 1,2,3, cad patrol pour les croyants, heatmap pour les sceptiques, center of mass pour les neutres 

======DATASETS======



agents_0.json:
    - répartition uniforme du charisme
    - répartition aléatoires des relations

agents_1.json:
    - répartition normale(0.5,0.3) du charisme
    - répartition aléatoires des relations


agents_2.json:
    - le charisme est à 0 pour tous les agents
    - les relations sont aléatoires

agents_3.json:
    - le charisme est à 0 pour tous les agents
    - relation ennemi pour tout le monde

