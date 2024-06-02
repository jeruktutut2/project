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
    console.log(new Date(), "grcp user client: connecting to 0.0.0.0:"+process.env.PROJECT_USER_APPLICATION_PORT);
    const client = new userProto.UserService("0.0.0.0:"+process.env.PROJECT_USER_APPLICATION_PORT, grpc.credentials.createInsecure())
    console.log(new Date(), "grcp user client: connected to 0.0.0.0:"+process.env.PROJECT_USER_APPLICATION_PORT);
    return client
}

export default {
    setClient,
    client
}