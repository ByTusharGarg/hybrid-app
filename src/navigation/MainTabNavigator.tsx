import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { MainTabParamList } from './types';
import DiscoveryStackNavigator from './DiscoveryStackNavigator';
import StoreScreen from '../features/store/components/StoreScreen';
import ProfileScreen from '../features/profile/components/ProfileScreen';
import { Text } from 'react-native';
import Animated, { useAnimatedStyle, withSpring, useSharedValue } from 'react-native-reanimated';

const Tab = createBottomTabNavigator<MainTabParamList>();

function AnimatedTabIcon({ icon, focused, color }: { icon: string, focused: boolean, color: string }) {
  const scale = useSharedValue(1);

  React.useEffect(() => {
    scale.value = withSpring(focused ? 1.2 : 1, { damping: 10, stiffness: 100 });
  }, [focused]);

  const animatedStyle = useAnimatedStyle(() => {
    return {
      transform: [{ scale: scale.value }],
    };
  });

  return (
    <Animated.View style={animatedStyle}>
      <Text style={{ color, fontSize: 20 }}>{icon}</Text>
    </Animated.View>
  );
}

export default function MainTabNavigator() {
  return (
    <Tab.Navigator screenOptions={{ 
      headerShown: false,
      tabBarStyle: { backgroundColor: '#1c1917', borderTopColor: '#292524', height: 60, paddingBottom: 10, paddingTop: 5 },
      tabBarActiveTintColor: '#ec4899',
      tabBarInactiveTintColor: '#78716c',
      animation: 'shift'
    }}>
      <Tab.Screen 
        name="Discovery" 
        component={DiscoveryStackNavigator} 
        options={{ 
          tabBarIcon: ({ color, focused }) => <AnimatedTabIcon icon="🔥" focused={focused} color={color} />,
          tabBarLabel: 'Aura'
        }} 
      />
      <Tab.Screen 
        name="Store" 
        component={StoreScreen} 
        options={{ 
          tabBarIcon: ({ color, focused }) => <AnimatedTabIcon icon="🛍️" focused={focused} color={color} />,
          tabBarLabel: 'Store'
        }} 
      />
      <Tab.Screen 
        name="Profile" 
        component={ProfileScreen} 
        options={{ 
          tabBarIcon: ({ color, focused }) => <AnimatedTabIcon icon="👤" focused={focused} color={color} />,
          tabBarLabel: 'Profile'
        }} 
      />
    </Tab.Navigator>
  );
}
