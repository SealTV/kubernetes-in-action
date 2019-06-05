FROM ubuntu:latest

RUN apt-get update; apt-get -y install git
ADD git-sync.sh /bin/git-sync.sh
ENTRYPOINT [ "/bin/git-sync.sh" ]