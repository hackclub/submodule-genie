FROM golang

MAINTAINER Harrison Shoebridge <harrison@theshoebridges.com>

RUN apt-get update && apt-get install -y cron ssh

ADD . /go/src/github.com/hackclub/submodule-genie
RUN go install github.com/hackclub/submodule-genie

ADD ./scripts/cron-update-submodules.sh /cron-update-submodules.sh
RUN chmod 777 /cron-update-submodules.sh

# Environment variables used by the cron-update-submodules.sh script
ENV SG_DIRECTORY /lecture-hall
ENV SG_REMOTE origin
ENV SG_BRANCH master
ENV SG_FORK_OWNER paked
ENV SG_FORK_REPO lecture-hall
ENV SG_FORK_GIT_REPO git@github.com:paked/lecture-hall.git
ENV SG_OWNER hackclub
ENV SG_REPO lecture-hall
ENV SG_UPSTREAM_REMOTE git@github.com:hackclub/lecture-hall.git
ENV SG_UPSTREAM_BRANCH master
ENV SG_TOKEN <put your token here>
# Set SG_TOKEN while `docker run`ing

# Setup SSH keys
ENV HOME /root
ADD ssh/ /root/.ssh/
RUN chmod 600 /root/.ssh/*

RUN ssh-keyscan github.com > /root/.ssh/known_hosts
RUN cat /root/.ssh/known_hosts

WORKDIR ${SG_DIRECTORY}

# Clone fork repo
RUN git clone ${SG_FORK_GIT_REPO} ${SG_DIRECTORY}
RUN git submodule init
RUN git submodule update --depth 50 

ENTRYPOINT /cron-update-submodules.sh
