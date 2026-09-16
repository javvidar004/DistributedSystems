// app/(main)/search/page.tsx
'use client';

import { useState } from 'react';
import { User } from '@/types';
import { useQuery } from '@tanstack/react-query'; 
import { getUsers } from '@/lib/api';
import SearchBar from '@/components/search/SearchBar';
import RecipeList from '@/components/recipes/RecipeList'; // Reutilizamos el componente de la lista

/**
 * Página de búsqueda de usuarios.
 * Filtra localmente los registros obtenidos del servicio de users.
 */
export default function SearchPage() {
  const [results, setResults] = useState<User[]>([]);
  const [hasSearched, setHasSearched] = useState(false);

  const { data: allUsers = [], isLoading, isError } = useQuery<User[]>({
    queryKey: ['users'],
    queryFn: getUsers,
  });

  const handleSearch = (query: string) => {
    setHasSearched(true);
    if (!query) {
      setResults([]);
      return;
    }

    const lowercasedQuery = query.toLowerCase();
    const filteredResults = (allUsers ?? []).filter(user => 
      user.name.toLowerCase().includes(lowercasedQuery) || 
      user.email.toLowerCase().includes(lowercasedQuery) ||
      user.lastName.toLowerCase().includes(lowercasedQuery) ||
      user.workPosition.toLowerCase().includes(lowercasedQuery)
    );

    setResults(filteredResults);
  };

  return (
    <div className="container mx-auto">
      {/* Encabezado de la página */}
      <h1 className="text-3xl font-bold text-dark-background mb-4">
        Buscar Usuarios
      </h1>
      <p className="text-gray-600 mb-6">
        Encuentra usuarios por nombre, apellido, correo o puesto de trabajo.
      </p>

      {/* Barra de Búsqueda */}
      <SearchBar onSearch={handleSearch} />

      {/* Sección de Resultados */}
      <div className="mt-8">
        {isLoading ? (
          <p className="text-gray-600">Loading users...</p>
        ) : isError ? (
          <p className="text-red-600">Could not load users.</p>
        ) : hasSearched ? (
          <RecipeList recipes={results} />
        ) : (
          <div className="text-center py-12 bg-white rounded-lg shadow-md">
            <h2 className="text-xl font-semibold text-gray-700">Comienza tu búsqueda</h2>
            <p className="text-gray-500 mt-2">Escribe algo en la barra de búsqueda para ver resultados.</p>
          </div>
        )}
      </div>
    </div>
  );
}