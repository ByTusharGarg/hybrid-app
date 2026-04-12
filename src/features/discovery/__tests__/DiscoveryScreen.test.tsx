import React from 'react';
import { render, fireEvent } from '@testing-library/react-native';
import DiscoveryScreen from '../components/DiscoveryScreen';
import { Provider } from 'react-redux';
import { store } from '../../store/store';

describe('DiscoveryScreen', () => {
  it('toggles between Real Users and AI Personas', () => {
    const { getByText, queryByText } = render(
      <Provider store={store}>
        <DiscoveryScreen />
      </Provider>
    );

    // Initial state should be "Real Users"
    expect(getByText('Alice')).toBeTruthy(); // A human mock
    expect(queryByText('The Intellectual')).toBeNull(); // An AI mock

    // Toggle to AI Personas
    const toggleButton = getByText('AI Personas');
    fireEvent.press(toggleButton);

    expect(getByText('The Intellectual')).toBeTruthy();
    expect(queryByText('Alice')).toBeNull();
  });
});
