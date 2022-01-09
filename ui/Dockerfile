FROM node:16.13.1-alpine as build

WORKDIR /app
ENV PATH /path/node_modules/.bin:$PATH
COPY package.json ./
COPY package-lock.json ./
RUN npm install --silent
RUN npm install react-scripts -g --silent
COPY . ./
RUN npm run build

FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.template /etc/nginx/templates/default.conf.template

ENV NGINX_PORT 80
ENV CATALOGUE_URL "http://catalogue:8085/"
ENV FEATURESTORE_URL "http://featurestore:8085/"
ENV METRICSTORE_URL "http://metricstore:8085/"

CMD ["nginx", "-g", "daemon off;"]