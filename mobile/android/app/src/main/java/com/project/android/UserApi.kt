package com.project.android

import retrofit2.Call
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST

interface UserApi {
    @POST("/api/v1/users/login")
//    suspend fun login(username: String, password: String)
    suspend fun login(@Body userRequest: LoginUserRequest): Response<LoginUserResponse>
}