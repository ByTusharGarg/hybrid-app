import React from 'react';
import { View, Text, TouchableOpacity, ScrollView } from 'react-native';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigation } from '@react-navigation/native';
import { setSubscriptionTier } from '../../store/slices/appSlice';
import { RootState } from '../../store/store';

export default function SubscriptionScreen() {
  const navigation = useNavigation();
  const dispatch = useDispatch();
  const currentTier = useSelector((state: RootState) => state.app.subscriptionTier);

  const handleUpgrade = () => {
    dispatch(setSubscriptionTier('premium'));
    navigation.goBack();
  };

  return (
    <ScrollView className="flex-1 bg-stone-900" contentContainerStyle={{ padding: 24, paddingTop: 60 }}>
      <TouchableOpacity onPress={() => navigation.goBack()} className="mb-8">
        <Text className="text-stone-400 text-lg">✕ Close</Text>
      </TouchableOpacity>
      
      <Text className="text-3xl font-extrabold text-white mb-2">Upgrade to Premium</Text>
      <Text className="text-stone-400 mb-8">Unlock exclusive AI personas and daily credits.</Text>

      <View className={`rounded-3xl p-6 mb-6 ${currentTier === 'free' ? 'bg-indigo-600' : 'bg-stone-800'}`}>
        <Text className="text-white font-bold text-xl mb-2">Aura Free</Text>
        <Text className="text-indigo-200 mb-4">$0.00 / month</Text>
        <Text className="text-white opacity-80">- Basic real matches</Text>
        <Text className="text-white opacity-80">- 100 intro credits</Text>
        <Text className="text-white opacity-80 mt-1">- Standard texting limits</Text>
      </View>

      <View className={`rounded-3xl p-6 mb-8 border-2 ${currentTier === 'premium' ? 'border-pink-500 bg-stone-800' : 'border-transparent bg-stone-800'}`}>
        <Text className="text-pink-500 font-bold text-xl mb-2">Aura Premium</Text>
        <Text className="text-stone-300 mb-4">$14.99 / month</Text>
        <Text className="text-white opacity-80">- Unlimited AI Personas</Text>
        <Text className="text-white opacity-80">- 500 Credits monthly</Text>
        <Text className="text-white opacity-80 mt-1">- Advanced generative images</Text>
        <Text className="text-white opacity-80 mt-1">- Priority matching</Text>
      </View>

      {currentTier !== 'premium' ? (
        <TouchableOpacity 
          className="bg-pink-500 py-4 rounded-full items-center shadow-lg shadow-pink-500/30"
          onPress={handleUpgrade}
        >
          <Text className="text-white font-bold text-lg">Upgrade Now</Text>
        </TouchableOpacity>
      ) : (
        <View className="bg-stone-800 py-4 rounded-full items-center">
          <Text className="text-pink-500 font-bold text-lg">Active Plan</Text>
        </View>
      )}
    </ScrollView>
  );
}
