package com.project.android

import android.hardware.biometrics.BiometricManager.Strings

data class LoginUserRequest(val email: String, val password: String)
data class ErrorResponse(val field: String, val message: String)
data class LoginUserResponse(val data: String, val error: List<ErrorResponse>)