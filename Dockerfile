# Travix Fireball ADK, Including Yeoman generators
FROM node:4.4

MAINTAINER Travix

# NPM packages
RUN npm install -g yo \
	&& npm cache clean

# Docker user to be created to intereact with container. This user is
# different than root
ENV DOCKER_USER docker

# Add a yeoman user to prevent the 'EACCES, permission denied' exception
RUN adduser --disabled-password --gecos "" ${DOCKER_USER}; \
	echo "${DOCKER_USER} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers \
	&& chown -R ${DOCKER_USER} /usr/local/lib/node_modules/

# Set up the appix CLI tool and the apps folder
COPY ./bin/appix-linux /home/${DOCKER_USER}/bin/appix
COPY ./bin/entrypoint.sh /home/${DOCKER_USER}/bin/entrypoint.sh
RUN chmod 755 /home/${DOCKER_USER}/bin/entrypoint.sh \
	&& mkdir /home/${DOCKER_USER}/apps \
	&& chown -R ${DOCKER_USER} /home/${DOCKER_USER}/apps

# Switch to user
USER ${DOCKER_USER}
ENV PATH /home/${DOCKER_USER}/bin:$PATH
ENV HOME /home/${DOCKER_USER}
WORKDIR /home/${DOCKER_USER}

# Copy the deployment key in the image
RUN mkdir -p /home/${DOCKER_USER}/.ssh
COPY .ssh/* /home/${DOCKER_USER}/.ssh/

# Set up the Yeoman generator
RUN npm install -g bitbucket:xivart/travix-fireball-ui-generator

# Set up the RWD website
RUN cd /home/${DOCKER_USER}/ \
	&& git clone git@bitbucket.org:xivart/rwd_cheaptickets_nl.git rwd -b fireball-poc --depth 1
RUN cd rwd/ \
	&& git submodule deinit bin/common \
	&& git rm bin/common \
	&& rm -rf bin/common \
	&& git submodule deinit fireballApps \
	&& git rm fireballApps \
	&& rm -rf fireballApps \
	&& ln -s /home/${DOCKER_USER}/apps fireballApps
RUN cd rwd/cms/src \
	&& git submodule update --init --remote --depth 1 \
	&& cd ../.. \
	&& npm install
EXPOSE 3001

VOLUME ["/home/docker/apps"]
WORKDIR /home/${DOCKER_USER}/apps
ENTRYPOINT ["/home/docker/bin/entrypoint.sh"]
CMD ["/home/docker/rwd/npm", "start"]