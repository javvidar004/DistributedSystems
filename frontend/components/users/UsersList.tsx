// components/users/UsersList.tsx
import { User } from '@/types';
import UserCard from './UserCard';

interface UsersListProps {
  users: User[];
}

/**
 * Muestra una lista de usuarios en formato de cuadrícula.
 */
const UsersList = ({ users }: UsersListProps) => {
  if (users.length === 0) {
    return (
      <div className="text-center py-12 bg-white rounded-lg shadow-md">
        <h2 className="text-xl font-semibold text-gray-700">No se encontraron usuarios</h2>
        <p className="text-gray-500 mt-2">Intenta con otro término de búsqueda o crea un nuevo usuario.</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {users.map((user) => (
        // Aquí podrías usar una tarjeta diferente para los resultados si quisieras,
        // por ejemplo, una que no tenga los botones de "Editar" y "Borrar".
        // Por ahora, reutilizamos la misma.
        <UserCard user={user} />
      ))}
    </div>
  );
};

export default UsersList;