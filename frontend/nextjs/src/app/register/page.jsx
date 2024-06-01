"use client"

import { useState } from "react";
import "./style.css";
import axios from "axios";
import { useRouter } from 'next/navigation'

export default function Register() {
    const [usernameState, setUsernameState] = useState("")
    const [emailState, setEmailState] = useState("")
    const [passwordState, setPasswordState] = useState("")
    const [confirmPasswordState, setConfirmPasswordState] = useState("")
    const [utcState, setUtcState] = useState("")
    const [usernameErrorResponseState, setUsernameErrorResponseState] = useState("")
    const [emailErrorResponseState, setEmailErrorResponseState] = useState("")
    const [passwordErrorResponseState, setPasswordErrorRespoonseState] = useState("")
    const [confirmPasswordErrorResponseState, setConfirmPasswordErrorResponseState] = useState("")
    const [utcErrorResponseState, setUtcErrorResponseState] = useState("")
    const [messageErrorResponseState, setMessageErrorResponseState] = useState("")
    let registerBody = {
        username: usernameState,
        email: emailState,
        password: passwordState,
        confirmPassword: confirmPasswordState,
        utc: utcState
    }
    const router = useRouter()
    const register = async (e, registerBody) => {
        e.preventDefault()
        setUsernameErrorResponseState("")
        setEmailErrorResponseState("")
        setPasswordErrorRespoonseState("")
        setConfirmPasswordErrorResponseState("")
        setUtcErrorResponseState("")
        setMessageErrorResponseState("")
        try {
            await axios.post("/api/v1/user/register", registerBody)
            router.push('/login')
        } catch (error) {
            console.log(error);
            error.response.data.error.forEach(element => {
                if (element.field === "username") {
                    setUsernameErrorResponseState(element.message)
                } else if (element.field === "email") {
                    setEmailErrorResponseState(element.message)
                } else if (element.field === "password") {
                    setPasswordErrorRespoonseState(element.message)
                } else if (element.field === "confirmpassword") {
                    setConfirmPasswordErrorResponseState(element.message)
                } else if (element.field === "utc") {
                    setUtcErrorResponseState(element.message)
                } else if (element.field === "message") {
                    setMessageErrorResponseState(element.message)
                }

            })
        }
    }
    return (
        <>
            <div className="register-container">
                <div className="register-box">
                    <h2 className="register-box-register">Register</h2>
                    {messageErrorResponseState && <p className="register-box-register-error-message">{messageErrorResponseState}</p>}
                    <div className="register-box-username">
                        <label htmlFor="username">Username</label>
                        <input type="text" id="username" name="username" value={usernameState} onChange={e => setUsernameState(e.target.value)}/>
                        {usernameErrorResponseState && <p className="message-red">{usernameErrorResponseState}</p>}
                    </div>
                    <div className="register-box-email">
                        <label htmlFor="email">Email</label>
                        <input type="text" id="email" name="email" value={emailState} onChange={e => setEmailState(e.target.value)}/>
                        {emailErrorResponseState && <p className="message-red">{emailErrorResponseState}</p>}
                    </div>
                    <div className="register-box-password">
                        <label htmlFor="password">Password</label>
                        <input type="password" id="password" name="password" value={passwordState} onChange={e => setPasswordState(e.target.value)}/>
                        {passwordErrorResponseState && <p className="message-red">{passwordErrorResponseState}</p>}
                    </div>
                    <div className="register-box-confirmpassword">
                        <label htmlFor="confirmpassword">Confirm Password</label>
                        <input type="password" id="confirmpassword" name="confirmpassword" value={confirmPasswordState} onChange={e => setConfirmPasswordState(e.target.value)}/>
                        {confirmPasswordErrorResponseState && <p className="message-red">{confirmPasswordErrorResponseState}</p>}
                    </div>
                    <div className="register-box-utc">
                        <label htmlFor="utc">UTC</label>
                        <input type="text" id="utc" name="utc" value={utcState} onChange={e => setUtcState(e.target.value)} />
                        {utcErrorResponseState && <p className="message-red">{utcErrorResponseState}</p>}
                    </div>
                    <button className="register-box-register-button" onClick={(e) => register(e, registerBody)}>Register</button>
                    <div className="register-box-login"><p>Already have an account? <a href="/login">Login</a></p></div>
                </div>
            </div>
        </>
    )
}