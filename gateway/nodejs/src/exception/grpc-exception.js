const grpcErrorHandler = (err, res) => {
    let httpCode = 0
    if (err.code === 3) {
        httpCode = 400
    } else if (err.code === 5) {
        httpCode = 403
    } else if (err.code === 13) {
        httpCode = 500
    } else {
        httpCode = 500
    }
    return res.status(httpCode).json({data: "", error: JSON.parse(err.details)})
}

export default {
    grpcErrorHandler
}