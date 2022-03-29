#!/usr/bin/env bash

set -euo pipefail

: "${ENV:=dev}"

echo "Beginning bootstrapping for \"${ENV}\" environment..." >&2

pulumi login "file://$(pwd)"

pushd bootstrap
  export PULUMI_CONFIG_PASSPHRASE=''
  pulumi destroy --yes --stack "bootstrap-${ENV}" || true
  pulumi stack rm --force --preserve-config --yes --stack "bootstrap-${ENV}" || true
  pulumi stack init --stack "bootstrap-${ENV}"
  pulumi up --yes --stack "bootstrap-${ENV}"

  BUCKET="$(pulumi stack output --stack bootstrap-${ENV} 'bucket name')"
  SECRETS_PROVIDER="$(pulumi stack output --stack bootstrap-${ENV} 'secrets provider')"
  unset PULUMI_CONFIG_PASSPHRASE
popd

pulumi login "${BUCKET}"

for APPLICATION in "ingester" "hashtags" "visualise"; do
  DEPLOYMENT="bootstrap"
  if [ -d "./${APPLICATION}/${DEPLOYMENT}/" ]; then
    pushd "${APPLICATION}/${DEPLOYMENT}"
      # setting backup for GNU & BSD compatability
      sed -i '.bak' '/^encryptedkey:/d' "Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml"
      rm "Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml.bak"
      pulumi stack init --stack "${APPLICATION}-${DEPLOYMENT}-${ENV}" --secrets-provider "${SECRETS_PROVIDER}"
      pulumi up --yes --stack "${APPLICATION}-${DEPLOYMENT}-${ENV}"
    popd
  else
    echo "${APPLICATION}/${DEPLOYMENT} does not exist or is not a directory, skipping ${DEPLOYMENT}..." >&2
  fi
done

echo "bootstrap complete." >&2
