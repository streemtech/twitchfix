FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/

ADD twitchfix /
CMD ["/twitchfix"]
