## Serveur de jeux pour Vampires VS Werewolves

Ce serveur a vocation à pouvoir servir à la place de l'officiel, et vise donc une compatibilité maximum du point de vue des IA (joueurs).
Par contre, il est plutôt implémenté dans un optique de débug : il essaie donc d'offrir le plus d'infos pertinentes aux utilisateurs, et n'implémente pas strictement les règles en cas de mal fonctionnement de l'IA (type non réponse). 

```
Usage of twilight:
  -columns int
    	total number of columns (default 10)
  -humans int
    	quantity of humans group (default 16)
  -map string
    	path to the map to load (or save if randomly generating)
  -monster int
    	quantity of monster in the start case (default 8)
  -rand
    	use a randomly generated map
  -seed int
      use seed to generate same map
  -rows int
    	total number of rows (default 10)
```

Comme avec le serveur officiel, les bot se connectent sur le port 9000.
La visualisation de la partie peut se faire sur navigateur à http://localhost:8080/.
Pour des raisons de gain de temps, vue.js a été utilisé pour le rendu réactif du front.

Le code pour les simulations est un copié collé de celui d'origine, re-écrit en go, il devrait donc être correct.
