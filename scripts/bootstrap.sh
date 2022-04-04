#!/usr/bin/env bash

set -euo pipefail

: "${ENV:=dev}"
: "${RETAIN:=false}"

echo "beginning bootstrapping for \"${ENV}\" environment..." >&2

pulumi login "file://$(pwd)"

pushd bootstrap

export PULUMI_CONFIG_PASSPHRASE=''
if [ "$RETAIN" == "false" ]; then
  pulumi destroy --yes --stack "bootstrap-${ENV}" || true
  pulumi stack rm --force --preserve-config --yes --stack "bootstrap-${ENV}" || true
fi
if ! pulumi stack select "bootstrap-${ENV}"; then
  pulumi stack init --stack "bootstrap-${ENV}"
fi
pulumi up --yes --stack "bootstrap-${ENV}"

BUCKET="$(pulumi stack output --stack "bootstrap-${ENV}" 'bucket name')"
SECRETS_PROVIDER="$(pulumi stack output --stack "bootstrap-${ENV}" 'secrets provider')"

pulumi stack unselect
unset PULUMI_CONFIG_PASSPHRASE

popd


pulumi login "${BUCKET}"

for APPLICATION in "ingester" "hashtags" "visualise"; do

  DEPLOYMENT="inf"

  if [ ! -f "${APPLICATION}/${DEPLOYMENT}/Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml" ]; then
    echo "${APPLICATION}/${DEPLOYMENT}/Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml does not exist, skipping..." >&2
    continue
  fi

  pushd "${APPLICATION}/${DEPLOYMENT}"

  if pulumi stack select "${APPLICATION}-${DEPLOYMENT}-${ENV}"; then
    echo "\"${APPLICATION}-${DEPLOYMENT}-${ENV}\" exists, skipping..." >&2
    pulumi stack unselect
    popd
    continue
  fi

  echo "initialising \"${APPLICATION}-${DEPLOYMENT}-${ENV}\"..." >&2
  # using sed with backup for GNU & BSD compatability
  sed -i '.bak' '/^encryptedkey:/d' "Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml"
  rm "Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml.bak"
  pulumi stack init \
    --stack "${APPLICATION}-${DEPLOYMENT}-${ENV}" \
    --secrets-provider "${SECRETS_PROVIDER}"
  
  if [ ! -f "${ENV}.env" ]; then
    echo "no secrets configured for \"${APPLICATION}-${DEPLOYMENT}-${ENV}\"..." >&2
    pulumi stack unselect
    popd
    continue
  fi

  echo "configuring secrets for \"${APPLICATION}-${DEPLOYMENT}-${ENV}\"..." >&2
  while IFS='=' read KEY VALUE; do
    pulumi config set --secret "$KEY" "$VALUE"
  done < <(sed -e 's/^[[:space:]]*#.*// ; /^[[:space:]]*$/d' "${ENV}.env" )

  pulumi stack unselect
  popd

done

pulumi logout

echo "bootstrap complete." >&2
