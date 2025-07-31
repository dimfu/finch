FROM nginx:latest

EXPOSE 3000

COPY nginx/nginx.conf /etc/nginx/nginx.conf
