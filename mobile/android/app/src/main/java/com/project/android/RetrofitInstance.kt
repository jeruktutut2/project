package com.project.android

import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

object RetrofitInstance {
    private  const val baseUrl = "http://10.0.2.2:10001"

    private fun getInstance(): Retrofit{
        // Membuat instance dari HttpLoggingInterceptor
        val logging = HttpLoggingInterceptor().apply {
            level = HttpLoggingInterceptor.Level.BODY // Mengatur level logging
        }

        // Menambahkan interceptor ke OkHttpClient
        val httpClient = OkHttpClient.Builder()
            .addInterceptor(logging)
            .build()

        return Retrofit.Builder()
            .baseUrl(baseUrl)
            .addConverterFactory(GsonConverterFactory.create())
            .client(httpClient)
            .build()
    }

    val userApi: UserApi = getInstance().create(UserApi::class.java)
}