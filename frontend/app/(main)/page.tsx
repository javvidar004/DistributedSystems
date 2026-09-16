// app/(main)/page.tsx
'use client';
import WelcomeHeader from '@/components/dashboard/WelcomeHeader';
import QuickActions from '@/components/dashboard/QuickActions';
import { useCurrentUser } from '@/lib/userContext';

/**
 * Homepage para usuarios autenticados (Dashboard).
 * Muestra un saludo, acciones rápidas, un vistazo a la próxima comida
 * y una vista previa de las recetas favoritas.
 */
export default function DashboardPage() {
  const userData = useCurrentUser();

  return (
    <div className="container mx-auto space-y-8">
      {/* 1. Saludo de Bienvenida */}
      <WelcomeHeader name={userData?.name ?? ''} />

      {/* 2. Acciones Rápidas */}
      <QuickActions />
    </div>
  );
}