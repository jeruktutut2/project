# PROJECT USER

## install grpc  
npm i @grpc/grpc-js  
npm i @grpc/proto-loader  

## why don't use prisma
because if i use dockerfile i need to npx generate or (maybe) copy all file to docker (image)  
https://github.com/prisma/prisma/discussions/20207  
npm uninstall prisma @prisma/client  
delete prisma folder  

## test  
npx jest path/to/your/test-file.js --runInBand  
npx jest tests/integration-test/services/user-service.test.js --detectOpenHandles  
npx jest  
https://stackoverflow.com/questions/62214949/testing-grpc-functions  