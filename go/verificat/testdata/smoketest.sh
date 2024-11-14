#!/usr/bin/env zsh
#
#   Smoketest for Verificat
#
#   Requires VPN access
#
#   Make sure your ENV VARs are set or in a local .env
#
[[ -n VCAT ]] && VCAT="http://localhost:4330"

print "\n::: Running smoketest for Verificat at $VCAT :::\n"

print -n "Source .env for EnvVars... "
set -a; source .env
gh_token=`env | grep GH_TOKEN= | sed 's/=.*/=<REDACTED>/'`
backstage=`env | grep BACKSTAGE=`
port=`env | grep PORT=`
print "loaded:\n $gh_token\n $backstage\n $port\n"

print -n "Healthz endpoint... "
health=$(curl -s $VCAT/healthz)
[[ $health == "ok" ]] && print $health || print "NOT OK"

print -n "Almanac download... "
almanac=$(curl -s $VCAT/almanac)
almsize=$(printf %s "$almanac" | wc -c)
print "$almsize bytes"

print -n "Admin service check... "
adminalive=$(curl -s -X POST $VCAT/v0/admin)
if [[ -n $adminalive ]]; then
  print $adminalive | jq
else
  print $?
  print "Sorry, no admin data was found.\n"
fi

print -n "Homepage copyright... "
colophon=$(curl -s $VCAT | grep Diesel)
print $colophon
