FROM golang:1.15

# Note: Latest version of kubectl may be found at:
# https://aur.archlinux.org/packages/kubectl-bin/
ENV KUBE_LATEST_VERSION="v1.17.12"

RUN apt-get update && apt-get install -y --no-install-recommends \
		openssh-client \
		curl \
		bash \
  && wget -q https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/amd64/kubectl -O /usr/local/bin/kubectl \
  && chmod +x /usr/local/bin/kubectl \
	&& rm -rf /var/lib/apt/lists/*
