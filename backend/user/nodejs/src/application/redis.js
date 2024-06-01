import Redis from "ioredis";

let redis = null
const newConnection = async () => {
    redis = new Redis({
        host: process.env.PROJECT_USER_REDIS_HOST,
        port: process.env.PROJECT_USER_REDIS_PORT,
        db: process.env.PROJECT_USER_REDIS_DATABASE
    })
    console.log(new Date(), "redis: connected");
    console.log(new Date(), "redis:", await redis.ping());

    return redis
}

export default {
    redis,
    newConnection
}