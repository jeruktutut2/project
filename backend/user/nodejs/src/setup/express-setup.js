import { web } from "../application/express.js";

const listen = async () => {
    const app = web.listen(process.env.PROJECT_USER_APPLICATION_PORT, () => {
        console.log(new Date(), "express: started on port "+process.env.PROJECT_USER_APPLICATION_PORT);
    })
    return app
}

const close = async (app) => {
    app.close( () => {
        console.log(new Date(), "express: has closed all connections.");
    })
}

export default {
    listen,
    close
}