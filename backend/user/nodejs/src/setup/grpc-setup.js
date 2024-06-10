import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import logInterceptor from "../interceptors/log.js";
import { Register, Login, Logout } from "../grpc/user-grpc.js";

const listen = () => {
    const options = {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true,
    }

    const protoPath = import.meta.dirname + "/../protofiles/user.proto"
    const packageDefinition = protoLoader.loadSync(protoPath, options)
    const userProto = grpc.loadPackageDefinition(packageDefinition).protofiles

    const server = new grpc.Server({interceptors: [logInterceptor.setLog]})
    server.addService(userProto.UserService.service, { Register, Login, Logout })
    server.bindAsync(
        process.env.PROJECT_USER_APPLICATION_HOST,
        grpc.ServerCredentials.createInsecure(),
        (error, port) => {
            if (error) {
                console.log("error:", error);
            }
            console.log((new Date()).toISOString(), "grpc: started on host", process.env.PROJECT_USER_APPLICATION_HOST);
        }
    )
}

export default {
    listen
}