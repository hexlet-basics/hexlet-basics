---
version: 2
updates:
  - package-ecosystem: github-actions
    directory: /
    schedule:
      interval: weekly
    labels: [dependencies, automated, actions]
    groups:
      actions-minor-patch:
        update-types: [minor, patch]
      actions-major:
        update-types: [major]
  - package-ecosystem: bundler
    directory: /
    schedule:
      interval: weekly
    labels: [dependencies, automated, ruby]
    groups:
      bundler-minor-patch:
        update-types: [minor, patch]
      bundler-major:
        update-types: [major]
