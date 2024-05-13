import redisApp from "../application/redis.js";

const newConnection = async () => {
    redisApp.redis = await redisApp.newConnection()
}

const closeConnection = async () => {
    if (redisApp.redis) {
        redisApp.redis.quit()
        console.log(new Date(), "redis: disconnected");
    }
}

export default {
    newConnection,
    closeConnection
}