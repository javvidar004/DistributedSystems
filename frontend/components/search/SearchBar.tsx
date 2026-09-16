// components/search/SearchBar.tsx
'use client';
import { useState } from 'react';
import { MagnifyingGlassIcon } from '@heroicons/react/24/solid';

interface SearchBarProps {
  onSearch: (query: string) => void;
}

/**
 * Componente reutilizable para la barra de búsqueda.
 * Contiene un campo de texto y un botón para iniciar la búsqueda.
 */
const SearchBar = ({ onSearch }: SearchBarProps) => {
  const [query, setQuery] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault(); // Evita que la página se recargue
    onSearch(query);
  };

  return (
    <form onSubmit={handleSubmit} className="flex items-center gap-2">
      <input
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Buscar por nombre, apellido o correo..."
        className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-shadow shadow-sm"
      />
      <button
        type="submit"
        className="flex items-center justify-center px-4 py-3 bg-secondary text-white font-semibold rounded-lg hover:bg-opacity-90 transition-colors shadow-sm"
        aria-label="Buscar"
      >
        <MagnifyingGlassIcon className="h-6 w-6" />
      </button>
    </form>
  );
};

export default SearchBar;