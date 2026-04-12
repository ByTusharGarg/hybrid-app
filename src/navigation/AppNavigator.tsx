import React from 'react';
import { useSelector } from 'react-redux';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { NavigationContainer } from '@react-navigation/native';
import { RootStackParamList } from './types';
import { RootState } from '../features/store/store';

import AuthNavigator from './AuthNavigator';
import MainTabNavigator from './MainTabNavigator';
import SubscriptionScreen from '../features/profile/components/SubscriptionScreen';

const RootStack = createNativeStackNavigator<RootStackParamList>();

export default function AppNavigator() {
  const isAuthenticated = useSelector((state: RootState) => state.app.isAuthenticated);

  return (
    <NavigationContainer>
      <RootStack.Navigator screenOptions={{ headerShown: false, animation: 'fade' }}>
        {!isAuthenticated ? (
          <RootStack.Screen name="Auth" component={AuthNavigator} />
        ) : (
          <>
            <RootStack.Screen name="Main" component={MainTabNavigator} />
            <RootStack.Screen 
              name="Subscription" 
              component={SubscriptionScreen} 
              options={{ presentation: 'modal', animation: 'slide_from_bottom' }} 
            />
          </>
        )}
      </RootStack.Navigator>
    </NavigationContainer>
  );
}
