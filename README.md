[![github action status](https://github.com/hexlet-basics/hexlet-basics/workflows/docker/badge.svg)](https://actions-badge.atrox.dev/hexlet-basics/hexlet-basics/goto)

# hexlet-basics

## Участие

* Обсуждение в канале #hexlet-volunteers слака http://slack-ru.hexlet.io

## Setup

### Required

* docker
* docker compose
* ruby >= 3.0.0

### Steps

**Add to _/etc/hosts_:**
  127.0.0.1 code-basics.test ru.code-basics.test en.code-basics.test

```sh
$ make project-setup
$ make compose # run server
# open code-basics.test

$ make app-test # run tests

# load language
# make app-language-load L=php

$ make app-db-prepare # sometimes, when fixtures were changed
```

### Production
Kube access

```
make local-cluster-setup
```


##
[![Hexlet Ltd. logo](https://raw.githubusercontent.com/Hexlet/assets/master/images/hexlet_logo128.png)](https://ru.hexlet.io/pages/about?utm_source=github&utm_medium=link&utm_campaign=exercises-javascript)

This repository is created and maintained by the team and the community of Hexlet, an educational project. [Read more about Hexlet (in Russian)](https://ru.hexlet.io/pages/about?utm_source=github&utm_medium=link&utm_campaign=exercises-javascript).
##
