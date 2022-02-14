import { render, screen } from '@testing-library/react';
import App from './App';

test('chat room header exists', () => {
  render(<App />);
  const linkElement = screen.getByText(/chat room/i);
  expect(linkElement).toBeInTheDocument();
});
