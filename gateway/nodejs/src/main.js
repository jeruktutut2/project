import exprss from "./setups/express.js";
import userRoute from "./routes/user-route.js";
import dotenv from 'dotenv';
import grpcSetup from "./setups/grpc.js";

dotenv.config();

grpcSetup.client = await grpcSetup.setClient()

let app = await exprss.setExpress(userRoute.userRouter())
app = await exprss.startExpress(app)
const signal = ["SIGBREAK", "SIGINT", "SIGUSR1", "SIGUSR2", "SIGTERM"]
signal.forEach((eventType) => {
    process.on(eventType, async () => {
        console.log(new Date(), eventType, "stop process");
        await exprss.stopExpress(app)
        process.exit(0)
    });
});