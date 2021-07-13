import React from 'react';
import { StyleSheet, Text, View, TouchableOpacity, Button} from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';

export default function HomeScreen({ navigation }) {
    /*
     * HomeScreen is the component used to render the home of Briscola game App
     */
  return (
            <View styles={styles.container}>
                <Text style={styles.textTitle}> {navigation.name} </Text>
                <View styles={{flexDirection:'row'}}>
                    <Button //style={styles.buttonStyle}
                            title="Start Game"
                            onPress={() =>
                                          navigation.navigate('StartGame')
                                      }
                    >
                    </Button>
                    <Button
                        title="Insert Question"
                        onPress={
                            () => navigation.navigate('InsertQuestion')
                        }
                        //style={styles.buttonStyle}
                    >
                    </Button>
                </View>
            </View>
  );
}

const styles = StyleSheet.create({
    container: {
        //flex: 1,
        flexDirection: 'row',
        //width: 100+"%",
        //height: 100+"%",
        //backgroundColor: '#1E4A62',
        //alignItems: 'center',
        //justifyContent: 'center',
    },
    buttomContainer: {
        flexDirection: 'row',
        //backgroundColor: '#1E4A62',
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
        backgroundColor: '#4280ff',
        //width: 30+"%",
        //height: 30+"%",
        //borderRadius : wp('10%'),
    },
    buttonText : {
      fontSize: wp('5%'),
    }
});