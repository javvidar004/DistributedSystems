// components/recipes/RecipeCard.tsx
import { User } from '@/types';

interface RecipeCardProps {
  recipe: User;
}

/**
 * Tarjeta individual para mostrar un resumen de un usuario.
 */
const RecipeCard = ({ recipe }: RecipeCardProps) => {
  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden transform hover:-translate-y-1 transition-transform duration-300">
      <div className="p-6">
        <h3 className="text-xl font-bold text-dark-background truncate">{recipe.name} {recipe.lastName}</h3>
        <p className="text-gray-600 mt-2 truncate">{recipe.email}</p>
        <p className="text-gray-600 mt-1">{recipe.workPosition}</p>
        <p className="text-gray-600 mt-1 font-medium">${recipe.salary.toFixed(2)}</p>
      </div>
    </div>
  );
};

export default RecipeCard;