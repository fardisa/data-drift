</br>
<p align="center">
  <a href="https://www.data-drift.io">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="./datadrift-logo-light.png" width="200px">
      <source media="(prefers-color-scheme: light)" srcset="./datadrift-logo-dark.png" width="200px">
      <img src="./datadrift-logo-dark.png" width="200px" alt="Datadrift logo" />
    </picture>
  </a>
</p>

<p align="center">
  <a href="https://discord.gg/GNEyCsrEve"><img src="https://dcbadge.vercel.app/api/server/GNEyCsrEve?style=flat-square&theme=discord" alt="Discord"></a>
  <a href="https://github.com/data-drift/data-drift/stargazers"><img src="https://img.shields.io/github/stars/data-drift/data-drift?style=flat-square" alt="Github Stars"></a>
  <a href="https://github.com/data-drift/data-drift/actions/workflows/datadrift-build.yml"><img src="https://img.shields.io/github/actions/workflow/status/data-drift/data-drift/datadrift-build.yml?style=flat-square" alt="Data-Drift Build"></a>
  <a href="https://main--64be84b7fe2172aa386216b8.chromatic.com/?path=/story/drift-dualtable--simple-case"><img src="https://img.shields.io/badge/storybook-visit-FF4785.svg?style=flat-square&logo=storybook" alt="Storybook"></a>
  <a href="https://pypi.org/project/driftdb/"><img src="https://img.shields.io/pypi/v/driftdb?style=flat-square" alt="DataGit version"></a>
</p>

<h1 align="center" >Metrics Observability & Troubleshooting</h1>

<p align="center">Datadrift is an open-source metric observability framework that helps data teams deliver trusted and reliable metrics.
</p>

<h3 align="center">
  <a href="https://www.data-drift.io/">Website</a> |
  <a href="https://github.com/data-drift/data-drift/issues">Issues</a> |
  <a href="https://www.data-drift.io/blog">Blog</a> |
  <a href="https://pypi.org/project/driftdb/">Doc</a> |
  <a href="https://github.com/orgs/data-drift/projects/3">Roadmap</a> |
  <a href="https://discord.gg/GNEyCsrEve">Discord</a>
</h3>

<p align="center">
  <a href="https://www.data-drift.io">
    <img src="./datadrift-overview.png" alt="DataDrift " />
  </a>
</p>

Data monitoring tools fail by focusing on static tests (eg. null, unique, expected values) and metadata monitoring (eg. column-level).
</br>

Datadrift monitors your metrics, sends alerts when anomalies are detected and automates root cause analysis. 
</br>
**Data teams detect and solve data issues faster with Datadrift's row-level monitoring & troubleshooting.**

</br>

# 🚀 Quickstart

## dbt installation

```
pip install driftdb
```

[Here is a quick demo](https://app.claap.io/sammyt/drift-db-demo-a18-c-ApwBh9kt4p-07oQMdsIzt_e).
For a step-by-step guide on the dbt installation, see the [docs](https://pypi.org/project/driftdb/).

## Python installation

Install the monitor in your pipeline.

```python
>>> from driftdb.connectors import LocalConnector
>>> LocalConnector().snapshot_table(table_dataframe=dataframe, table_name="revenue")
```

For a step-by-step guide on the python installation, see the [docs](https://pypi.org/project/driftdb/).

## Datadrift cloud

We are in development and we would love to do the installation with you. [Fill the form on our website](https://www.data-drift.io/) so we can do a 15min demo. If the tool solves your problem then the installation requires 30min.

</br>

# ⚡️ Key Features

## 🔮 Metrics monitoring & custom alerting

Get full visibility into metrics variation and pro-actively detect data quality issues. Become aware of unknown unknowns with metric drift custom alerting.

  <a href="https://www.data-drift.io">
    <img src="./datadrift-new-drift-alert.png" alt="DataDrift new drift custom alerting" width="800px"/>
  </a>

</br>

## 🧑‍🎤 Automated root cause analysis & troubleshooting

Operationalize your monitoring and solve your underlying data quality issue with lineage drill-down to understand the root cause of the problem.

  <a href="https://www.data-drift.io">
    <img src="./datadrift-metric-troubleshooting.png" alt="DataDrift diff compare table" width="800px"/>
  </a>

</br>

## 💎 Shared understanding of metric variation

Give visibility to data analysts and data consumers with shared explanation of metric variation.

  <a href="https://www.data-drift.io">
    <img src="./datadrift-changelog-dark.png" alt="DataDrift metric drift changelog" width="800px"/>
  </a>

</br>

## 🧠 And much more

We are in the early days of Datadrift. Just open a new [issue](https://github.com/data-drift/data-drift/issues) to tell us more about it and see how we could help!

</br>

# 💚 Community

We 💚 contributions big and small. In priority order (although everything is appreciated) with the most helpful first:

- [Star this repo](https://github.com/data-drift/data-drift) to help us get visibility and build awesome open-source tools
- [Join our Discord server](https://discord.gg/X2RUXFAm) to be part of our thriving community
- [Open an issue](https://github.com/data-drift/data-drift/issues) to share your idea or a bug you might have spotted
- [Become a Design Partner](https://www.data-drift.io/design-partner) to co-built a product you & users love

</br>

# 🗓 Upcoming features

Track planning on [Github Projects](https://github.com/orgs/data-drift/projects/3) and help us prioritising by upvoting or creating [issues](https://github.com/data-drift/data-drift/issues).
