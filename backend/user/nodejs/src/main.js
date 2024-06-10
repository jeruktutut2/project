import grpcSetup from "./setup/grpc-setup.js";
import mysqlUtil from "./utils/mysql-util.js";
import redisUtil from "./utils/redis-util.js";

mysqlUtil.mysqlPool = mysqlUtil.newConnection()
redisUtil.redis = await redisUtil.newConnection()
grpcSetup.listen()

const signal = ["SIGBREAK", "SIGINT", "SIGUSR1", "SIGUSR2", "SIGTERM"]
signal.forEach((eventType) => {
    process.on(eventType, async () => {
        console.log(new Date(), eventType, "stop process");
        mysqlUtil.closeConnection(mysqlUtil.mysqlPool)
        await redisUtil.closeConnection()
        process.exit(0)
    });
});