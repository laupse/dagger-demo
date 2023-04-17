# DAGGER à DEVOXX FRANCE 2023

🎥 COMING SOON

Ce répertoire contient les sources relatives au Tool-in-Action présenté à Devoxx France 2023

## API PLAYGROUND

L'interface principal pour accéder à Dagger, c'est via l'API GraphQL.

Graçe au [playground](https://play.dagger.cloud/), il est possible de le faire sans installation

Pour une version deja remplie spécifique au sources de ce repertoire: https://play.dagger.cloud/playground/aGxZ9AwGRc4

## GO SDK

### Prérequis 

* 🐋 [Docker](https://www.docker.com/get-started/)
* 🔵 [Go Version >=1.20](https://go.dev/doc/install)
* 🧙 [Mage](https://magefile.org/)
    * `go install github.com/magefile/mage@300bbc868ba8f2c15b35e09df7e8804753cac00d`

### Usage

Avec les SDKs Dagger, on ne s'occupe pas de la partie GraphQL. Depuis le dossier racine, on peut lancer la commande `mage` pour voir l'ensemble des targets éxecutable (fonction GO) :

```
$ mage
Targets:
  all                Test, Build, Run
  build              Compile local go file into binary
  buildConcurrent    Build binary in multiple architecture concurrently
  image              Download image used in the engine
  run                Compile and run the binary
  secret             Try to leak a secret in dagger
  service            Run redis cmd against a redis server declared in Dagger
  test               Run unit test in math folder
```



