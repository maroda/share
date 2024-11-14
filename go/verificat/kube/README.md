# Kubernetes Resources

## Verificat Manifests

Order of operations for v0.0.1 (subject to change), use `kubectl apply -f`:

1. `verificat-namespace.yaml`
2. `verificat-app.yaml` with changes:
   1. `<BASE64_JSON_CONFIG_WITH_SECRET>` contains the auth with PAT
   2. `<GITHUB_TOKEN>` is only the PAT
   3. These will eventually be retrievable via ***ESO + Kube Secret***, see example below.
3. `verificat-ingress.yaml` is for using a local Traefik load balancer.

## External Secrets Operator with Localstack

### About

In order to use AWS Secret Manager (ASM) values inside Kubernetes, there must be middleware. External Secrets Operator (ESO) is an option that can use an arbitrary external store to fill Kubernetes Secrets.

This guide is for how to set up ESO locally using Localstack to provide ASM (and the IAM to get there) and Orbstack's kubernetes engine.

### Config Walkthrough

This set of configurations follow the [Getting Started walkthrough](https://external-secrets.io/latest/introduction/getting-started/) for ESO.

1. Run Orbstack: `orb start` and then Kubernetes: `orb start k8s`
2. Configure Orbstack to bridge container IP addresses to macOS (required for using Localstack DNS):

   1. This has the potential of breaking the local AWS VPN if it tries to reconnect because it can't operate in parallel with Orbstack's bridge interfaces.
   2. In **OrbStack Desktop** go to **Settings > Network**
   3. Check on **"Allow access to container domains & IPs"**
3. Run Localstack: `localstack start -d --network ls`
4. Get the DNS for Localstack's container (see [Access via endpoint URL](https://docs.localstack.cloud/references/network-troubleshooting/endpoint-url/#from-your-container))

   1. `docker inspect localstack-main | jq -r '.[0].NetworkSettings.Networks | to_entries | .[].value.IPAddress'`
5. Set `dnsConfig.nameservers` in `values.yaml` with the Localstack container IP
6. Create AWS IAM Principal

   1. This requires having a `localstack` profile in `~/.aws/config`
      1. Check STS: `aws --profile localstack sts get-caller-identity`
      2. The ARN should be: `arn:aws:iam::000000000000:root`
   2. Add a user: `aws --profile localstack iam create-user --user-name <USER>`
   3. Add an API Key: `aws --profile localstack iam create-access-key --user-name <USER>`
7. Create a secret for the Localstack AWS credentials (used by ESO to connect to ASM):

   1. `echo '<ACCESSKEY>' > ./lstack-access-key`
   2. `echo '<SECRETKEY>' > ./lstack-secret-access-key`
   3. Make sure you're using Orbstack: `kubectx orbstack`
   4. `kubectl create secret generic awssm-secret --from-file=./lstack-access-key --from-file=./lstack-secret-access-key`
8. Install External Secret Operator using Helm

   1. `helm install external-secrets external-secrets/external-secrets -n external-secrets --create-namespace --values ./values.yaml`
   2. Install the Store: `kubectl apply -f basic-secret-store.yaml`
   3. Define the External Secret: `kubectl apply -f basic-external-secret.yaml`
9. Inspect the new external secret: `kubectl describe externalsecret example`

To Uninstall:

1. Get all external secrets: `kubectl get SecretStores,ClusterSecretStores,ExternalSecrets --all-namespaces`
2. Remove all of them with `kubectl delete`
3. `helm delete external-secrets --namespace external-secrets`
