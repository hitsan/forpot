FROM denoland/deno:ubuntu-2.3.1

RUN apt-get -y update \
 && apt-get -y --no-install-recommends install openssh-server \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*
RUN ssh-keygen -A \
 && mkdir -p /run/sshd
RUN echo 'root:password' | chpasswd
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
