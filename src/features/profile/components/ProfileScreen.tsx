import React from 'react';
import { View, Text, TouchableOpacity } from 'react-native';
import { useSelector, useDispatch } from 'react-redux';
import { RootState } from '../../store/store';
import { logout } from '../../store/slices/appSlice';
import { useNavigation } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { RootStackParamList } from '../../../navigation/types';

export default function ProfileScreen() {
  const dispatch = useDispatch();
  const navigation = useNavigation<NativeStackNavigationProp<RootStackParamList>>();
  const { userProfile, credits, subscriptionTier } = useSelector((state: RootState) => state.app);

  return (
    <View className="flex-1 bg-stone-900 pt-16 px-6">
      <View className="items-center mb-10">
        <View className="w-24 h-24 bg-stone-700 rounded-full mb-4 items-center justify-center">
          <Text className="text-3xl text-stone-300">👤</Text>
        </View>
        <Text className="text-2xl font-bold text-white">{userProfile.name}</Text>
        <Text className="text-stone-400 capitalize">{subscriptionTier} Member</Text>
      </View>

      <View className="bg-stone-800 rounded-2xl p-5 mb-6 flex-row justify-between items-center">
        <View>
          <Text className="text-stone-400 text-sm mb-1">Available Credits</Text>
          <Text className="text-pink-500 text-2xl font-bold">{credits} 🪙</Text>
        </View>
        <TouchableOpacity 
          className="bg-stone-700 px-4 py-2 rounded-full"
          onPress={() => navigation.navigate('Subscription')}
        >
          <Text className="text-white font-medium">Get More</Text>
        </TouchableOpacity>
      </View>

      <TouchableOpacity 
        className="bg-stone-800 rounded-2xl p-5 mb-4 flex-row justify-between items-center"
        onPress={() => navigation.navigate('Subscription')}
      >
        <Text className="text-white font-medium text-lg">Manage Subscription</Text>
        <Text className="text-stone-500">›</Text>
      </TouchableOpacity>

      <TouchableOpacity 
        className="mt-auto mb-10 bg-red-500/10 py-4 rounded-xl items-center border border-red-500/20"
        onPress={() => dispatch(logout())}
      >
        <Text className="text-red-500 font-bold">Log Out</Text>
      </TouchableOpacity>
    </View>
  );
}
