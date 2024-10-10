# {{ projectName }}

{% for shield in shields %}[![{{ shield.Description }}](https://img.shields.io/badge/{{ shield.Name }}-{{ shield.Value }}-success.svg?style=flat)](){% if not forloop.Last %}
{% endif %}{% endfor %}

## Install

```bash

```

## Run

```bash
```

## Build

```bash
```
