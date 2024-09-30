# {{ projectName }}
{% for shield in shields %}[![{{ shield.Description }}](https://img.shields.io/badge/{{ shield.Name }}-{{ shield.Value }}-success.svg?style=flat)](){% endfor %}
{% if pythonVersion != "" %}[![python version](https://img.shields.io/badge/Python-{{ pythonVersion }}-success.svg?style=flat)](){% endif %}
{% if djangoVersion != "" %}[![django version](https://img.shields.io/badge/Django-{{ djangoVersion }}-success.svg?style=flat)](){% endif %}
[![build status](https://img.shields.io/badge/build-pass-success.svg?style=flat)]()

## Technology stack

```bash
```

## Build

```bash
```

## Internationalization

```bash
```
