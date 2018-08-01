#****************************
#* (C) Fabian Salamanca
#****************************

FROM debian:latest

#RUN apt-get update
RUN apt-get -y update && apt-get -y install git apt-utils
RUN useradd -d /usr/local/goapps -m goapps
RUN su - goapps -c pwd
COPY goapp /usr/local/goapps/
ADD tps /usr/local/goapps/tps/

EXPOSE 8800

ENTRYPOINT ["su","-","goapps","-c","/usr/local/goapps/goapp"]
