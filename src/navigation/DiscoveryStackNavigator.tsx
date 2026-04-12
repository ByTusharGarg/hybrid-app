import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { DiscoveryStackParamList } from './types';
import DiscoveryScreen from '../features/discovery/components/DiscoveryScreen';
import ChatScreen from '../features/chat/components/ChatScreen';

const DiscoveryStack = createNativeStackNavigator<DiscoveryStackParamList>();

export default function DiscoveryStackNavigator() {
  return (
    <DiscoveryStack.Navigator screenOptions={{ headerShown: false }}>
      <DiscoveryStack.Screen name="DiscoveryFeed" component={DiscoveryScreen} />
      <DiscoveryStack.Screen name="Chat" component={ChatScreen} />
    </DiscoveryStack.Navigator>
  );
}
