FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY tsconfig.json ./
COPY config ./config
COPY src ./src

RUN npm run build

FROM node:22-alpine AS runtime-deps

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --omit=dev

FROM node:22-alpine AS runtime

WORKDIR /app

ENV NODE_ENV=production

COPY --from=runtime-deps /app/node_modules ./node_modules
COPY --from=builder /app/dist/config ./dist/config
COPY --from=builder /app/dist/src ./dist/src
COPY package.json ./

RUN find dist/src -name "*.test.js" -delete \
	&& find dist/src -name "*.test.js.map" -delete

USER node

EXPOSE 3000

CMD ["node", "dist/src/index.js"]
