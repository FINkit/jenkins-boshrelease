# BOSH release of Jenkins

[![Build Status](https://travis-ci.org/FINkit/jenkins-boshrelease.svg?branch=master)](https://travis-ci.org/FINkit/jenkins-boshrelease)

## 1. Overview

This is a BOSH release of a Jenkins Master and its Slaves, configured with a base set of plugins.

## 2. Release

### Create

```
bosh -e MY_ENV \
  create-release
```
### Upload

```
bosh -e MY_ENV \
  upload-release
```

## 3. Configuration
### Disabling installed plugins
Installed plugins may be disabled at deployment time by providing a list of plugins to be disabled in an operations file.
e.g. disable-plugins.yml
```
---
- type: replace
  path: /instance_groups/name=jenkins-master/jobs/name=jenkins-master/properties?/jenkins/disabled_plugins
  value: |
    gatling
    scoverage
    whitesource
```

## 4. Feature Toggles
The use of certain features can be controlled using the feature ops file at [/operations/feature-toggles.yml](operations/feature-toggles.yml).<br/>
e.g. draining Jenkins master
```
---
- type: replace
  path: /instance_groups/name=jenkins-master/jobs/name=jenkins-master/properties?/toggle/jenkins/drain
  value: false
```
Short lived toggles should eventually be enabled by default within the job spec or removed.
