import React from 'react';
import { View, Text, TouchableOpacity, TextInput } from 'react-native';
import { useDispatch } from 'react-redux';
import { login } from '../../store/slices/appSlice';
import Animated, { useSharedValue, useAnimatedStyle, withRepeat, withTiming, withSequence } from 'react-native-reanimated';

export default function LoginScreen() {
  const dispatch = useDispatch();
  const scale = useSharedValue(1);

  React.useEffect(() => {
    scale.value = withRepeat(
      withSequence(
        withTiming(1.05, { duration: 1000 }),
        withTiming(1, { duration: 1000 })
      ),
      -1,
      true
    );
  }, []);

  const animatedButtonStyle = useAnimatedStyle(() => ({
    transform: [{ scale: scale.value }]
  }));

  const handleLogin = () => {
    dispatch(login());
  };

  return (
    <View className="flex-1 bg-stone-900 justify-center px-8">
      <View className="items-center mb-12">
        <Text className="text-4xl font-extrabold text-white tracking-tight mb-2">Aura</Text>
        <Text className="text-lg text-stone-400">Match. Chat. Connect.</Text>
      </View>

      <View className="space-y-4">
        <TextInput 
          placeholder="Email Address" 
          placeholderTextColor="#78716c"
          className="bg-stone-800 text-white px-5 py-4 rounded-2xl text-base"
        />
        <TextInput 
          placeholder="Password" 
          placeholderTextColor="#78716c"
          secureTextEntry
          className="bg-stone-800 text-white px-5 py-4 rounded-2xl text-base"
        />
      </View>

      <Animated.View style={animatedButtonStyle} className="mt-8">
        <TouchableOpacity 
          onPress={handleLogin}
          className="bg-indigo-500 py-4 rounded-full items-center shadow-lg shadow-indigo-500/30"
        >
          <Text className="text-white font-bold text-lg">Continue to Aura</Text>
        </TouchableOpacity>
      </Animated.View>

      <View className="mt-8 flex-row justify-center">
        <Text className="text-stone-400">Don't have an account? </Text>
        <TouchableOpacity>
          <Text className="text-indigo-400 font-bold">Sign up</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}
