import Redis from "ioredis";

let redis

const newConnection = async () => {
    console.log(new Date(), "redis: connecting to", process.env.PROJECT_USER_REDIS_HOST + ":" + process.env.PROJECT_USER_REDIS_PORT);
    const redisConnection = new Redis({
        host: process.env.PROJECT_USER_REDIS_HOST,
        port: process.env.PROJECT_USER_REDIS_PORT,
        db: process.env.PROJECT_USER_REDIS_DATABASE
    })
    console.log(new Date(), "redis: connected to", process.env.PROJECT_USER_REDIS_HOST + ":" + process.env.PROJECT_USER_REDIS_PORT);
    console.log(new Date(), "redis:", await redisConnection.ping());

    return redisConnection
}

const closeConnection = async (redis) => {
    if (redis) {
        console.log(new Date(), "redis: closing connection", process.env.PROJECT_USER_REDIS_HOST + ":" + process.env.PROJECT_USER_REDIS_PORT);
        await redis.quit()
        console.log(new Date(), "redis: closed connection", process.env.PROJECT_USER_REDIS_HOST + ":" + process.env.PROJECT_USER_REDIS_PORT);
    }
}

export default {
    redis,
    newConnection,
    closeConnection
}