import dotenv from 'dotenv';
import redisSetup from "./setup/redis-setup.js";
import grpcSetup from "./setup/grpc-setup.js";

dotenv.config();
await redisSetup.newConnection()
grpcSetup.listen()

const signal = ["SIGBREAK", "SIGINT", "SIGUSR1", "SIGUSR2", "SIGTERM"]
signal.forEach((eventType) => {
    process.on(eventType, async () => {
        console.log(new Date(), eventType, "stop process");
        await redisSetup.closeConnection()
        process.exit(0)
    });
});