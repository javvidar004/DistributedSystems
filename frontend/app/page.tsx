// app/(auth)/signin/page.tsx
"use client";

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { userSignIn } from '../lib/api';

/**
 * Página de inicio de sesión (Sign In).
 * Muestra un formulario para que el usuario ingrese con su correo/usuario y contraseña.
 * Está diseñada para ser usada dentro del AuthLayout, que la centra en la pantalla.
 */
export default function SignInPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const form = new FormData(e.currentTarget);
      const credentials = {
        email: (form.get('email') as string) || '',
        password: (form.get('password') as string) || '',
      };

      if (!credentials.email || !credentials.password) {
        setError('Please enter your email and password.');
        setLoading(false);
        return;
      }

      // Call the shared API helper which stores the userId cookie.
      const resp = await userSignIn(credentials);

      if (resp?.user?.id) {
        router.push('/dashboard');
      } else {
        setError('Authentication failed: no user returned');
      }
    } catch (err: any) {
      console.error('Sign in error', err);
      // Try to extract friendly message
      const msg = err?.response?.data?.message || err?.message || 'Sign in failed';
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      {/* Título del formulario */}
      <h2 className="text-2xl font-bold text-center text-dark-background">
        Login
      </h2>
      
      {/* Formulario de inicio de sesión */}
      <form onSubmit={handleSubmit} className="mt-8 space-y-6">
        {/* Campo para Email o Nombre de Usuario */}
        <div>
          <label htmlFor="email" className="sr-only">Email</label>
          <input
            id="email"
            name="email"
            type="email"
            autoComplete="email"
            required
            className="w-full px-4 py-3 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-shadow"
            placeholder="email/username"
          />
        </div>
        
        {/* Campo para Contraseña */}
        <div>
          <label htmlFor="password" className="sr-only">Password</label>
          <input
            id="password"
            name="password"
            type="password"
            autoComplete="current-password"
            required
            className="w-full px-4 py-3 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-shadow"
            placeholder="password"
          />
        </div>
        
        {/* Enlace para recuperar contraseña */}
        <div className="text-sm text-right">
          <Link href="/forgot-password" className="font-medium text-primary hover:underline">
            Forgot password?
          </Link>
        </div>
        
        {/* Botón de envío */}
        <button
          type="submit"
          disabled={loading}
          className="w-full py-3 px-4 bg-secondary text-white font-semibold rounded-md hover:bg-opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-secondary transition-transform transform hover:scale-105 disabled:opacity-60"
        >
          {loading ? 'Signing in...' : 'Log in'}
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