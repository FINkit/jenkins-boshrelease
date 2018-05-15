#! /bin/bash -e

JENKINS_INPUT_JOB_LIST=plugins.txt
if [ ! -f "$JENKINS_INPUT_JOB_LIST" ]
then
    echo "ERROR File not found: $JENKINS_INPUT_JOB_LIST"
    exit 1
fi

REF=temp
mkdir -p $REF
rm -f $REF/*
COUNT_PLUGINS_INSTALLED=0
JENKINS_UC=https://updates.jenkins-ci.org
while read -r spec || [ -n "$spec" ]; do

    plugin=(${spec//:/ });
    [[ ${plugin[0]} =~ ^# ]] && continue
    [[ ${plugin[0]} =~ ^[[:space:]]*$ ]] && continue
    [[ -z ${plugin[1]} ]] && plugin[1]="latest"

    if [ -z "$JENKINS_UC_DOWNLOAD" ]; then
      JENKINS_UC_DOWNLOAD=$JENKINS_UC/download
    fi

    echo "Downloading ${plugin[0]}:${plugin[1]}"
    curl --connect-timeout 60 --retry 5 --retry-delay 10 -sSL -f "${JENKINS_UC_DOWNLOAD}/plugins/${plugin[0]}/${plugin[1]}/${plugin[0]}.hpi" -o "$REF/${plugin[0]}.jpi"
    unzip -qqt "$REF/${plugin[0]}.jpi"
    (( COUNT_PLUGINS_INSTALLED += 1 ))
done  < "$JENKINS_INPUT_JOB_LIST"

echo "---------------------------------------------------"
if (( "$COUNT_PLUGINS_INSTALLED" > 0 ))
then
    echo "INFO: Successfully installed $COUNT_PLUGINS_INSTALLED plugins."
else
    echo "INFO: No changes, all plugins previously installed."
fi
echo "---------------------------------------------------"

cd temp
DATE=`date +%Y%m%d%H%M`
PLUGIN_ARCHIVE=plugins-${DATE}.zip
OUTPUT_DIR=../../../blobs/jenkins

mkdir -p ${OUTPUT_DIR}

jar cvf ${OUTPUT_DIR}/${PLUGIN_ARCHIVE} *

cd ..
rm -fr $REF

echo "cd ../.."
echo "bosh2 add-blob blobs/jenkins/${PLUGIN_ARCHIVE} jenkins/${PLUGIN_ARCHIVE}"
echo "bosh2 upload-blobs"

