# Alveus

Dynamically creates GitHub workflows to allow for progressive delivery of Kubernetes resources across environments.

This is designed as a replacement for [Kargo](https://docs.kargo.io/) and has some similarities to [ConcourseCI](https://concourse-ci.org/).

This initial version of Alveus uses GitHub actions as it's "execution platform".

Alveus focuses several concepts: 

- [Sources](docs/sources.md) - A source is an artifact, it's what flows down the river. It's what is promoted to each environment.

- [Rivers](docs/rivers.md) - A river is the series of GitHub workflows, the underlying execution platform

- [Drains](docs/drains.md) - A drain is the collection of steps performed when a source artifact is reaches it.

## Getting Started

```shell
alveus init
```

This will create a .alveus.yml file for you to fill in.

It will also create an initial GitHub workflow.

Create a PR with your changes to `.alveus.yml`. Alveus will then regenerate the workflow based on the contents of `.alveus.yml`.

You'll see everything that's going to happen when you deploy. This is the point of Alveus, to pull back the current on "deployment magic".

## Roadmap

- [ ] Abstract specifics into a "plugin"-like system 
    - [ ] ArgoCD
    - [ ] GitHub Workflows
