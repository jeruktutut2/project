import axios from "axios";

export async function login(e, body) {
    e.preventDefault()
    console.log("body:", body);
    const response = await axios.post("/api/v1/user/login", body)
    console.log("login:", body);
    console.log("response:", response);
 }