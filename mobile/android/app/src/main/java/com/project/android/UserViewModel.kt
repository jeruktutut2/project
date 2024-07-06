package com.project.android

import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.launch

class UserViewModel: ViewModel() {
    private val userApi = RetrofitInstance.userApi
    fun login(email: String, password: String) {
        Log.i("login", "login")
        viewModelScope.launch {
            val loginUserRequest = LoginUserRequest(email, password)
            Log.i("responselogin", loginUserRequest.toString())
            val response = userApi.login(loginUserRequest)
            Log.i("responselogin", response.toString())
            Log.i("responselogin", response.code().toString())
            Log.i("responselogin", response.message())
            Log.i("responselogin", response.body().toString())
            Log.i("responselogin", response.errorBody().toString())
        }
    }
}