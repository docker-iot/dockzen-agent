FROM scratch
ADD agent /
ADD data/server_url.json /data/
ENV http_proxy=http://10.112.1.184:8080
CMD ["/agent"]
