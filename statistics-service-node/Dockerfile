FROM node:20-alpine

WORKDIR /app

ENV NODE_ENV=production

COPY package.json package-lock.json ./
RUN npm ci --omit=dev

COPY config ./config
COPY src ./src

USER node

EXPOSE 3000

CMD ["npm", "start"]
