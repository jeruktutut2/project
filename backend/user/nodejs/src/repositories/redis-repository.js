const set = async (redis, key, value) => {
    return redis.set(key, value)
}

const del = async (redis, sessionId) => {
    return await redis.del(sessionId)
}

export default {
    set,
    del
}