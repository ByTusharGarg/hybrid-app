export interface Product {
  id: string;
  name: string;
  price: number;
  imageUrl: string;
}

export const mockProducts: Product[] = [
  { id: 'p1', name: 'Virtual Rose', price: 10, imageUrl: 'https://placehold.co/100' },
  { id: 'p2', name: 'Premium Gift Box', price: 50, imageUrl: 'https://placehold.co/100' },
  { id: 'p3', name: 'Concert Ticket', price: 100, imageUrl: 'https://placehold.co/100' },
];
