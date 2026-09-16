// components/configurations/PasswordForm.tsx
'use client';

import React, { useState } from 'react';
import { useCurrentUser } from '@/lib/userContext';
import { updatePassword } from "@/lib/api";

/**
 * Formulario para cambiar la contraseña del usuario.
 * Pide la contraseña actual y la nueva contraseña dos veces para confirmación.
 */
const PasswordForm = () => {
  const user = useCurrentUser();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  if (!user) {
    return <p>Loading user...</p>;
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setLoading(true);

    try {
      const form = new FormData(e.currentTarget);
      const currentPassword = form.get('currentPassword') as string;
      const newPassword = form.get('newPassword') as string;
      const confirmPassword = form.get('confirmPassword') as string;

      // Validaciones básicas
      if (!currentPassword || !newPassword || !confirmPassword) {
        setError('Por favor, completa todos los campos.');
        setLoading(false);
        return;
      }
      if (newPassword !== confirmPassword) {
        setError('La nueva contraseña y la confirmación no coinciden.');
        setLoading(false);
        return;
      }

      const response = await updatePassword({
        email: user.email,
        currentPassword,
        newPassword,
        confirmPassword,
      });

      setSuccess(response?.message || 'Contraseña actualizada con éxito.');
    } catch (error: any) {
      console.error('Error al actualizar la contraseña:', error);
      setError(error?.message || 'Error al actualizar la contraseña. Inténtalo de nuevo más tarde.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label
          htmlFor="currentPassword"
          className="block text-sm font-medium text-gray-700"
        >
          Contraseña Actual
        </label>
        <input
          id="currentPassword"
          name="currentPassword"
          type="password"
          className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
        />
      </div>
      <div>
        <label
          htmlFor="newPassword"
          className="block text-sm font-medium text-gray-700"
        >
          Nueva Contraseña
        </label>
        <input
          id="newPassword"
          name="newPassword"
          type="password"
          className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
        />
      </div>
      <div>
        <label
          htmlFor="confirmPassword"
          className="block text-sm font-medium text-gray-700"
        >
          Confirmar Nueva Contraseña
        </label>
        <input
          id="confirmPassword"
          name="confirmPassword"
          type="password"
          className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
        />
      </div>
      
      {/* Botón de Actualizar */}
      <div className="text-right pt-2">
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-2 bg-secondary text-white font-semibold rounded-lg hover:bg-opacity-90 transition-colors shadow-md"
        >
          {loading ? 'Actualizando...' : 'Actualizar Contraseña'}
        </button>
        {error && <p className="text-sm text-red-600 mt-2">{error}</p>}
        {success && <p className="text-sm text-green-600 mt-2">{success}</p>}
      </div>
    </form>
  );
};

export default PasswordForm;