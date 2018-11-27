FROM debian:stable-slim

RUN set -x \
	&& apt-get update -y \
	&& apt-get upgrade -y \
	&& apt-get install --no-install-recommends --no-install-suggests -y ca-certificates supervisor

# Configure supervisors
COPY config/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Configure vbrest
RUN useradd --create-home --shell /bin/bash vbrestservice
COPY vbrest /home/vbrestservice/vbrest
RUN chown vbrestservice:vbrestservice /home/vbrestservice/vbrest
RUN chmod 500 /home/vbrestservice/vbrest

# Start supervisord
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
