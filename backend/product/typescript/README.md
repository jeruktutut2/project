# Product

## init project
npm init  

## zod  
npm i zod  
https://www.npmjs.com/package/zod  

## express  
npm install express  
npm install --save-dev @types/express  
https://www.npmjs.com/package/express  

## jest  
npm install --save-dev jest @types/jest  
https://www.npmjs.com/package/jest  

## babel  
npm install --save-dev babel-jest @babel/preset-env  
https://babeljs.io/setup#installation  
```
{
  "scripts": {
    "test": "jest"
  },
  "jest": {
    "transform": {
      "^.+\\.[t|j]sx?$": "babel-jest"
    }
  }
}
```
create file: babel.config.json
```
{
  "presets": ["@babel/preset-env"]
}
```
npm install --save-dev @babel/preset-typescript  
npm install --save-dev @jest/globals  
https://jestjs.io/docs/getting-started#using-typescript  
add "@babel/preset-typescript" to babel.config.json  

## typescript  
npm install --save-dev typescript  
https://www.npmjs.com/package/typescript  
"main": "index.js",  

## init typescript project  
npx tsc --init  
"target": "es2016"  
"module": "commonjs"  
"moduleResolution": "Node"  
"include": [
    "src/**/*"
]  
"outDir": "./dist"  

## mysql2
npm install --save mysql2  
npm install --save-dev @types/node  


## Create database mysql  
CREATE DATABASE databasename;  

## test  
npx jest path/to/your/test-file.js --runInBand  
npx jest tests/integration-test/services/user-service.test.js --detectOpenHandles  
npx jest  
https://stackoverflow.com/questions/62214949/testing-grpc-functions  

## elasticsearch  
docker pull elasticsearch  
docker pull elasticsearch:7.17.22  
docker run -d --name project-elasticsearch -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:7.17.22  
https://www.elastic.co/guide/en/elasticsearch/client/javascript-api/current/getting-started-js.html  
npm install @elastic/elasticsearch  