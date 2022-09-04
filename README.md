# Helm Cleanup Action

[![CI](https://github.com/jonstacks/helm-cleanup/actions/workflows/ci.yml/badge.svg)](https://github.com/jonstacks/helm-cleanup/actions/workflows/ci.yml)

Rather than relying on the `helm` CLI, you can use the `helm-cleanup` action to cleanup a single chart or multiple charts matching specific filters.
My original intent is for this to be used with a scheduled GitHub Action to cleanup helm releases that are used for review apps after X amount of time.


## Usage

```yaml

steps:
- uses: actions/checkout@v3
- uses: azure/setup-kubectl@v1
- uses: azure/setup-helm@v1
  with:
    version: v3.8.0
- uses: jonstacks/helm-cleanup@v0
  with:
    ### 
    ### Filters. At least one filter is required to prevent deleting all releases
    ### accidentally. Filters are additive.
    ###
    last-modified-older-than: 168h # 7 days
    release-name-filter: 'review-app-.*'

    ### Optional ###

    debug: true                    # default: false
    description: some description  # default: nil
    dry-run: true                  # default: false
    keep-history: true             # default: false
    kube-context: my-context       # default: current kube context
    namespace: my-namespace        # default: current namespace
    no-hooks: true                 # default: false
    timeout: 10m                   # default: 5m
    wait: true                     # default: false
```
