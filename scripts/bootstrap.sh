#!/usr/bin/env bash

set -euo pipefail

: "${ENV:=dev}"
: "${RETAIN:=true}"

echo "beginning bootstrapping for \"${ENV}\" environment..." >&2

export PULUMI_CONFIG_PASSPHRASE=''

APPLICATIONS=(
  "bootstrap"
  "shared-network"
  "shared-kafka"
  "ingester"
  "snowflake"
  "hashtags"
  "visualise"
)

BACKEND=''
SECRETS_PROVIDER=''

for APPLICATION in "${APPLICATIONS[@]}"; do

  while read -r PULUMI_PROJECT_FILE; do

    pulumi login "${BACKEND:-file://$(pwd)}"

    DEPLOYMENT="$(basename "$(dirname $PULUMI_PROJECT_FILE)" )"
    pushd "${APPLICATION}/${DEPLOYMENT}" 2> /dev/null

    STACK="${APPLICATION}-${DEPLOYMENT}-${ENV}"

    if ! pulumi stack select "${STACK}"; then
      echo "initialising \"${APPLICATION}-${DEPLOYMENT}-${ENV}\"..." >&2
      # using sed with backup for GNU & BSD compatability
      sed -i '.bak' '/^encryptedkey:/d' "Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml"
      rm "Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml.bak"
      pulumi stack init \
        --stack "${APPLICATION}-${DEPLOYMENT}-${ENV}" \
        --secrets-provider "${SECRETS_PROVIDER:-passphrase}"
    elif [ "$RETAIN" == "false" ]; then
      pulumi destroy --yes || true
      pulumi stack rm --preserve-config --force --yes || true
      # using sed with backup for GNU & BSD compatability
      sed -i '.bak' '/^encryptedkey:/d' "Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml"
      rm "Pulumi.${APPLICATION}-${DEPLOYMENT}-${ENV}.yaml.bak"
      pulumi stack init \
        --stack "${APPLICATION}-${DEPLOYMENT}-${ENV}" \
        --secrets-provider "${SECRETS_PROVIDER:-passphrase}"
    fi

    if [ -f "./${ENV}.env.sh" ]; then
      echo "configuring secrets for \"${APPLICATION}-${DEPLOYMENT}-${ENV}\"..." >&2
      ( source "${ENV}.env.sh" )
    fi

    if [ "${APPLICATION}" == "bootstrap" ]; then
      pulumi up --yes
      BACKEND="$(pulumi stack output --stack "${STACK}" 'pulumi backend' 2>/dev/null || true)"
      SECRETS_PROVIDER="$(pulumi stack output --stack "${STACK}" 'pulumi secrets provider' 2>/dev/null || true)"
      echo "${BACKEND} - $SECRETS_PROVIDER"
    fi 

    pulumi stack unselect
    popd 2>/dev/null

  done < <(find "$APPLICATION" -not -path '*/.*' -mindepth 2 -maxdepth 2 -name Pulumi.yaml)

done

pulumi logout

echo "bootstrap complete." >&2
