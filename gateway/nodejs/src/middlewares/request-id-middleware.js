import { v4 as uuid } from "uuid";

const setRequestId = async (req, res, next) => {
    const requestId = uuid().toString()
    req.requestId = requestId
    next()
}

export default {
    setRequestId
}