# Project_Go

# Descriptif du projet
Le but du projet est de faire un programme de recherche de chemin le plus court (ou presque) d'un point de départ à un point d'arrivé sur une carte aléatoire. On utilisera le fameux algorithme A* qui est très efficace pour résoudre ce genre de problème.
Dans un premier temps, on implémentera un programme qui puisse fonctionner pour une carte de type labyrinthe (qui soit faisable bien entendu) puis on généralisera pour qu'il puisse fonctionner sur tout type de carte.
Le but ensuite est que plusieurs clients puissent utiliser ce process sur la même carte avec chaqu'un un point de départ et d'arrivé différent, il faudra donc parralleliser le programme. Egalement pour optimiser le programme on pourra faire un A* qui part du départ et un autre qui part de l'arrivée en parrallèle pour augmenter l'éfficacité de la recherche.
 
# Descriptif du programme

# Serveur :

Répond à la demande d'un ou plusieurs clients à la fois qui veulent calculer un chemin d'un point A à un point B (différent) sur une carte en utilisant soit le A* normal ou le A* double.

A* double : un qui part du point de départ et un autre du point d'arrivé dès qu'il se rejoignent ils reconstituent le chemin. Sur des petites cartes ce n'est pas plus rapide c'est même plus lent.

# Client :

Envoie des coordonnées de départ (format : {x,y}) et d'arrivée et le type d'algo utilisé : "normal" ou "double"
