import React from 'react';
import { render } from '@testing-library/react-native';
import ChatScreen from '../components/ChatScreen';
import { Provider } from 'react-redux';
import { store } from '../../store/store';

describe('ChatScreen', () => {
  it('renders correctly', () => {
    const { getByText } = render(
      <Provider store={store}>
        <ChatScreen />
      </Provider>
    );

    expect(getByText(/Chat with AI/i)).toBeTruthy();
  });
});
