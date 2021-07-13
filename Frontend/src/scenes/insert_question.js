import React from 'react';
import { StyleSheet, Text, View, TouchableOpacity} from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';

export default function InsertQuestion({ navigation }) {
    /*
     * HomeScreen is the component used to render the home of Briscola game App
     */
  return (
            <View style={styles.container}>
                <Text style={styles.textTitle}> {navigation.name} </Text>
                <TouchableOpacity style={styles.buttonStyle}
                        onPress={() =>
                                      navigation.navigate('StartGame')
                                  }
                >
                    <Text style={styles.buttonText}>Emacs</Text>
                </TouchableOpacity>
            </View>
  );
}

const styles = StyleSheet.create({
  container: {
      alignItems: 'center',
      marginTop: hp('15%'),
  },
  image: {
      flex: 1,
      resizeMode: "cover",
      justifyContent: "center",
  },
  textTitle : {
      fontSize : wp('10%'),
      color: '#4280ff',
      fontWeight: 'bold'
  },
  buttonStyle : {
      marginTop: hp('10%'),
      backgroundColor: '#4280ff',
      padding: wp('10%'),
      borderRadius : wp('10%'),
  },
  buttonText : {
      fontSize: wp('10%'),
  }
});