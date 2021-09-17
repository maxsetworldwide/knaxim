FROM node:12.22-alpine3.14 as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# production stage
FROM nginx:latest as production-stage
COPY --from=build-stage /app/dist /var/www/frontend
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]



