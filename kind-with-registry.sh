#!/bin/sh
set -o errexit

# run an arbitrary core to create the buildkit container
cloak -p cloak.yaml do <<'EOF'
{
    core {
      image(ref: "index.docker.io/alpine") {
        exec(input: { args: ["echo", "start"] }) {
          fs {
            id
          }
        }
      }
    }
  }
EOF


# create registry container unless it already exists
reg_name='kind-registry'
reg_port='5001'
if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always --name "${reg_name}" --network=ci registry:2
    # -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
fi

if [ "$(docker inspect -f '{{.State.Running}}' "kubegraph-control-plane" 2>/dev/null || true)" != 'true' ]; then
  # create a cluster with the local registry enabled in containerd
  cat <<EOF |  kind create cluster --name kubegraph --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    # port forward 80 on the host to 80 on this node
    extraPortMappings:
      - containerPort: 30000
        hostPort: 30000
containerdConfigPatches:
  - |-
    [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
      endpoint = ["http://${reg_name}:5000"]
EOF

fi

kind export kubeconfig --name kubegraph --kubeconfig kind-ci.yaml


# connect the registry to the cluster network if not already connected
if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "dagger-buildkitd")" = 'null' ]; then
  docker network connect kind dagger-buildkitd
fi

# connect the registry to the cluster network if not already connected
if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${reg_name}")" = 'null' ]; then
  docker network connect kind "${reg_name}"
fi


# Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF


