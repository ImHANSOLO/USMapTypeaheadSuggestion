# Frontend Dockerfile

FROM node:16-alpine
WORKDIR /app
COPY .env .env
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build
RUN npm install -g serve
EXPOSE 8080
CMD ["serve","-s","dist","-l","8080"]
