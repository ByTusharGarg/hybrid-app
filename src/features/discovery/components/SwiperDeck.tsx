import React, { useState } from 'react';
import { View, Text, TouchableOpacity, Dimensions } from 'react-native';
import { Gesture, GestureDetector } from 'react-native-gesture-handler';
import Animated, { 
  useSharedValue, 
  useAnimatedStyle, 
  withSpring, 
  withTiming, 
  runOnJS,
  interpolate,
  Extrapolation
} from 'react-native-reanimated';
import { Profile } from '../../../shared/mock/mockProfiles';

interface SwiperDeckProps {
  profiles: Profile[];
}

const { width: SCREEN_WIDTH } = Dimensions.get('window');
const SWIPE_THRESHOLD = SCREEN_WIDTH * 0.3;

export function SwiperDeck({ profiles }: SwiperDeckProps) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [history, setHistory] = useState<number[]>([]);

  // The active card profile
  const currentProfile = profiles[currentIndex];
  // Next profile for the background visually
  const nextProfile = profiles[currentIndex + 1];

  const translateX = useSharedValue(0);
  const translateY = useSharedValue(0);

  const performSwipe = (direction: 'left' | 'right' | 'up') => {
    let toX = 0;
    let toY = 0;
    
    if (direction === 'left') toX = -SCREEN_WIDTH * 1.5;
    if (direction === 'right') toX = SCREEN_WIDTH * 1.5;
    if (direction === 'up') toY = -SCREEN_WIDTH * 1.5;

    translateX.value = withTiming(toX, { duration: 300 });
    translateY.value = withTiming(toY, { duration: 300 }, () => {
      runOnJS(finalizeSwipe)();
    });
  };

  const finalizeSwipe = () => {
    setHistory(prev => [...prev, currentIndex]);
    setCurrentIndex(prev => prev + 1);
    translateX.value = 0;
    translateY.value = 0;
  };

  const handleRewind = () => {
    if (history.length > 0) {
      const lastIndex = history[history.length - 1];
      setHistory(prev => prev.slice(0, -1));
      
      // Move card to previous position instantly
      translateX.value = -SCREEN_WIDTH;
      translateY.value = 0;
      
      setCurrentIndex(lastIndex);
      
      // Animate it back to center
      translateX.value = withSpring(0);
    }
  };

  const panGesture = Gesture.Pan()
    .onChange((event) => {
      translateX.value += event.changeX;
      translateY.value += event.changeY;
    })
    .onEnd((event) => {
      if (translateX.value > SWIPE_THRESHOLD) {
        runOnJS(performSwipe)('right');
      } else if (translateX.value < -SWIPE_THRESHOLD) {
        runOnJS(performSwipe)('left');
      } else if (translateY.value < -SWIPE_THRESHOLD && Math.abs(translateX.value) < SWIPE_THRESHOLD) {
        runOnJS(performSwipe)('up');
      } else {
        translateX.value = withSpring(0);
        translateY.value = withSpring(0);
      }
    });

  const animatedCardStyle = useAnimatedStyle(() => {
    const rotateZ = interpolate(
      translateX.value,
      [-SCREEN_WIDTH / 2, 0, SCREEN_WIDTH / 2],
      [-15, 0, 15],
      Extrapolation.CLAMP
    );

    return {
      transform: [
        { translateX: translateX.value },
        { translateY: translateY.value },
        { rotateZ: `${rotateZ}deg` }
      ]
    };
  });

  const animatedNextCardStyle = useAnimatedStyle(() => {
    const scale = interpolate(
      Math.abs(translateX.value),
      [0, SCREEN_WIDTH / 2],
      [0.9, 1],
      Extrapolation.CLAMP
    );
    return {
      transform: [{ scale }],
      opacity: scale
    };
  });

  if (currentIndex >= profiles.length) {
    return (
      <View className="flex-1 justify-center items-center">
        <Text className="text-white font-bold text-xl mb-4">Out of Profiles!</Text>
        <TouchableOpacity onPress={handleRewind} className="bg-stone-800 px-6 py-3 rounded-full">
          <Text className="text-pink-500 font-bold">⏪ Rewind</Text>
        </TouchableOpacity>
      </View>
    );
  }

  const renderCard = (profile: Profile, isTopCard: boolean) => {
    return (
      <View className="absolute inset-0 m-4 rounded-3xl overflow-hidden bg-stone-900 border border-stone-800 shadow-2xl">
         <View className="flex-1 items-center justify-center bg-stone-800">
           <Text className="text-6xl">📸</Text>
         </View>
         <View className="p-6 bg-stone-900">
           <Text className="text-3xl font-extrabold text-white mb-2">{profile.name}</Text>
           {profile.prompt && <Text className="text-stone-400 text-lg mb-4">{profile.prompt}</Text>}
         </View>
      </View>
    );
  };

  return (
    <View className="flex-1">
       {/* Deck Area */}
       <View className="flex-1 relative">
         {nextProfile && (
           <Animated.View style={[{ position: 'absolute', top: 0, left: 0, right: 0, bottom: 0 }, animatedNextCardStyle]}>
             {renderCard(nextProfile, false)}
           </Animated.View>
         )}
         {currentProfile && (
           <GestureDetector gesture={panGesture}>
             <Animated.View style={[{ position: 'absolute', top: 0, left: 0, right: 0, bottom: 0 }, animatedCardStyle]}>
               {renderCard(currentProfile, true)}
             </Animated.View>
           </GestureDetector>
         )}
       </View>

       {/* Controls Area */}
       <View className="flex-row justify-evenly items-center py-6">
         <TouchableOpacity 
           onPress={handleRewind}
           disabled={history.length === 0}
           className={`w-14 h-14 rounded-full items-center justify-center ${history.length === 0 ? 'bg-stone-800/50' : 'bg-stone-800 border border-stone-700'}`}
         >
           <Text className="text-2xl opacity-80">⏪</Text>
         </TouchableOpacity>

         <TouchableOpacity 
           onPress={() => performSwipe('left')}
           className="w-16 h-16 rounded-full bg-stone-800 border border-stone-700 items-center justify-center shadow-lg"
         >
           <Text className="text-3xl">❌</Text>
         </TouchableOpacity>

         <TouchableOpacity 
           onPress={() => performSwipe('up')}
           className="w-14 h-14 rounded-full bg-indigo-500 items-center justify-center shadow-lg shadow-indigo-500/30"
         >
           <Text className="text-2xl">🌟</Text>
         </TouchableOpacity>

         <TouchableOpacity 
           onPress={() => performSwipe('right')}
           className="w-16 h-16 rounded-full bg-pink-500 items-center justify-center shadow-lg shadow-pink-500/30"
         >
           <Text className="text-3xl">❤️</Text>
         </TouchableOpacity>
       </View>
    </View>
  );
}
