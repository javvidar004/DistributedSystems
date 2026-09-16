// app/(auth)/signup/page.tsx
"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { createUser } from '@/lib/api';

/**
 * Página para crear usuarios desde el panel autenticado.
 * Envía la estructura que espera el servicio de users.
 */
export default function CreateUserPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const form = new FormData(e.currentTarget);
      const payload = {
        name: (form.get('name') as string) || '',
        lastName: (form.get('lastName') as string) || '',
        email: (form.get('email') as string) || '',
        workPosition: (form.get('workPosition') as string) || '',
        salary: Number(form.get('salary') || 0),
      };

      if (!payload.email || !payload.name || !payload.lastName || !payload.workPosition) {
        setError('Please fill in all required fields.');
        setLoading(false);
        return;
      }

      await createUser(payload);
      router.push('/users');
    } catch (err) {
      console.error('Signup error', err);
      setError((err as Error).message || 'An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      {/* Título del formulario */}
      <h2 className="text-2xl font-bold text-center text-dark-background">
        Create User
      </h2>

      {/* Formulario de registro */}
      <form onSubmit={handleSubmit} className="mt-8 space-y-4">
        <div className="flex flex-col sm:flex-row sm:space-x-4 space-y-4 sm:space-y-0">
          {/* Campo de Nombre */}
          <div className="w-full">
            <label htmlFor="name" className="sr-only">First Name</label>
            <input
              id="name"
              name="name"
              type="text"
              required
              className="w-full px-4 py-3 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-shadow"
              placeholder="First Name"
            />
          </div>
          {/* Campo de Apellido */}
          <div className="w-full">
            <label htmlFor="lastName" className="sr-only">Last Name</label>
            <input
              id="lastName"
              name="lastName"
              type="text"
              required
              className="w-full px-4 py-3 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-shadow"
              placeholder="Last Name"
            />
          </div>
        </div>

        <div className="flex flex-col sm:flex-row sm:space-x-4 space-y-4 sm:space-y-0">
          {/* Campo de Puesto */}
          <div className="w-full">
            <label htmlFor="workPosition" className="sr-only">Work Position</label>
            <input
              id="workPosition"
              name="workPosition"
              type="text"
              required
              className="w-full px-4 py-3 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-shadow"
              placeholder="Work Position"
            />
          </div>
          {/* Campo de Salario */}
          <div className="w-full">
            <label htmlFor="salary" className="sr-only">Salary</label>
            <input
              id="salary"
              name="salary"
              type="number"
              step="0.01"
              required
              className="w-full px-4 py-3 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-shadow"
              placeholder="Salary"
            />
          </div>
        </div>

        {/* Campo de Email */}
        <div>
          <label htmlFor="email" className="sr-only">Email</label>
          <input
            id="email"
            name="email"
            type="email"
            autoComplete="email"
            required
            className="w-full px-4 py-3 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-shadow"
            placeholder="Email address"
          />
        </div>

        {/* Campo de Contraseña */}
        <div>
          <label htmlFor="password" className="sr-only">Password</label>
          <input
            id="password"
            name="password"
            type="password"
            required
            className="w-full px-4 py-3 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-shadow"
            placeholder="Password"
          />
        </div>

        {/* Botón de envío */}
        <button
          type="submit"
          disabled={loading}
          className="w-full py-3 px-4 bg-secondary text-white font-semibold rounded-md hover:bg-opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-secondary transition-transform transform hover:scale-105 disabled:opacity-60"
        >
          {loading ? 'Creating account...' : 'Sign Up'}
        </button>

        {error && (
          <p className="text-sm text-red-600 text-center" role="alert">
            {error}
          </p>
        )}
      </form>
    </>
  );
}