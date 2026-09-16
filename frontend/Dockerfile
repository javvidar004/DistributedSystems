FROM node:25-alpine3.22

WORKDIR /app

COPY . .

RUN npm install

EXPOSE 3000

CMD [ "npm", "run", "dev" ]