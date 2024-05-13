import dotenv from 'dotenv';
import redisSetup from "./setup/redis-setup.js";
import expressSetup from "./setup/express-setup.js";

dotenv.config();

await redisSetup.newConnection()

const app = await expressSetup.listen()

const signal = ["SIGBREAK", "SIGINT", "SIGUSR1", "SIGUSR2", "SIGTERM"]
signal.forEach((eventType) => {
    process.on(eventType, async () => {
        console.log(new Date(), eventType, "stop process");
        await expressSetup.close(app)
        await redisSetup.closeConnection()
        process.exit(0)
    });
});