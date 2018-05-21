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
