// components/configurations/ProfileForm.tsx
 'use client';
import React, { useState } from 'react';
import { User } from '@/types';
import { updateUserProfile } from '@/lib/api';

interface ProfileFormProps {
  user?: User;
}

/**
 * Formulario para editar los detalles del perfil del usuario.
 * El email se muestra pero no se puede editar.
 */
const ProfileForm = ({ user }: ProfileFormProps) => {
  const userData = user;

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  if (!userData) {
    return <p>Loading user...</p>;
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setLoading(true);

    try {
      const form = new FormData(e.currentTarget);
      const payload = {
        email: userData.email,
        name: (form.get('name') as string) || '',
        lastName: (form.get('lastName') as string) || '',
        workPosition: (form.get('workPosition') as string) || '',
        salary: Number(form.get('salary') || 0),
      };

      if (!payload.name || !payload.lastName || !payload.workPosition) {
        setError('Please complete the profile fields.');
        setLoading(false);
        return;
      }

      await updateUserProfile(payload);
      setSuccess('Profile updated successfully.');
    } catch (err: any) {
      console.error('Failed to update profile', err);
      setError(err?.message || 'Failed to update profile');
    } finally {
      setLoading(false);
    }
  };

  return (
  <form onSubmit={handleSubmit} className="space-y-4">
      {/* Campo de Email (deshabilitado) */}
      <div>
        <label htmlFor="email" className="block text-sm font-medium text-gray-700">
          Email
        </label>
        <input
          id="email"
          type="email"
          value={userData.email}
          disabled
          className="mt-1 w-full px-4 py-2 bg-gray-100 border border-gray-300 rounded-md cursor-not-allowed"
        />
        <p className="text-xs text-gray-500 mt-1">El correo electrónico no se puede cambiar.</p>
      </div>

      {/* Campos Editables */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label htmlFor="name" className="block text-sm font-medium text-gray-700">
            Nombre
          </label>
          <input
            id="name"
            name="name"
            type="text"
            defaultValue={userData.name}
            className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>
        <div>
          <label htmlFor="lastName" className="block text-sm font-medium text-gray-700">
            Apellido
          </label>
          <input
            id="lastName"
            name="lastName"
            type="text"
            defaultValue={userData.lastName}
            className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label htmlFor="workPosition" className="block text-sm font-medium text-gray-700">
            Puesto de trabajo
          </label>
          <input
            id="workPosition"
            name="workPosition"
            type="text"
            defaultValue={userData.workPosition}
            className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>
        <div>
          <label htmlFor="salary" className="block text-sm font-medium text-gray-700">
            Salario
          </label>
          <input
            id="salary"
            name="salary"
            type="number"
            step="0.01"
            defaultValue={userData.salary}
            className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>
      </div>
      
      {/* Botón de Guardar */}
      <div className="text-right pt-2">
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-2 bg-primary text-white font-semibold rounded-lg hover:bg-blue-500 transition-colors shadow-md disabled:opacity-60"
        >
          {loading ? 'Guardando...' : 'Guardar Cambios'}
        </button>
        {error && <p className="text-sm text-red-600 mt-2">{error}</p>}
        {success && <p className="text-sm text-green-600 mt-2">{success}</p>}
      </div>
    </form>
  );
};

export default ProfileForm;