const getSessionId = async (req, res, next) => {
    const authorization = req.signedCookies["Authorization"]
    req.sessionId = authorization === undefined ? "" : authorization
    next()
}

export default {
    getSessionId
}