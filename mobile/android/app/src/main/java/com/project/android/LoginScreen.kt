package com.project.android

import android.widget.Toast
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.material3.Button
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.*
import androidx.compose.runtime.internal.composableLambda
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController

@Composable
//fun LoginScreen(onLoginSuccess: () -> Unit) {
fun LoginScreen(navController: NavHostController, userViewModel: UserViewModel) {
    var username by remember {
        mutableStateOf("email17@email.com")
    }
    var password by remember {
        mutableStateOf("password@A1")
    }
    val context = LocalContext.current.applicationContext
    Column (modifier = Modifier.fillMaxSize(),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(text = "Login", fontSize = 20.sp, fontWeight = FontWeight.Bold)
        Spacer(modifier = Modifier.height(5.dp))
        OutlinedTextField(value = username, onValueChange = { username = it }, label = { Text(text = "email")})
        Spacer(modifier = Modifier.height(5.dp))
        OutlinedTextField(value = password, onValueChange = { password = it}, label = { Text(text = "password")}, visualTransformation = PasswordVisualTransformation())
        Spacer(modifier = Modifier.height(5.dp))
        Button(onClick = {
            userViewModel.login(username, password)
            Toast.makeText(context, "login", Toast.LENGTH_SHORT).show()
//            if (authenticate(username, password)) {
////                onLoginSuccess()
////                navController.popBackStack("home", inclusive = false)
//                navController.navigate("home")
//                Toast.makeText(context, "successfully login", Toast.LENGTH_SHORT).show()
//            } else {
//                Toast.makeText(context, "not successfully login", Toast.LENGTH_SHORT).show()
//            }
        }) {
            Text(text = "Login")
        }
    }
}

private fun authenticate(username: String, password: String): Boolean{
    val validUsername = "admin"
    var validPassword = "admin"
    return username == validUsername && password == validPassword
}

@Composable
//fun NavGraph(navController: NavHostController) {
fun NavGraph(userViewModel: UserViewModel) {
//    NavHost(navController = navController, startDestination = "login") {
//        composable("login"){
//            LoginScreen(onLoginSuccess = {
//                navController.navigate("home"){
//                    popUpTo(0)
//                }
//            })
//        }
//        composable("home"){
//            HomeScreen()
//        }
//    }
    val context = LocalContext.current.applicationContext
    val navController = rememberNavController()
    NavHost(navController = navController, startDestination = "login") {
        composable("login"){
            LoginScreen(navController = navController, userViewModel)
        }
        composable("home"){
//            Toast.makeText(context, "home", Toast.LENGTH_SHORT).show()
            HomeScreen()
        }
    }
}