# Changelog

## [0.1.2](https://github.com/hexlet-basics/hexlet-basics/compare/hexlet_basics-v0.1.1...hexlet_basics-v0.1.2) (2026-07-16)


### Bug Fixes

* **api:** return 422 instead of 500 on lesson check without version_id ([2e0a4b8](https://github.com/hexlet-basics/hexlet-basics/commit/2e0a4b8270aeb6e89711e24895435e9d300512a5))
* **blog:** trim card bottom spacing and localize post date ([d1c1e61](https://github.com/hexlet-basics/hexlet-basics/commit/d1c1e6156cff02f4d2eef94ac9c7b3755c866303))
* **i18n:** load ru and es dayjs locales for date formatting ([e6637e2](https://github.com/hexlet-basics/hexlet-basics/commit/e6637e28e318ab732bff5af05ae5c84986214843))
* **languages:** match registration heading size to home page ([b494433](https://github.com/hexlet-basics/hexlet-basics/commit/b494433844997f8d4edb17008688754680976fe7))
* **shiki:** register diff language for code highlighting ([415ab1b](https://github.com/hexlet-basics/hexlet-basics/commit/415ab1baeea8d5c984969bff04e8798ff51dc931))
* **telegram:** point bot links to community and dedupe via translation ([4c4de41](https://github.com/hexlet-basics/hexlet-basics/commit/4c4de41794f9b1663eb33589ee2defd5d54639c7))


### Performance Improvements

* **web:** drop redundant landing-page includes that broke eager loading ([28459ee](https://github.com/hexlet-basics/hexlet-basics/commit/28459ee1d4bbd863ab9830e5c69270e6341a195b))


### Miscellaneous

* **deps:** add bundler checksum to Gemfile.lock ([830c3c3](https://github.com/hexlet-basics/hexlet-basics/commit/830c3c37ac92fadcbc77ff74f680943fd00e9fd2))
* **deps:** bump mantine-datatable to 9.4.0 and i18next-cli to 1.66.0 ([c00f2a4](https://github.com/hexlet-basics/hexlet-basics/commit/c00f2a4c8a8fe79aa44b1e4b9604bb3be01c2142))
* **k8s:** delete migrate/notify hook pods on success ([88738d0](https://github.com/hexlet-basics/hexlet-basics/commit/88738d08fd2cb4d7c9462dfa985bbf608c8c772b))

## [0.1.1](https://github.com/hexlet-basics/hexlet-basics/compare/hexlet_basics-v0.1.0...hexlet_basics-v0.1.1) (2026-07-15)


### Bug Fixes

* **k8s:** deploy web-check DaemonSet into release namespace ([e0b0bba](https://github.com/hexlet-basics/hexlet-basics/commit/e0b0bbad07ff6ea89d5f8d3919f68c1f8ddbfeff))


### Performance Improvements

* enable goldiloader in production ([c088b6b](https://github.com/hexlet-basics/hexlet-basics/commit/c088b6b0e3eb751d4d026e754d5950685cadc36f))


### Miscellaneous

* bump .ruby-version to 4.0.6 ([349129e](https://github.com/hexlet-basics/hexlet-basics/commit/349129ed70f22ef2a32d99fae78c8fc42fc4e51a))
* **dependabot:** add npm (pnpm), docker, and terraform ecosystems ([da3cd52](https://github.com/hexlet-basics/hexlet-basics/commit/da3cd525216bbc6e2d0451bc43df407104cedaf0))
* **deps:** bump gem and npm dependencies, refresh skills ([cd3e789](https://github.com/hexlet-basics/hexlet-basics/commit/cd3e789a107361b82c345094c135fc610106e6fb))
* **deps:** bump library/ruby in the docker-minor-patch group ([#730](https://github.com/hexlet-basics/hexlet-basics/issues/730)) ([65ff54a](https://github.com/hexlet-basics/hexlet-basics/commit/65ff54a339c6ac24f6387e41e6172720a8762167))
* **deps:** bump the actions-major group across 1 directory with 3 updates ([#727](https://github.com/hexlet-basics/hexlet-basics/issues/727)) ([842d347](https://github.com/hexlet-basics/hexlet-basics/commit/842d347c2dc690e0b107ae7e0910a04f94843dd2))
* **deps:** bump the terraform-minor-patch group ([#731](https://github.com/hexlet-basics/hexlet-basics/issues/731)) ([b4c6601](https://github.com/hexlet-basics/hexlet-basics/commit/b4c660181bc7a11f6bec7d97e70f92e105ff6128))
* **deps:** drop unused sqlite3, unpin tapioca ([3349aa9](https://github.com/hexlet-basics/hexlet-basics/commit/3349aa987a87a9b0b198062a103047d507d78bb2))
* **deps:** update sorbet rbi and synced artifacts ([2c0f7ad](https://github.com/hexlet-basics/hexlet-basics/commit/2c0f7ad60ae6694f0f840bb3731832b62f6747e2))
* **k8s:** add memory limit to nginx containers ([9487839](https://github.com/hexlet-basics/hexlet-basics/commit/9487839dbe45b8cc2e0a9213c3da00b967025ffa))
* **k8s:** drop values.yaml commit from helm-upgrade-app ([01fc2e5](https://github.com/hexlet-basics/hexlet-basics/commit/01fc2e5d15a7364c09edd6c5ed5390575c0a2a77))
* **k8s:** pin busybox, IfNotPresent pull policy, job TTL cleanup ([98dd049](https://github.com/hexlet-basics/hexlet-basics/commit/98dd0499e1ecb387263a019ed0db749db8b83a6d))
