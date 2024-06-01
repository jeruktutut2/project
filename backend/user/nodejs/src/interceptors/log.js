import grpc from "@grpc/grpc-js";

const setLog = (methodDescriptor, call) => {
    let body
    let requestid = ""
    let md
    const listener = (new grpc.ServerListenerBuilder())
    .withOnReceiveMetadata((metadata, next) => {
        md = metadata.toJSON()
        requestid = metadata.get("requestid")[0]
        next(metadata)
    })
    .withOnReceiveMessage((message, next) => {
        if (methodDescriptor.path === "/protofiles.UserService/Login") {
            body = {email: message.email}
        }
        const requestLog = {grpcRequestTime: new Date() , app: "project-user", requestId: requestid, urlPath: methodDescriptor.path, body: body, metadata: md}
        console.log(JSON.stringify(requestLog));
        next(message)
    })
    .build()
    const responder = (new grpc.ResponderBuilder())
    .withStart((next) => {
        next(listener)
    })
    .withSendStatus((status, next) => {
        const responseLog = { grpcResponseTime: new Date(), app: "project-user", requestId: requestid, body: status}
        console.log(JSON.stringify(responseLog));
        next(status)
    })
    .build()
    return new grpc.ServerInterceptingCall(call, responder)
}

export default {
    setLog
}