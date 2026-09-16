'use client';
// app/(main)/configurations/page.tsx
import ProfileForm from '@/components/configurations/ProfileForm';
import PasswordForm from '@/components/configurations/PasswordForm';
import { useCurrentUser } from '@/lib/userContext';

/**
 * Página de "Configuraciones".
 * Permite al usuario modificar la información de su perfil y cambiar su contraseña.
 * Está dividida en secciones para facilitar la gestión.
 */
export default function ConfigurationsPage() {
  const user = useCurrentUser();
  return (
    <div className="container mx-auto max-w-3xl">
      {/* Encabezado */}
      <h1 className="text-3xl font-bold text-dark-background mb-6">
        Configuración de la Cuenta
      </h1>

      <div className="space-y-8">
        {/* Sección de Información del Perfil */}
        <div className="bg-white p-6 rounded-lg shadow-lg">
          <h2 className="text-xl font-bold text-dark-background border-b pb-3 mb-4">
            Información del Perfil
          </h2>
          <ProfileForm user={user ?? undefined} />
        </div>

        {/* Sección para Cambiar Contraseña */}
        <div className="bg-white p-6 rounded-lg shadow-lg">
          <h2 className="text-xl font-bold text-dark-background border-b pb-3 mb-4">
            Cambiar Contraseña
          </h2>
          <PasswordForm />
        </div>
      </div>
    </div>
  );
}