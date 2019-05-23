#!/bin/sh

# exit immediately on failure
set -e

# This collection of commands builds on the set in edge-cloud/vault/setup.sh
# It configures MC access.


# these are commented out but are used to set the dme/mcorm secrets
#vault kv put jwtkeys/mcorm secret=12345 refresh=60m

# set mcorm approle
cat > /tmp/mcorm-pol.hcl <<EOF
path "auth/approle/login" {
  capabilities = [ "create", "read" ]
}

path "jwtkeys/data/mcorm" {
  capabilities = [ "read" ]
}

path "jwtkeys/metadata/mcorm" {
  capabilities = [ "read" ]
}

path "secret/data/accounts/sql" {
  capabilities = [ "read" ]
}

path "secret/data/accounts/noreplyemail" {
  capabilities = [ "read" ]
}
EOF
vault policy write mcorm /tmp/mcorm-pol.hcl
rm /tmp/mcorm-pol.hcl
vault write auth/approle/role/mcorm period="720h" policies="mcorm"
# get mcorm app roleID and generate secretID
vault read auth/approle/role/mcorm/role-id
vault write -f auth/approle/role/mcorm/secret-id

# set crm approle
cat > /tmp/crm-pol.hcl <<EOF
path "auth/approle/login" {
  capabilities = [ "create", "read" ]
}

path "secret/data/cloudlet/*" {
  capabilities = [ "read" ]
}
EOF
vault policy write crm /tmp/crm-pol.hcl
rm /tmp/crm-pol.hcl
vault write auth/approle/role/crm period="720h" policies="crm"
# get crm app roleID and generate secretID
vault read auth/approle/role/crm/role-id
vault write -f auth/approle/role/crm/secret-id

# set rotator approle - rotates dme/mcorm secret
cat > /tmp/rotator-pol.hcl <<EOF
path "auth/approle/login" {
  capabilities = [ "create", "read" ]
}

path "jwtkeys/data/*" {
  capabilities = [ "create", "update", "read" ]
}

path "jwtkeys/metadata/*" {
  capabilities = [ "read" ]
}
EOF
vault policy write rotator /tmp/rotator-pol.hcl
rm /tmp/rotator-pol.hcl
vault write auth/approle/role/rotator period="720h" policies="rotator"
# get rotator app roleID and generate secretID
vault read auth/approle/role/rotator/role-id
vault write -f auth/approle/role/rotator/secret-id

# mexenv approle
# This is used by mexdind to access registry secrets only.
# It does not have full access of crm role, so we can put it in
# the local_dind.yml setup file. Once we make our demo
# apps public, we can remove this app role.
cat > /tmp/mexenv-pol.hcl <<EOF
path "auth/approle/login" {
  capabilities = [ "create", "read" ]
}

path "secret/data/cloudlet/openstack/mexenv.json" {
  capabilities = [ "read" ]
}
EOF
vault policy write mexenv /tmp/mexenv-pol.hcl
rm /tmp/mexenv-pol.hcl
vault write auth/approle/role/mexenv period="720h" policies="mexenv"
# get mexenv app roleID and generate secretID
vault read auth/approle/role/mexenv/role-id
vault write -f auth/approle/role/mexenv/secret-id