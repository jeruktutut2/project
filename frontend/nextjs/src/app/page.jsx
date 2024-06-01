"use client"
import axios from "axios";
import { useRouter } from 'next/navigation'

export default function Home() {
    const router = useRouter()
    const logout = async(e) => {
        try {
            const response = await axios.post("/api/v1/user/logout")
            router.push('/login')
        } catch (error) {
            console.log(error);
        }
    }
    return (
        <>
            <h1>/</h1>
            <button onClick={(e) => logout(e)}>logout</button>
        </>  
    );
}
