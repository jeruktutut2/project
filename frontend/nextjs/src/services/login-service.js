import axios from "axios";

export async function login(body) {
    const response = await axios.post("/api/v1/user/login", body)
    return response
 }