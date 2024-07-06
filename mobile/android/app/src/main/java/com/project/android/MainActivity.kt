package com.project.android

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.core.splashscreen.SplashScreen.Companion.installSplashScreen
import androidx.lifecycle.ViewModelProvider
import androidx.navigation.compose.rememberNavController
import com.project.android.ui.theme.ProjectTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        installSplashScreen()
        enableEdgeToEdge()
//        Text(text = "Mantap")
        setContent {
            ProjectTheme {
//                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
//                    Greeting(
//                        name = "Android",
//                        modifier = Modifier.padding(innerPadding)
//                    )
//                }
//                val navController = rememberNavController()
//                NavGraph(navController = navController)
                val userViewModel = ViewModelProvider(this)[UserViewModel::class.java]
                NavGraph(userViewModel)
            }
        }
    }
}

//@Composable
//fun LoginScreen(onLoginSuccess: () -> Unit) {
//    val context = LocalContext.current.applicationContext
//
//    Column (modifier = Modifier
//        .fillMaxSize()
//        .padding(
//            horizontal = 26.dp,
//            vertical = 140.dp
//        ),
//        verticalArrangement = Arrangement.Bottom,
//        horizontalAlignment = Alignment.CenterHorizontally
//    ) {
//
//    }
//}
