const setErrorMessage = (message) => {
    return JSON.stringify([{field: "message", message: message}])
}

const setInternalServerErrorMessage = () => {
    return JSON.stringify([{field: "message", message: "internal server error"}])
}

export default {
    setErrorMessage,
    setInternalServerErrorMessage
}