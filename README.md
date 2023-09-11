[![github action status](https://github.com/hexlet-basics/hexlet-basics/workflows/push/badge.svg)](https://actions-badge.atrox.dev/hexlet-basics/hexlet-basics/goto)

# hexlet-basics

## Setup

### Requirements

* docker
* docker compose V2
* ruby >= 3.0.0
* make

### Steps

**Add to _/etc/hosts_:**
  127.0.0.1 code-basics.test

**Clone project**

Some lsp servers are fully workable only when the root dir is the same inside and outside the container. That is why we set WORKDIR to `/opt/projects/hexlet-basics`. So, if it is possible, clone this project to that directory.

**Run**

```bash
make project-setup
make compose # run server
# open code-basics.test

make app-test # run tests

# load language
# make app-language-load L=php

make app-db-prepare # sometimes, when fixtures were changed
```

To manage loaded languages and set other settings, you need to sign in (login: `full@test.io`, password: `password`)

### Production

Kube access

```bash
# make k8s-macos-setup or make k8s-ubuntu-setup
export DIGITALOCEAN_ACCESS_TOKEN=<your token>
make local-cluster-setup
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

## FAQ

#### I can't make commits with my changes

- try to reinitialize your fork repository.
[Read more](https://www.airplane.dev/blog/fixing-fatal-not-a-git-repository-error)
- try to execute `git commit` command with flag `--no-verify`.
_Example_: `git commit --no-verify -m '[#254] rename variable' `

---

[![Hexlet Ltd. logo](https://raw.githubusercontent.com/Hexlet/assets/master/images/hexlet_logo128.png)](https://hexlet.io/?utm_source=github&utm_medium=referral&utm_campaign=hexlet&utm_content=hexlet-basics)

This repository is created and maintained by the team and the community of Hexlet, an educational project. [Read more about Hexlet](https://hexlet.io/?utm_source=github&utm_medium=referral&utm_campaign=hexlet&utm_content=hexlet-basics).

See most active contributors on [hexlet-friends](https://friends.hexlet.io/).
