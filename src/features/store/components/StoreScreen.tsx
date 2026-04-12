import React from 'react';
import { View, Text, ScrollView, TouchableOpacity, Image } from 'react-native';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/store';
import { mockProducts } from '../../../shared/mock/mockProducts';

export default function StoreScreen() {
  const credits = useSelector((state: RootState) => state.app.credits);

  return (
    <View className="flex-1 bg-stone-950 pt-16 px-4">
      <View className="flex-row justify-between items-center mb-6 px-2">
        <Text className="text-3xl font-extrabold text-white">Aura Store</Text>
        <View className="bg-stone-800 px-4 py-2 rounded-full flex-row items-center">
          <Text className="text-pink-500 font-bold mr-1">{credits}</Text>
          <Text className="text-stone-300">🪙</Text>
        </View>
      </View>

      <ScrollView className="flex-1" showsVerticalScrollIndicator={false}>
        <View className="bg-indigo-600 rounded-3xl p-6 mb-8 mx-2 overflow-hidden relative">
          <Text className="text-white font-bold text-2xl mb-2">Gift Pack</Text>
          <Text className="text-indigo-200 mb-4 w-2/3">Boost your matches with premium digital gifts and super likes.</Text>
          <TouchableOpacity className="bg-white/20 self-start px-4 py-2 rounded-full backdrop-blur-md">
            <Text className="text-white font-medium">Get 500 🪙</Text>
          </TouchableOpacity>
        </View>

        <Text className="text-xl font-bold text-stone-200 mb-4 px-2">Send to AI Personas</Text>
        <View className="flex-row flex-wrap justify-between px-2">
          {mockProducts.map((product) => (
            <View key={product.id} className="w-[48%] bg-stone-900 rounded-2xl p-4 mb-4 border border-stone-800">
              <Text className="text-3xl mb-2">{product.imageUrl}</Text>
              <Text className="text-white font-bold text-lg mb-1">{product.name}</Text>
              <Text className="text-stone-400 text-sm mb-3 h-10">{product.description}</Text>
              <TouchableOpacity className="bg-pink-500/20 py-2 rounded-xl items-center border border-pink-500/30">
                <Text className="text-pink-500 font-bold">{product.priceCredits} 🪙</Text>
              </TouchableOpacity>
            </View>
          ))}
        </View>
      </ScrollView>
    </View>
  );
}
