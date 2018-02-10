FROM scratch

EXPOSE 8080
VOLUME /srv/scripts

ADD httpcmd-linux /srv/httpcmd

ENTRYPOINT ["/srv/httpcmd", "-port", "8080", "-scriptPath", "/srv/scripts/start.sh"]