"use client";

import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { PlusIcon } from '@heroicons/react/24/solid';
import UserTable from '@/components/users/UserTable';
import { getUsers } from '@/lib/api';
import { User } from '@/types';

/**
 * Página de usuarios.
 * Muestra las opciones principales: crear nuevo usuario y buscar.
 * También renderiza la lista de usuarios obtenida del servicio.
 */
export default function UsersPage() {
  const { data: allUsers = [], isLoading, isError } = useQuery<User[]>({
    queryKey: ['users'],
    queryFn: getUsers,
  });


  return (
    <div className="container mx-auto">
      {/* Encabezado y Acciones Principales */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8 gap-4">
        <h1 className="text-3xl font-bold text-dark-background">
          Users
        </h1>
        <div className="flex items-center gap-2">
          <Link 
            href="/users/new"
            className="flex items-center justify-center gap-2 px-4 py-2 bg-primary text-white font-semibold rounded-lg hover:bg-blue-500 transition-colors shadow-md"
          >
            <PlusIcon className="h-5 w-5" />
            <span>New User</span>
          </Link>
        </div>
      </div>

      {/* Lista de Usuarios */}
      {isLoading ? (
        <p className="text-gray-600">Loading users...</p>
      ) : isError ? (
        <p className="text-red-600">Could not load users.</p>
      ) : (
        <UserTable users={allUsers ?? []} />
      )}
    </div>
  );
}