/*
 * App.js contains skeleton of Quiz App 
 */
import { StatusBar } from 'expo-status-bar';
import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import Quiz from "./src/scenes/quiz.js"

export default function App() {
  /*
   * Function call at start of the App to render everything 
   */
  return (
      <Quiz styles={styles.container}></Quiz>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    width: 100+"%",
    height: 100+"%",
    backgroundColor: '#1E4A62',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
