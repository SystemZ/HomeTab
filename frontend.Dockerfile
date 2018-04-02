FROM nginx:mainline-alpine
COPY frontend/dist /usr/share/nginx/html