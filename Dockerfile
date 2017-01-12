FROM iron/go
MAINTAINER John Wesonga <johnwesonga@gmail.com>
COPY server server
ENTRYPOINT ["/server"]



