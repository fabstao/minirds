#****************************
#* (C) Fabian Salamanca
#****************************

FROM debian:latest

#RUN apt-get update
RUN apt-get update && apt-get -y install git
RUN useradd -d /usr/local/goapps -m goapps
RUN su - goapps -c pwd
COPY goapp /usr/local/goapps/
ADD tps /usr/local/goapps/tps/

EXPOSE 8188

ENTRYPOINT ["su","-","goapps","-c","/usr/local/goapps/goapp"]
