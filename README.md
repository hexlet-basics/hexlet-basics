[![github action status](https://github.com/hexlet-basics/hexlet-basics/workflows/docker/badge.svg)](https://actions-badge.atrox.dev/hexlet-basics/hexlet-basics/goto)

# hexlet-basics

## Setup

### Requirements

* docker
* docker compose
* ruby >= 3.0.0
* make

### Steps

**Add to _/etc/hosts_:**
  127.0.0.1 code-basics.test ru.code-basics.test

```bash
make project-setup
make compose # run server
# open code-basics.test

make app-test # run tests

# load language
# make app-language-load L=php

make app-db-prepare # sometimes, when fixtures were changed
```

### Production

Kube access

```bash
# make k8s-macos-setup or make k8s-ubuntu-setup
export DIGITALOCEAN_ACCESS_TOKEN=<your token>
make local-cluster-setup
```

---

[![Hexlet Ltd. logo](https://raw.githubusercontent.com/Hexlet/assets/master/images/hexlet_logo128.png)](https://hexlet.io/?utm_source=github&utm_medium=referral&utm_campaign=hexlet&utm_content=hexlet-basics)

This repository is created and maintained by the team and the community of Hexlet, an educational project. [Read more about Hexlet](https://hexlet.io/?utm_source=github&utm_medium=referral&utm_campaign=hexlet&utm_content=hexlet-basics).

See most active contributors on [hexlet-friends](https://friends.hexlet.io/).
