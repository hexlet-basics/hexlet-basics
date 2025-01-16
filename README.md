[![github action status](https://github.com/hexlet-basics/hexlet-basics/workflows/push/badge.svg)](https://actions-badge.atrox.dev/hexlet-basics/hexlet-basics/goto)

# hexlet-basics

## Setup

### Requirements

* docker
* ruby >= 3.0.0
* make

### Run

```bash
make setup
make dev # run server
# open code-basics.test

make test # run tests

# load language
# make language-load L=php

make sync # sometimes, when fixtures were changed
```

To manage loaded languages and set other settings, you need to sign in (login: `full@test.io`, password: `password`)

### Production

Kube access

```bash
# make k8s-macos-setup or make k8s-ubuntu-setup
export TWC_TOKEN=<your token>
```

### Deploy

* Create new tag via command:

  ```bash
  make next-tag
  ```

* Wait notification about ready tag in Slack channel `#sideprojects-code-auto` or wait [Github Actions](https://github.com/hexlet-basics/hexlet-basics/actions/workflows/release.yml)
* Change version in [k8s/hb-app-chart/values.yaml](/k8s/hb-app-chart/values.yaml) and then:

  ```bash
  make -C k8s helm-upgrade-app
  ```

## TODO

1. theme switcher
1. https://github.com/DavidWells/analytics
1. auth: vk id, yandex id
1. fix XPaging

---

[![Hexlet Ltd. logo](https://raw.githubusercontent.com/Hexlet/assets/master/images/hexlet_logo128.png)](https://hexlet.io/?utm_source=github&utm_medium=referral&utm_campaign=hexlet&utm_content=hexlet-basics)

This repository is created and maintained by the team and the community of Hexlet, an educational project. [Read more about Hexlet](https://hexlet.io/?utm_source=github&utm_medium=referral&utm_campaign=hexlet&utm_content=hexlet-basics).

See most active contributors on [hexlet-friends](https://friends.hexlet.io/).
