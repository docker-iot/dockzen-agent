FROM scratch
ADD agent /
ADD data/server_url.json /data/server_url/
CMD ["/agent"]
