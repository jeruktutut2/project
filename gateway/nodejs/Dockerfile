FROM node:21.7.3 AS builder

WORKDIR /app

COPY package*.json ./
    
RUN npm install

COPY . .

RUN npm run build

FROM node:21.7.3-alpine3.20

WORKDIR /app

# COPY package*.json ./
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/src/protofiles ./dist/protofiles

RUN npm install --only=production

# COPY --from=builder /usr/src/app/dist ./

EXPOSE 10002

ENTRYPOINT ["node","dist/main.js"]