# Changelog

## [0.1.4](https://github.com/hexlet-basics/hexlet-basics/compare/hexlet_basics-v0.1.3...hexlet_basics-v0.1.4) (2026-07-24)


### Bug Fixes

* **app:** reload once on stale chunk import errors after deploy ([dcb801a](https://github.com/hexlet-basics/hexlet-basics/commit/dcb801afb5e548b3db75d494c9941aaf929276d3))
* **lessons:** quiet code-language Sentry noise and highlight lua ([fe101ee](https://github.com/hexlet-basics/hexlet-basics/commit/fe101eef7b5a189186830674b6fdf16f2dc5ec65))
* **lessons:** stop ShikiError on unmapped course code languages ([25964a9](https://github.com/hexlet-basics/hexlet-basics/commit/25964a9e74e11d5b1d948d8cddcb8a70aaa19d14))


### Miscellaneous

* update dependencies and sync sorbet RBIs ([3e0fc5f](https://github.com/hexlet-basics/hexlet-basics/commit/3e0fc5f1d68e31b97dc96fdcd92c19f215bb1a27))

## [0.1.3](https://github.com/hexlet-basics/hexlet-basics/compare/hexlet_basics-v0.1.2...hexlet_basics-v0.1.3) (2026-07-23)


### Features

* add admin-only Solid Queue monitor dashboard ([dc3293a](https://github.com/hexlet-basics/hexlet-basics/commit/dc3293a1c28ffa8cbe5512e6757719a5c843b5f0))
* **admin:** add flipper feature flags with an admin-gated UI ([fd703df](https://github.com/hexlet-basics/hexlet-basics/commit/fd703df6c89119b16bb0b644cb07f33634a536b4))
* **admin:** hide lesson reviews without student questions ([95a293d](https://github.com/hexlet-basics/hexlet-basics/commit/95a293de77382547a85f1c85f80d24709ec6933f)), closes [#619](https://github.com/hexlet-basics/hexlet-basics/issues/619)
* **lessons:** hint about editor autocomplete in the trainer ([48b6b88](https://github.com/hexlet-basics/hexlet-basics/commit/48b6b88a9fcbb278a87643b2ff27878258a72185))
* **lessons:** let users resize the theory and editor panes ([d8147a1](https://github.com/hexlet-basics/hexlet-basics/commit/d8147a1f339055a1e555f4959ef58df3f7bb7d34)), closes [#585](https://github.com/hexlet-basics/hexlet-basics/issues/585)
* **lessons:** screen-reader support and autocomplete hint in the editor ([3157c17](https://github.com/hexlet-basics/hexlet-basics/commit/3157c17674fc86e5424dd9cfaed26c917e2e7a23))
* **security:** throttle abusive traffic with rack-attack ([0285e79](https://github.com/hexlet-basics/hexlet-basics/commit/0285e79d26023291bcfb93a3f243ef4966d6d15f))
* **ui:** label the two Telegram links so they are distinguishable ([03af597](https://github.com/hexlet-basics/hexlet-basics/commit/03af597847f3092164558ce2482765b32ddd725e)), closes [#634](https://github.com/hexlet-basics/hexlet-basics/issues/634)
* **ui:** link the footer "Categories" heading to the categories page ([8d031b3](https://github.com/hexlet-basics/hexlet-basics/commit/8d031b3b1cbe42f033054a7bf83b26cfe9e64013)), closes [#637](https://github.com/hexlet-basics/hexlet-basics/issues/637)


### Bug Fixes

* **admin:** don't 500 when course form submits serialized cover attachment ([5d5cc2f](https://github.com/hexlet-basics/hexlet-basics/commit/5d5cc2fe05d0b5c7c56b264d2644af2ebb0f5f37)), closes [#740](https://github.com/hexlet-basics/hexlet-basics/issues/740)
* **auth:** return redirect instead of 500 on invalid password/passkey token ([91b293f](https://github.com/hexlet-basics/hexlet-basics/commit/91b293f9c2bc86554ff8e813cdd2c546ac3eaf95))
* **blog:** sync URL symmetrically while scrolling blog posts ([35534f8](https://github.com/hexlet-basics/hexlet-basics/commit/35534f89004b28d15929c34cd1dacaaa670c5b94)), closes [#587](https://github.com/hexlet-basics/hexlet-basics/issues/587)
* **build:** resolve monaco workers under Vite 8 / Rolldown ([1eea370](https://github.com/hexlet-basics/hexlet-basics/commit/1eea3707d875ca9c1f3baa44831adefd92c8a7ea))
* **ci:** return image_tag as String in non-production ([a45a465](https://github.com/hexlet-basics/hexlet-basics/commit/a45a46593b4cf5b8813f497ecaa1cd989bf1a2ca))
* **exercise-loader:** surface silent skip and docker failures in version result ([3fea8d1](https://github.com/hexlet-basics/hexlet-basics/commit/3fea8d16b7776b5296a8af9f8def849e1e5e74ae))
* **i18n:** add missing passkey error translations ([4a8a1cf](https://github.com/hexlet-basics/hexlet-basics/commit/4a8a1cfda7252dc3cd42704eb5da46384308680a))
* **i18n:** restore blanked frontend locale values ([a11ff41](https://github.com/hexlet-basics/hexlet-basics/commit/a11ff41d1be71c9be96ab66ff44051e5323b3d80))
* **i18n:** translate leftover Russian labels on the English course page ([c4dcd99](https://github.com/hexlet-basics/hexlet-basics/commit/c4dcd990e00bdf32738dcfd6b26501bb2235d38e)), closes [#571](https://github.com/hexlet-basics/hexlet-basics/issues/571)
* **languages:** add missing success.warning flash translation ([f020ae3](https://github.com/hexlet-basics/hexlet-basics/commit/f020ae3fe58f247f5ced98e4a25f51d5231d9f2b)), closes [#739](https://github.com/hexlet-basics/hexlet-basics/issues/739)
* **lessons:** drop copy button that overlapped code in trainer tabs ([1cd8444](https://github.com/hexlet-basics/hexlet-basics/commit/1cd84447d5ff4edf17d425c9323c177cf06294cf)), closes [#661](https://github.com/hexlet-basics/hexlet-basics/issues/661)
* **lessons:** keep text selectable in the solution tab ([08307e4](https://github.com/hexlet-basics/hexlet-basics/commit/08307e4de76d1a5b6a578cff75271f2a95c1d374)), closes [#652](https://github.com/hexlet-basics/hexlet-basics/issues/652)
* **lessons:** make code readable by screen readers in the editor ([78524f2](https://github.com/hexlet-basics/hexlet-basics/commit/78524f24484123395eab35570374a122b3480431))
* **lessons:** make navigation and assistant tabs reachable on mobile ([8bc5cfd](https://github.com/hexlet-basics/hexlet-basics/commit/8bc5cfd7f22878c0fe224bf3a7ee5cc3eb61ec8c)), closes [#720](https://github.com/hexlet-basics/hexlet-basics/issues/720)
* **lessons:** render code blocks without a language and unknown languages ([df232d7](https://github.com/hexlet-basics/hexlet-basics/commit/df232d7005807be8e01ef7ccfff5e83fc0c1ba3c)), closes [#592](https://github.com/hexlet-basics/hexlet-basics/issues/592) [#597](https://github.com/hexlet-basics/hexlet-basics/issues/597)
* **lessons:** render native list markers in lesson theory ([1e175db](https://github.com/hexlet-basics/hexlet-basics/commit/1e175db8193fc90eb8fb79faaeeab1132734f64f)), closes [#600](https://github.com/hexlet-basics/hexlet-basics/issues/600)
* **locales:** keep current page when switching language ([1936f5c](https://github.com/hexlet-basics/hexlet-basics/commit/1936f5cefbeae15342bd085dec8e65fc0acb4318)), closes [#648](https://github.com/hexlet-basics/hexlet-basics/issues/648)
* **locales:** use US flag for English instead of wrong UM flag ([6bb09c0](https://github.com/hexlet-basics/hexlet-basics/commit/6bb09c073da2ee0cf3a87e4c7146912ffc0fcdb7)), closes [#622](https://github.com/hexlet-basics/hexlet-basics/issues/622)
* **meta:** restore return value of display_escaped_meta_tags ([76b92ff](https://github.com/hexlet-basics/hexlet-basics/commit/76b92ff429b339f8e94fc1cd2782e26d08d1f7b1))
* **passkey:** replace deprecated @github/webauthn-json with native WebAuthn API ([7092892](https://github.com/hexlet-basics/hexlet-basics/commit/70928926a57b938cf378e6d5b1fafd276ab8f720))
* **seo:** make twitter:site handle locale-dependent ([357e347](https://github.com/hexlet-basics/hexlet-basics/commit/357e3474eb6be8e6c94071f065980f0c28800a1a)), closes [#628](https://github.com/hexlet-basics/hexlet-basics/issues/628)
* **ui:** mount Mantine Notifications so toasts actually render ([6a811a3](https://github.com/hexlet-basics/hexlet-basics/commit/6a811a38e849282523e89c31a09520a674f18a48))


### Miscellaneous

* **db:** add flipper tables to schema and generated RBIs ([ce51ef7](https://github.com/hexlet-basics/hexlet-basics/commit/ce51ef7a0b3c4266fa9a2a5f4df086e3bb649e0e))
* **db:** drop migrations superseded by schema, keep last five ([8dbe4a9](https://github.com/hexlet-basics/hexlet-basics/commit/8dbe4a92369fd15fb3d49de1c03df0e760e1c184))
* **deps:** bump yandex-cloud/yandex ([#736](https://github.com/hexlet-basics/hexlet-basics/issues/736)) ([c8d58a6](https://github.com/hexlet-basics/hexlet-basics/commit/c8d58a6fe58a9b92fcc7a688274359370ce616cb))
* **deps:** sync lockfiles and sorbet rbis, exclude rubocop plugin rbis ([c3e52b6](https://github.com/hexlet-basics/hexlet-basics/commit/c3e52b6628b964ec1399231495bc5281e5828636))
* **deps:** update dependencies and harden pnpm supply-chain config ([687fe80](https://github.com/hexlet-basics/hexlet-basics/commit/687fe8073f93e37d2185c7a0a1e6835f25c8d9cf))
* **deps:** update sorbet gem RBIs ([31251e9](https://github.com/hexlet-basics/hexlet-basics/commit/31251e946ef6b82a80a81eeeb2175eccea93e80a))
* drop accidental tapioca RBI churn ([b05d1ee](https://github.com/hexlet-basics/hexlet-basics/commit/b05d1ee4f5df1e77d18310b7fa7371472d2632f2))
* **i18n:** add i18n-tasks to audit backend locale files ([3ab9d9c](https://github.com/hexlet-basics/hexlet-basics/commit/3ab9d9c71079878989ef57d6ea29f080503310c6))
* **k8s:** raise job worker memory and set JOB_CONCURRENCY=1 ([af1942a](https://github.com/hexlet-basics/hexlet-basics/commit/af1942abd7d63a21c8442f682cf18458908522ae))
* **lint:** enable new cops and correctness departments, fix offenses ([edf259b](https://github.com/hexlet-basics/hexlet-basics/commit/edf259bc4e5a89c1a300e55094483455e7aa71ef))
* **logging:** add lograge for single-line structured request logs ([3409de9](https://github.com/hexlet-basics/hexlet-basics/commit/3409de901e775165dd51436becaf0dc4ad909aec))
* **makefile:** rename editor-setup to sync-sorbet and add to sync ([679a53f](https://github.com/hexlet-basics/hexlet-basics/commit/679a53fd323281d56effb54f9113298f885a9c00))
* remove dead duplicate AdminConstraint ([1c58181](https://github.com/hexlet-basics/hexlet-basics/commit/1c5818173f317228a700f4b29925bccba632e645))
* **sorbet:** add generated RBIs for flipper gems ([ca77853](https://github.com/hexlet-basics/hexlet-basics/commit/ca77853b3b05e74ffcaf8b5cfd2af21e06532c02))
* **sorbet:** enforce strict sigil via rubocop-sorbet and raise ActionCable ([eaeda46](https://github.com/hexlet-basics/hexlet-basics/commit/eaeda464ed8aa8d3d4040070672ab78847ca9d9b))
* sync generated RBIs and locales ([96f2d3f](https://github.com/hexlet-basics/hexlet-basics/commit/96f2d3f17cf04ec393fb90c4d4641d91b120677b))


### Build System

* upgrade Node to 26.5.0 ([dc44ca6](https://github.com/hexlet-basics/hexlet-basics/commit/dc44ca61be3c93adb915a6930c37bceb227e17f6))

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
