FROM node:20.19

WORKDIR /app
COPY . .

RUN npm install --legacy-peer-deps
RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]
