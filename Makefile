# Set the shell to be /bin/bash and enable pipefail, which causes a pipeline to return a failure if any command fails
SHELL := /bin/bash -o pipefail

# Use kubectl with a specific Kubernetes context for interacting with the k3d cluster
KUBECTL := kubectl --context k3d-cluster

# .PHONY targets are not associated with files, they are always executed when called
.PHONY: create-k3d-cluster
.PHONY: delete-local-kube-cluster
.PHONY: build-pinger
.PHONY: run-local-kube-with-ping-pong-app

# Target: create-k3d-cluster
# This creates a new k3d Kubernetes cluster, but only after deleting any existing cluster
#
# Create the k3d cluster with 2 agent nodes and disable the default service load balancer (servicelb) and Traefik Ingress controller

create-k3d-cluster: delete-local-kube-cluster
	# Check if k3d is installed; if not, provide instructions on how to install it
	@which k3d >> /dev/null || echo "K3d must be installed to create local kube cluster\n==> wget -q -O - https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash" \
	&& k3d cluster create cluster --k3s-arg '--disable=servicelb@server:0' --k3s-arg '--disable=traefik@server:0' --agents 2

# Target: delete-local-kube-cluster
# Deletes any existing k3d cluster named 'cluster'
delete-local-kube-cluster:
	# Display a message and delete the cluster if it exists
	@echo "Deleting existing cluster..." && k3d cluster delete cluster

# Target: build-pinger
# Builds the Docker image for the 'pinger' service
build-pinger:
	# Build the Docker image for the pinger service from the Dockerfile in the app/pinger directory and tag it as ping:latest
	docker build -t pinger:latest app -f app/pinger/Dockerfile

# Target: build-ponger
# Builds the Docker image for the 'ponger' service
build-ponger:
	# Build the Docker image for the ponger service from the Dockerfile in the app/ponger directory and tag it as pong:latest
	docker build -t ponger:latest app -f app/ponger/Dockerfile

# Target: run-local-kube-with-ping-pong-app
# This builds both pinger and ponger, creates a local k3d Kubernetes cluster, and deploys the services
run-local-kube-with-ping-pong-app: build-pinger build-ponger create-k3d-cluster
	# Import the Docker images into the k3d cluster
	k3d image import pinger:latest --cluster cluster \
	&& k3d image import ponger:latest --cluster cluster \
	&& ${KUBECTL} create \
	  -f app/ponger/manifests \
	  -f app/pinger/manifests \
	&& echo "cluster available on kubernetes context k3d-cluster"


