FROM scratch
ADD agent /
ADD ./data/server_url.json /data/
CMD ["/agent"]
