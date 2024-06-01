"use client"
import "./style.css";
import { useState } from "react"; 
import axios from "axios";
import { useRouter } from 'next/navigation'

export default function Login() {
    const [emailState, setEmailState] = useState("")
    const [passwordState, setPasswordState] = useState("")
    const [emailErrorResponseState, setEmailErrorResponseState] = useState("")
    const [passwordErrorResponseState, setPasswordErrorResponseState] = useState("")
    const [messageErrorResponseState, setMessageErrorResponseState] = useState("")
    const [disableState, setDisableState] = useState(false)
    let loginBody = {email: emailState, password: passwordState}
    const router = useRouter()
    const login = async (e, loginBody) => {
        e.preventDefault()
        try {
            setEmailErrorResponseState("")
            setPasswordErrorResponseState("")
            setMessageErrorResponseState("")
            setDisableState(true)
            const response = await axios.post("/api/v1/user/login", loginBody)
            router.push('/')
        } catch (error) {
            setDisableState(false)
            error.response.data.error.forEach(function(element) {
                if (element.field === "email") {
                    setEmailErrorResponseState(element.message)
                } else if (element.field === "password") {
                    setPasswordErrorResponseState(element.message)
                } else if (element.field === "message") {
                    setMessageErrorResponseState(element.message)
                }
            })
        }
    }
    return (
        <>
            <div className="login-container">
                <div className="login-box">
                    <h2 className="login-box-login">Login</h2>
                    {messageErrorResponseState && <p className="login-box-login-error-message">{messageErrorResponseState}</p>}
                    <div className="email-box">
                        <label htmlFor="email">Email</label>
                        <input type="text" id="email" name="email" value={emailState} onChange={e => setEmailState(e.target.value)}/>
                        {emailErrorResponseState && <p className="message-red">{emailErrorResponseState}</p>}
                    </div>
                    <div className="password-box">
                        <div className="password-label">
                            <label htmlFor="password">Password</label>
                            <a href="#">Forgot your password?</a>
                        </div>
                        <input type="password" id="password" name="password" value={passwordState} onChange={e => setPasswordState(e.target.value)}/>
                        {passwordErrorResponseState && <p className="message-red">{passwordErrorResponseState}</p>}
                    </div>
                    <button type="submit" className="login" onClick={(e) => login(e, loginBody)} disabled={disableState}>Login</button>
                    <div className="login-box-signup"><p>Don&apos;t have an account? <a href="/register">Register</a></p></div>
                    <div className="line"></div>
                    <button className="login-box-google">Google</button>
                </div>
            </div>
        </>
    )
}