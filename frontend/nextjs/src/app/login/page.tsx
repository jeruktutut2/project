"use client"
import "./style.css";
import { useState } from "react"; 
import { login } from "../../services/login.js";

export default function Login() {
    const [emailState, setEmailState] = useState("")
    const [passwordState, setPasswordState] = useState("")
    let loginBody = {email: emailState, password: passwordState}
    return (
        <>
            <div className="login-container">
                <div className="login-box">
                    {/* <form onSubmit={(e) => login(e)}></form> */}
                    <h2 className="login-box-login">Login</h2>
                    <div className="email-box">
                        <label htmlFor="email">Email</label>
                        <input type="text" id="email" name="email" value={emailState} onChange={e => setEmailState(e.target.value)}/>
                    </div>
                    <div className="password-box">
                        <div className="password-label">
                            <label htmlFor="password">Password</label>
                            <a href="#">Forgot your password?</a>
                        </div>
                        <input type="password" id="password" name="password" value={passwordState} onChange={e => setPasswordState(e.target.value)}/>
                    </div>
                    <button type="submit" className="login" onClick={(e) => login(e, loginBody)}>Login</button>
                </div>
            </div>
        </>
    )
}