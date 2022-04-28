# Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

FROM bhojpur/platform-postgres
ARG TRIGGER_REBUILD=1

USER root

ENV PROTOC_ZIP=protoc-3.7.1-linux-x86_64.zip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/$PROTOC_ZIP && \
    unzip -o $PROTOC_ZIP -d /usr/local bin/protoc && \
    unzip -o $PROTOC_ZIP -d /usr/local 'include/*' && \
    rm -f $PROTOC_ZIP

RUN curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
RUN curl -L https://download.docker.com/linux/static/stable/x86_64/docker-19.03.5.tgz | tar xz && \
    mv docker/docker /usr/bin && \
    rm -rf docker

RUN curl -o /usr/bin/k3s -L https://github.com/rancher/k3s/releases/download/v1.0.1/k3s && \
    chmod +x /usr/bin/k3s

ENV GORPA_APPLICATION_ROOT=/application/piro
RUN curl -L https://github.com/bhojpur/gorpa/releases/download/v1.0.0/gorpa-v1.0.0-Linux-x86_64 && \
    mv gorpa-v1.0.0-Linux-x86_64 /usr/bin/gorpa && \
    rm README.md

RUN curl -L https://get.helm.sh/helm-v3.3.0-rc.2-linux-amd64.tar.gz | tar xz linux-amd64/helm && \
    mv linux-amd64/helm /usr/bin && \
    rm -r linux-amd64
