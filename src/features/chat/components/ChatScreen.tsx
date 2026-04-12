import React, { useState, useEffect } from 'react';
import { View, Text, TextInput, TouchableOpacity, Alert, ScrollView } from 'react-native';
import { useSelector, useDispatch } from 'react-redux';
import { RootState } from '../../store/store';
import { deductCredits } from '../../store/slices/appSlice';
import { mockProducts, Product } from '../../../shared/mock/mockProducts';
import { useNavigation, useRoute } from '@react-navigation/native';
import { mockProfiles } from '../../../shared/mock/mockProfiles';

interface Message {
  id: string;
  sender: 'user' | 'ai';
  text: string;
  isProduct?: boolean;
  productInfo?: Product;
}

export default function ChatScreen() {
  const route = useRoute<any>();
  const navigation = useNavigation();
  const matchId = route.params?.matchId;
  const currentProfile = mockProfiles.find(p => p.id === matchId);
  const dispatch = useDispatch();
  const credits = useSelector((state: RootState) => state.app.credits);

  const [messages, setMessages] = useState<Message[]>([
    { id: '1', sender: 'ai', text: `Hello! It is so nice to finally chat with you, I'm ${currentProfile?.name || 'Aura'}. I love gifts!` },
  ]);
  const [inputText, setInputText] = useState('');
  const [isTyping, setIsTyping] = useState(false);

  useEffect(() => {
    const lastMsg = messages[messages.length - 1];
    if (lastMsg && lastMsg.sender === 'user' && !lastMsg.isProduct) {
      if (lastMsg.text.toLowerCase().includes('gift')) {
        // AI Suggests a gift
        setTimeout(() => {
          setMessages(prev => [
            ...prev,
            { id: Date.now().toString(), sender: 'ai', text: 'I would love a gift! How about this?' },
            { id: Date.now().toString() + '_p', sender: 'ai', text: '', isProduct: true, productInfo: mockProducts[0] }
          ]);
        }, 1000);
      } else {
        // Normal generic response
        setIsTyping(true);
        setTimeout(() => {
          setMessages(prev => [
            ...prev,
            { id: Date.now().toString(), sender: 'ai', text: 'That sounds amazing! Tell me more.' }
          ]);
          setIsTyping(false);
        }, 1500);
      }
    }
  }, [messages]);

  const handleSend = () => {
    if (inputText.trim()) {
      setMessages([...messages, { id: Date.now().toString(), sender: 'user', text: inputText }]);
      setInputText('');
    }
  };

  const handleBuy = (product: Product) => {
    if (credits >= product.price) {
      Alert.alert(
        'Purchase Confirmation',
        `Buy ${product.name} for ${product.price} credits?`,
        [
          { text: 'Cancel', style: 'cancel' },
          { 
            text: 'Buy', 
            onPress: () => {
              dispatch(deductCredits(product.price));
              Alert.alert('Success', `You purchased ${product.name}! Your credits balance is now ${credits - product.price}.`);
              setMessages(prev => [...prev, { id: Date.now().toString(), sender: 'user', text: `Here is a ${product.name} for you!` }]);
            } 
          }
        ]
      );
    } else {
      Alert.alert('Out of Credits', 'Please top-up your credits to buy this gift.');
    }
  };

  return (
    <View className="flex-1 bg-stone-950 pt-16">
      <View className="flex-row items-center p-4 border-b border-stone-800 bg-stone-900">
        <TouchableOpacity onPress={() => navigation.goBack()} className="mr-3 bg-stone-800 p-2 rounded-full">
          <Text className="text-white font-bold px-2">←</Text>
        </TouchableOpacity>
        <Text className="text-xl font-bold flex-1 text-center text-white">{currentProfile?.name || 'Aura'}</Text>
        <Text className="text-sm rounded-full bg-pink-500/20 text-pink-500 border border-pink-500/30 px-3 py-1 font-bold">
          🪙 {credits}
        </Text>
      </View>

      <ScrollView className="flex-1 px-4 pt-4">
        {messages.map(msg => (
          <View key={msg.id} className={`mb-4 w-3/4 ${msg.sender === 'user' ? 'self-end' : 'self-start'}`}>
            {!msg.isProduct ? (
              <View className={`p-4 rounded-2xl ${msg.sender === 'user' ? 'bg-indigo-500 rounded-tr-none' : 'bg-stone-800 shadow-sm border border-stone-700 rounded-tl-none'}`}>
                <Text className={msg.sender === 'user' ? 'text-white' : 'text-stone-200'}>{msg.text}</Text>
              </View>
            ) : (
              <View className="bg-stone-900 p-4 rounded-2xl shadow-xl border border-pink-500/20">
                <View className="h-32 bg-stone-800 rounded-xl mb-3 flex items-center justify-center">
                  <Text className="text-5xl">🎁</Text>
                </View>
                <Text className="font-bold text-lg text-white mb-2">{msg.productInfo?.name}</Text>
                <TouchableOpacity 
                  className={`py-3 rounded-xl items-center shadow-lg ${credits >= msg.productInfo!.price ? 'bg-pink-500 shadow-pink-500/30' : 'bg-stone-700'}`}
                  onPress={() => handleBuy(msg.productInfo!)}
                >
                  <Text className="text-white font-bold">Buy for {msg.productInfo?.price} credits</Text>
                </TouchableOpacity>
              </View>
            )}
          </View>
        ))}
        {isTyping && <Text className="text-stone-500 italic ml-2">Aura is typing...</Text>}
      </ScrollView>

      <View className="p-4 bg-stone-900 border-t border-stone-800 flex-row items-center pb-8">
        <TextInput
          className="flex-1 bg-stone-800 border border-stone-700 text-white rounded-full px-5 py-3 mr-3"
          placeholder="Type 'gift' to see magic..."
          placeholderTextColor="#78716c"
          value={inputText}
          onChangeText={setInputText}
          onSubmitEditing={handleSend}
        />
        <TouchableOpacity className="bg-indigo-500 rounded-full p-4 px-6 shadow-lg shadow-indigo-500/30" onPress={handleSend}>
          <Text className="text-white font-bold">Send</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}
