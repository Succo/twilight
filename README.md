## Serveur de jeux pour Vampire VS Werewolf

Ce serveur à vocation à remplacer l'officiel et à réimplémenter l'interface.

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
  -rows int
    	total number of rows (default 10)
```

La visualisation de la partie peut se faire sur navigateur à http://localhost:8080/.
Pour des raisons de gain de temps, vue.js a été utilisé pour le rendu réactif du front.

Le code pour les simulations est un copié collé de celui d'orgine re-écrit en go, il devrait donc être correcte.
