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
### Enabling persistent disks
Persistent disks may be enabled at deployment time by providing a list of such disks in an operations file, as shown in the example at [/operations/add-persistent-disks.yml](operations/add-persistent-disk.yml). Also, the required disk type must be specified in your cloud config, as shown in the example snippet below:
```
disk_types:
- disk_size: 1024
  name: default
- disk_size: 10_240
  name: 10GB
- disk_size: 100_240
  name: 100GB
```

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

## 5. Drain
Draining of both the Jenkins master and individual agents can be used to ensure running builds have completed before the nodes are shutdown during any restart.  A drain of Jenkins master will ensure all running builds have completed before shutdown, whereas a drain of an individual agent will only ensure running builds on that agent are completed.<br/>

Builds that take longer than the defined timeout period will be forcibly cancelled.  Any pending builds will resume when the node is active again.
