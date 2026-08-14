import { render, screen } from '@testing-library/react';
import App from './App';

test('renders user management title', () => {
  render(<App />);
  const titleElement = screen.getByText(/user management/i);
  expect(titleElement).toBeInTheDocument();
});
