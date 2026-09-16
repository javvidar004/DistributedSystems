// components/shopping-list/IngredientList.tsx
import { Ingredient } from '@/types';
import { CheckCircleIcon } from '@heroicons/react/24/outline';

interface IngredientListProps {
  ingredients: Ingredient[];
}

/**
 * Muestra la lista de ingredientes generada, permitiendo al usuario
 * marcar los que ya tiene.
 */
const IngredientList = ({ ingredients }: IngredientListProps) => {
  return (
    <div className="bg-white p-6 rounded-lg shadow-lg">
      <h2 className="text-2xl font-bold text-dark-background mb-4">Ingredientes a comprar:</h2>
      <ul className="space-y-3">
        {ingredients.map((ingredient) => (
          <li
            key={ingredient.id}
            className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <div className="flex items-center gap-4">
              <button className="group">
                <CheckCircleIcon className="h-7 w-7 text-gray-300 group-hover:text-green-500 transition-colors" />
              </button>
              <span className="font-medium text-gray-800">{ingredient.name}</span>
            </div>
            <span className="text-gray-600 font-semibold">
              {`${ingredient.quantity} ${ingredient.units}`}
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default IngredientList;