FROM busybox:1.33
ADD tmp/AppPlaygroundService /
ADD tmp/link_files/ /
CMD ["./AppPlaygroundService", "-c", "/etc/ASUS/app-playground-service.yaml", "serve"]
