FROM node:18.8.0-alpine

RUN mkdir app

WORKDIR /app

RUN npm i -g http-server

COPY ./package*.json ./

RUN npm install

ARG VITE_API_URL \
    PORT=3000

ENV VITE_API_URL=${VITE_API_URL}  \
    PORT=${PORT}

COPY . .

RUN npm run build

CMD [ "http-server", "dist"]