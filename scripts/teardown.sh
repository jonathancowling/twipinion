#!/usr/bin/env bash

set -euo pipefail

: "${ENV:=dev}"

echo "Beginning teardown for \"${ENV}\" environment..." >&2

pulumi login "file://$(pwd)"

pushd bootstrap
  export PULUMI_CONFIG_PASSPHRASE=''
  BUCKET="$(pulumi stack output --stack bootstrap-${ENV} 'bucket name')"
  unset PULUMI_CONFIG_PASSPHRASE
popd

pulumi login "${BUCKET}"

for APPLICATION in "ingester" "hashtags" "visualise"; do
  for DEPLOYMENT in "inf" "bootstrap"; do
    if [ -d "./${APPLICATION}/${DEPLOYMENT}/" ]; then
        pushd "${APPLICATION}/${DEPLOYMENT}"
          if pulumi stack select "${APPLICATION}-${DEPLOYMENT}-${ENV}"; then
            echo "tearing down \"${APPLICATION}-${DEPLOYMENT}-${ENV}\" (${APPLICATION}/${DEPLOYMENT})..." >&2
            pulumi destroy --yes --stack "${APPLICATION}-${DEPLOYMENT}-${ENV}"
            pulumi stack rm --yes --preserve-config --stack "${APPLICATION}-${DEPLOYMENT}-${ENV}"
            # setting backup for GNU & BSD compatability
            sed -i '.bak' '/^encryptedkey:/d' "Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml"
            rm "Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml.bak"
          else
            echo "stack \"${APPLICATION}-${DEPLOYMENT}-${ENV}\" does not exist, skipping" >&2
          fi
        popd
    else
        echo "\"${APPLICATION}/${DEPLOYMENT}\" does not exist or is not a directory, skipping ${DEPLOYMENT}..." >&2
    fi
  done
done

pulumi login "file://$(pwd)"

pushd bootstrap
  export PULUMI_CONFIG_PASSPHRASE=''
  pulumi destroy --yes --stack "bootstrap-${ENV}"
  pulumi stack rm --preserve-config --yes --stack "bootstrap-${ENV}"
  unset PULUMI_CONFIG_PASSPHRASE
popd

echo "teardown complete." >&2
