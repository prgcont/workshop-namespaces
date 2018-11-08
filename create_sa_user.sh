#!/bin/bash
set -o pipefail

printHelp() {
  echo "Script for adding namespace and user (service account) for workshop"
  echo "Usage:"
  echo "${0} [create|delete] STRING"
  echo "e.g. ${0} create mynamespace"
}

create_target_folder() {
    echo -n "Creating target directory to hold files in ${TARGET_FOLDER}..."
    mkdir -p "${TARGET_FOLDER}"
    printf "done"
}

create_namespace() {
    echo -e "\\nCreating namespace: ${NAMESPACE}"
    kubectl create namespace "${NAMESPACE}"
}

delete_namespace() {
  echo -e "\\nCreating namespace: ${NAMESPACE}"
  kubectl delete namespace "${NAMESPACE}"
}

delete_rolebinds() {
  kubectl -n ${NAMESPACE} delete clusterrolebinding kubeproxier-${SERVICE_ACCOUNT_NAME}
  kubectl -n ${NAMESPACE} delete clusterrolebinding crd-admin-${SERVICE_ACCOUNT_NAME}
}

create_service_account() {
  echo -e "\\nCreating a service account: ${SERVICE_ACCOUNT_NAME} on namespace: ${NAMESPACE}"
  kubectl create sa "${SERVICE_ACCOUNT_NAME}" --namespace "${NAMESPACE}"
}

get_secret_name_from_service_account() {
    echo -e "\\nGetting secret of service account ${SERVICE_ACCOUNT_NAME} from NAMESPACE ${NAMESPACE}"
    SECRET_NAME=$(kubectl get sa "${SERVICE_ACCOUNT_NAME}" --namespace "${NAMESPACE}" -o json | jq -r '.secrets[].name')
    echo "Secret name: ${SECRET_NAME}"
}

extract_ca_crt_from_secret() {
    echo -e -n "\\nExtracting ca.crt from secret..."
    kubectl get secret "${SECRET_NAME}" --namespace "${NAMESPACE}" -o json | jq \
    -r '.data["ca.crt"]' | base64 --decode > "${TARGET_FOLDER}/ca.crt"
    printf "done"
}

get_user_token_from_secret() {
    echo -e -n "\\nGetting user token from secret..."
    USER_TOKEN=$(kubectl get secret "${SECRET_NAME}" \
    --namespace "${NAMESPACE}" -o json | jq -r '.data["token"]' | base64 --decode)
    printf "done"
}

add_namespace_rbac() {
  kubectl -n ${NAMESPACE} create rolebinding --clusterrole=cluster-admin --serviceaccount=${NAMESPACE}:${SERVICE_ACCOUNT_NAME} ${SERVICE_ACCOUNT_NAME}-admin
  kubectl -n ${NAMESPACE} create clusterrolebinding --clusterrole=system:node-proxier --serviceaccount=${NAMESPACE}:${SERVICE_ACCOUNT_NAME} kubeproxier-${SERVICE_ACCOUNT_NAME}
  kubectl -n ${NAMESPACE} create clusterrolebinding --clusterrole=cluster-admin --serviceaccount=${NAMESPACE}:${SERVICE_ACCOUNT_NAME} crd-admin-${SERVICE_ACCOUNT_NAME}
}

set_kube_config_values() {
  context=${CONTEXT:-$(kubectl config current-context)}
  echo -e "\\nSetting current context to: $context"

  CLUSTER_NAME=${CLUSTER_NAME:-$(kubectl config get-contexts "$context" | awk '{print $3}' | tail -n 1)}
  echo "Cluster name: ${CLUSTER_NAME}"

  ENDPOINT=${ENDPOINT:-$(kubectl config view \
                                 -o jsonpath="{.clusters[?(@.name == \"${CLUSTER_NAME}\")].cluster.server}")}
  echo "Endpoint: ${ENDPOINT}"

  # Set up the config
  echo -e "\\nPreparing k8s-${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-conf"
  echo -n "Setting a cluster entry in kubeconfig..."
          kubectl config set-cluster "${CLUSTER_NAME}" \
          --kubeconfig="${KUBECFG_FILE_NAME}" \
          --server="${ENDPOINT}" \
          --certificate-authority="${TARGET_FOLDER}/ca.crt" \
          --embed-certs=true

  echo -n "Setting token credentials entry in kubeconfig..."
  kubectl config set-credentials \
          "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
          --kubeconfig="${KUBECFG_FILE_NAME}" \
          --token="${USER_TOKEN}"

  echo -n "Setting a context entry in kubeconfig..."
  kubectl config set-context \
          "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
          --kubeconfig="${KUBECFG_FILE_NAME}" \
          --cluster="${CLUSTER_NAME}" \
          --user="${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
          --namespace="${NAMESPACE}"

  echo -n "Setting the current-context in the kubeconfig file..."
  kubectl config use-context "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
          --kubeconfig="${KUBECFG_FILE_NAME}"
}

main() {
  if [ "$#" -lt 2 ]; then
    echo "Illegal number of parameters: Expected at 2 as minimum, got $#"
    printHelp
    exit 1
  fi

  ACTION="$1"
  NAMESPACE="$2"
  SERVICE_ACCOUNT_NAME=${NAMESPACE}
  KUBECFG_DIR="${3:-~/tmp/kube}"
  KUBECFG_FILE_NAME=${KUBECFG_DIR}/k8s-config-${NAMESPACE}
  TARGET_FOLDER=~/tmp/kube

  if [ ${ACTION} = 'create' ]; then
    create_target_folder
    create_namespace
    create_service_account
    get_secret_name_from_service_account
    extract_ca_crt_from_secret
    get_user_token_from_secret
    add_namespace_rbac
    set_kube_config_values

    echo -e "\\nAll done! Test with:"
    echo "KUBECONFIG=${KUBECFG_FILE_NAME} kubectl get pods"
  elif [ ${ACTION} = 'delete' ]; then
    # DELETE NAMESPACE
    delete_namespace
    delete_rolebinds
  else
    echo "Unknown action"
    printHelp
    exit 1
  fi
}

main "$@"
