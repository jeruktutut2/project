const setLogRequest = async (req, res, next) => {
    let body = ""
    if (req.originalUrl === "/api/v1/user/login") {
        body = { email: req.body.email}
    } else {
        body = req.body
    }
    const requestLog = {requestTime: new Date(), app: "project-gateway", method: req.method, requestId: req.requestId, host: req.headers.host, urlPath: req.originalUrl, protocol: req.protocol, body: body, userAgent: req.get('User-Agent'), remoteAddr: req.socket.remoteAddress, forwardedFor: req.headers['x-forwarded-for'], headers: req.headers}
    console.log(JSON.stringify(requestLog))
    next()
}

// const logResponse = async (res, httpCode, requestId, response) => {
const logResponse = async (requestId, response) => {
    // console.log("new Date():", new Date(), requestId, response);
    const resp = JSON.stringify({responseTime: new Date(), app: "project-gateway", requestId: requestId, response: response})
    console.log(resp);
    // return res.status(httpCode).json(response)
}

export default {
    setLogRequest,
    logResponse
}