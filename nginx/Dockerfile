FROM nginx:1.25.0

COPY nginx_active_active.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
