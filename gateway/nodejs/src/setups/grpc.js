import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";

let client
const setClient = async () => {
    const options = {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true,
    };
    const protoPath = import.meta.dirname + "/../protofiles/user.proto"
    var packageDefinition = protoLoader.loadSync(protoPath, options);
    const userProto = grpc.loadPackageDefinition(packageDefinition).protofiles
    console.log(new Date(), "grcp user client: connecting to "+process.env.PROJECT_GATEWAY_USER_APPLICATION_HOST);
    const client = new userProto.UserService(process.env.PROJECT_GATEWAY_USER_APPLICATION_HOST, grpc.credentials.createInsecure())
    console.log(new Date(), "grcp user client: connected to "+process.env.PROJECT_GATEWAY_USER_APPLICATION_HOST);
    return client
}

export default {
    setClient,
    client
}