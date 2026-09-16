// components/recipes/RecipeList.tsx
import { User } from '@/types';
import RecipeCard from './RecipeCard';

interface RecipeListProps {
  recipes: User[];
}

/**
 * Muestra una lista de usuarios en formato de cuadrícula.
 */
const RecipeList = ({ recipes }: RecipeListProps) => {
  if (recipes.length === 0) {
    return (
      <div className="text-center py-12 bg-white rounded-lg shadow-md">
        <h2 className="text-xl font-semibold text-gray-700">No se encontraron usuarios</h2>
        <p className="text-gray-500 mt-2">Intenta con otro término de búsqueda o crea un nuevo usuario.</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {recipes.map((recipe) => (
        // Aquí podrías usar una tarjeta diferente para los resultados si quisieras,
        // por ejemplo, una que no tenga los botones de "Editar" y "Borrar".
        // Por ahora, reutilizamos la misma.
        <RecipeCard key={recipe.id} recipe={recipe} />
      ))}
    </div>
  );
};

export default RecipeList;