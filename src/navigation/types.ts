import { NavigatorScreenParams } from '@react-navigation/native';

export type AuthStackParamList = {
  Login: undefined;
};

export type DiscoveryStackParamList = {
  DiscoveryFeed: undefined;
  Chat: {
    matchId?: string;
    personaId?: string;
  };
};

export type MainTabParamList = {
  Discovery: NavigatorScreenParams<DiscoveryStackParamList>;
  Store: undefined;
  Profile: undefined;
};

export type RootStackParamList = {
  Auth: NavigatorScreenParams<AuthStackParamList>;
  Main: NavigatorScreenParams<MainTabParamList>;
  Subscription: undefined;
};
