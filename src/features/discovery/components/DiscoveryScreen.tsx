import React, { useState } from 'react';
import { View, Text, TouchableOpacity } from 'react-native';
import { mockProfiles } from '../../../shared/mock/mockProfiles';
import { useNavigation } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { DiscoveryStackParamList } from '../../../navigation/types';
import { SwiperDeck } from './SwiperDeck';

export default function DiscoveryScreen() {
  const navigation = useNavigation<NativeStackNavigationProp<DiscoveryStackParamList>>();
  const [filter, setFilter] = useState<'human' | 'ai'>('human');

  const displayedProfiles = mockProfiles.filter(p => p.type === filter);


  return (
    <View className="flex-1 bg-stone-950 pt-16">
      <View className="flex-row justify-center p-3 gap-3">
        <TouchableOpacity
          className={`py-3 px-6 rounded-full ${filter === 'human' ? 'bg-pink-500' : 'bg-stone-800'}`}
          onPress={() => setFilter('human')}
        >
          <Text className={filter === 'human' ? 'text-white font-bold' : 'text-stone-400 font-medium'}>Real Users</Text>
        </TouchableOpacity>
        <TouchableOpacity
          className={`py-3 px-6 rounded-full ${filter === 'ai' ? 'bg-pink-500' : 'bg-stone-800'}`}
          onPress={() => setFilter('ai')}
        >
          <Text className={filter === 'ai' ? 'text-white font-bold' : 'text-stone-400 font-medium'}>AI Personas</Text>
        </TouchableOpacity>
      </View>
      <View className="flex-1 mt-2">
        <SwiperDeck profiles={displayedProfiles} />
      </View>
    </View>
  );
}
