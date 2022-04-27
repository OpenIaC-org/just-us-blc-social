import React, {useState} from 'react';
import { AuthContext } from '../context/context';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { 
    StyleSheet,
    Text,
    View,
    Image,
    TextInput,
    Button,
    TouchableOpacity,
 } from 'react-native';

 

 
export default function Login({ navigation }) {


    return (
        <View style={styles.container}>
            <Text style={styles.h1}>Just us</Text>
            <Image style={styles.image} source={require("../assets/lock.png")} />

            <TouchableOpacity style={styles.loginBtn} onPress={() => navigation.navigate("LoginUser")}>
                <Text style={styles.loginText} >LOGIN</Text>
            </TouchableOpacity>

            <TouchableOpacity>
                <Text style={styles.forgot_button}>Forgot Password?</Text>
            </TouchableOpacity>

            <TouchableOpacity style={styles.loginBtn} onPress={() => navigation.navigate("CreateUser")}>
                <Text style={styles.loginText}>CREATE NEW USER</Text>
            </TouchableOpacity>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: "#fff",
      alignItems: "center",
      justifyContent: "center",
    },
  
    image: {
      height: "30%",
      resizeMode: 'contain'
    },
  
    inputView: {
      backgroundColor: "#FFC0CB",
      borderRadius: 30,
      width: "70%",
      height: 45,
      marginBottom: 20,
  
      alignItems: "center",
    },
  
    TextInput: {
      height: 50,
      flex: 1,
      padding: 10,
      marginLeft: 20,
    },
  
    forgot_button: {
      height: 30,
      marginBottom: 30,
    },
  
    loginBtn: {
      width: "80%",
      borderRadius: 10,
      height: 50,
      alignItems: "center",
      justifyContent: "center",
      marginTop: 40,
      backgroundColor: "#157EFB",
    },
    h1: {
      fontSize: 80,
      fontWeight: "bold",
    },
  });

