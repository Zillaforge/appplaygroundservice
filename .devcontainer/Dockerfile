ARG GO_VER=1.22.4
ARG GO_ARCH=alpine
ARG INSTALL_KUBECTL=true
ARG KUBECTL_VER=1.29.1
ARG INSTALL_KUBEBUILDER=true
ARG KUBEBUILDER_VER=3.10.0

FROM golang:${GO_VER}-${GO_ARCH}

# https://blog.scottchayaa.com/post/2018/11/04/docker-arg-env-variable/#arg
# 1. ARG need to be declared again after FROM
# 2. To get benfit of ARG, you need to put ARG as close as to where it is used.
#    If you change value when building, it will invalidate cache and rebuild.


# Install Golang extension package tools
RUN go install github.com/ramya-rao-a/go-outline@v0.0.0-20210608161538-9736a4bde949    && \
    go install github.com/cweill/gotests/gotests@v1.7.0                                && \
    go install github.com/fatih/gomodifytags@v1.17.0                                   && \
    go install github.com/josharian/impl@v1.4.0                                        && \
    go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0                          && \
    go install github.com/go-delve/delve/cmd/dlv@v1.25.2                               && \
    go install honnef.co/go/tools/cmd/staticcheck@v0.4.7                               && \
    go install golang.org/x/tools/gopls@v0.15.3

ARG GO_ARCH
RUN if [ "${GO_ARCH}" == "alpine" ]; then \
        apk add zsh-vcs make libc6-compat; \
    fi

RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.2.1/zsh-in-docker.sh)" -- \
    -t agnoster \
    -p git \
    -p z \
    -p https://github.com/zsh-users/zsh-autosuggestions \
    -a 'ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE="fg=2"' \
    -p https://github.com/zsh-users/zsh-completions

    
RUN sed -i '/## Main/i prompt_newline() { \n\
 if [[ -n $CURRENT_BG ]]; then\n\
   echo -n " %{%k%F{$CURRENT_BG}%}$SEGMENT_SEPARATOR\n\
%{%k%F{$CURRENT_BG}%}$SEGMENT_SEPARATOR"\n\
 else\n\
   echo -n " %{%k%}"\n\
 fi\n\
 echo -n " %{%f%}"\n\
 CURRENT_BG=''\n\
}' /root/.oh-my-zsh/themes/agnoster.zsh-theme

RUN sed -i 's/  prompt_hg/  prompt_hg\n  prompt_newline/' /root/.oh-my-zsh/themes/agnoster.zsh-theme



ARG KUBECTL_VER
ARG INSTALL_KUBECTL
# Install kubectl
RUN if [ "${INSTALL_KUBECTL}" == "true" ]; then \
        export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac) && \
        export OS=$(uname | awk '{print tolower($0)}')                                  && \
        curl -LO "https://dl.k8s.io/release/v${KUBECTL_VER}/bin/${OS}/${ARCH}/kubectl"  && \
        chmod +x kubectl && mv kubectl /usr/local/bin/kubectl                           && \
        echo 'source <(kubectl completion zsh)'                    >> /root/.zshrc    && \
        echo 'alias k=kubectl'                                      >> /root/.zshrc    && \
        echo 'complete -o default -F __start_kubectl k'             >> /root/.zshrc    && \
        cd "$(mktemp -d)"                                           && \
        KREW="krew-${OS}_${ARCH}"                                   && \
        curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" && \
        tar zxvf "${KREW}.tar.gz"                                   && \
        echo 'export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"'    >> /root/.zshrc    && \
        ./"${KREW}" install krew                                    && \
        ./"${KREW}" install ctx                                     && \
        ./"${KREW}" install ns                                      && \
        ./"${KREW}" install neat                                    ;\
    fi

ARG INSTALL_KUBEBUILDER
ARG KUBEBUILDER_VER
# Install kubebuilder
RUN if [ "${INSTALL_KUBEBUILDER}" == "true" ]; then \
        export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac) && \
        export OS=$(uname | awk '{print tolower($0)}') && \
        export KUBEBUILDER_DL_URL=https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VER} && \
        curl -LO ${KUBEBUILDER_DL_URL}/kubebuilder_${OS}_${ARCH} && \
        chmod +x kubebuilder_${OS}_${ARCH} && \
        mv kubebuilder_${OS}_${ARCH} /usr/local/bin/kubebuilder; \
    fi




