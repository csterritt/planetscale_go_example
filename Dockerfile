# Select the base system
FROM ubuntu:bionic

# Setup the locales for the Ubuntu system.  Because the base image is a bare bones
# setup, this is needed to get things in the correct language.
# https://hub.docker.com/_/ubuntu/
# Always use update with the install subcommand
# https://docs.docker.com/develop/develop-images/dockerfile_best-practices/
RUN apt-get update && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8

# Set LANG to us.utf8
ENV LANG en_US.utf8

# Set tell the installer that we are working in a noninteractive ENV
ENV DEBIAN_FRONTEND noninteractive

# Install dumb-init
#RUN wget -O /usr/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64
#RUN chmod +x /usr/bin/dumb-init
RUN apt-get update && apt-get -y install dumb-init

# Install planetscale cli
COPY pscale_0.38.0_linux_amd64.deb pscale_0.38.0_linux_amd64.deb
RUN dpkg -i pscale_0.38.0_linux_amd64.deb
RUN rm pscale_0.38.0_linux_amd64.deb

# Set planetscale environmental variables
ENV PLANETSCALE_ORG=chris-st
ENV PLANETSCALE_SERVICE_TOKEN_NAME=your_service_token_name_here
ENV PLANETSCALE_SERVICE_TOKEN=your_service_token_here

# Set up run directory
RUN mkdir /app

# Install my program
COPY ps_ws_ex /app/ps_ws_ex

# Expose the port
EXPOSE 8080

# Runs "/usr/bin/dumb-init -- [whatever is in CMD below]"
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/usr/local/bin/pscale", "connect", "firstexample", "main", "--host", "0.0.0.0", "--execute", "/app/ps_ws_ex"]
