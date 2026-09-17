// components/users/UserCard.tsx
import { User } from '@/types';

interface UserCardProps {
  user: User;
}

/**
 * Tarjeta individual para mostrar un resumen de un usuario.
 */
const UserCard = ({ user }: UserCardProps) => {
  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden transform hover:-translate-y-1 transition-transform duration-300">
      <div className="p-6">
        <h3 className="text-xl font-bold text-dark-background truncate">{user.name} {user.lastName}</h3>
        <p className="text-gray-600 mt-2 truncate">{user.email}</p>
        <p className="text-gray-600 mt-1">{user.workPosition}</p>
        <p className="text-gray-600 mt-1 font-medium">${user.salary.toFixed(2)}</p>
      </div>
    </div>
  );
};

export default UserCard;