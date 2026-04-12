export interface Profile {
  id: string;
  name: string;
  type: 'human' | 'ai';
  imageUrl: string;
  prompt?: string;
  relationshipLevel?: string;
}

export const mockProfiles: Profile[] = [
  { id: '1', name: 'Alice', type: 'human', imageUrl: 'https://placehold.co/400', prompt: 'I am looking for a real connection.' },
  { id: '2', name: 'The Intellectual', type: 'ai', imageUrl: 'https://placehold.co/400', relationshipLevel: 'Acquaintance' },
  { id: '3', name: 'Bob', type: 'human', imageUrl: 'https://placehold.co/400', prompt: 'Coffee addict and dog lover.' },
  { id: '4', name: 'The Adventurer', type: 'ai', imageUrl: 'https://placehold.co/400', relationshipLevel: 'Friend' },
  { id: '5', name: 'Charlie', type: 'human', imageUrl: 'https://placehold.co/400', prompt: 'Pizza is my favorite vegetable.' },
  { id: '6', name: 'The Romantic', type: 'ai', imageUrl: 'https://placehold.co/400', relationshipLevel: 'Partner' },
  { id: '7', name: 'Diana', type: 'human', imageUrl: 'https://placehold.co/400', prompt: 'Always exploring new places.' },
  { id: '8', name: 'The Mentor', type: 'ai', imageUrl: 'https://placehold.co/400', relationshipLevel: 'Guide' },
  { id: '9', name: 'Eve', type: 'human', imageUrl: 'https://placehold.co/400', prompt: 'Movie buff.' },
  { id: '10', name: 'The Jester', type: 'ai', imageUrl: 'https://placehold.co/400', relationshipLevel: 'Entertainer' },
];
