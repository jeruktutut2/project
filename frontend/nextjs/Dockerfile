FROM node:lts-alpine3.20

RUN apk add --no-cache libc6-compat
# RUN npm i -g npm

EXPOSE 10000

ENV PORT 10000
ENV NODE_ENV production

WORKDIR /home/nextjs/app

COPY package.json .
COPY package-lock.json .

RUN npm install --omit=optional
RUN npx browserslist@latest --update-db
# RUN npx next telemetry disable

# need to install linux specific swc builds
RUN npm install -D @swc/cli @swc/core

COPY . .

RUN npm run build

CMD [ "npm", "start" ]