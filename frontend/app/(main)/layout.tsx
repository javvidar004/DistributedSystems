// app/(main)/layout.tsx
'use client';
import { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import Navbar from '@/components/layout/Navbar';
import LeftSidebar from '@/components/layout/LeftSidebar';

import { getUserData } from '@/lib/api';
import { UserProvider } from '@/lib/userContext';

/**
 * Main layout for authenticated users.
 * It manages the state of both sidebars and provides them along with the navbar
 * to all child pages within this route group.
 */
export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [isLeftSidebarOpen, setLeftSidebarOpen] = useState(false);

  const toggleLeftSidebar = () => setLeftSidebarOpen(!isLeftSidebarOpen);

  const { data: userData, isError } = useQuery({
    queryKey: ['current-user'],
    queryFn: getUserData,
    staleTime: 5 * 60 * 1000,
    refetchOnWindowFocus: false,
    retry: false,
  });

  const router = useRouter();

  useEffect(() => {
    if (isError) {
      router.push('/');
    }
  }, [isError, router]);

  return (
    <div className="min-h-screen bg-background text-gray-800">
      <Navbar
        onToggleLeftSidebar={toggleLeftSidebar}
      />
      <LeftSidebar isOpen={isLeftSidebarOpen} onClose={() => setLeftSidebarOpen(false)} />

      <UserProvider value={userData ?? null}>
        <main className="pt-16 transition-all duration-300">
          <div className="container mx-auto p-4 lg:p-6">
            {children}
          </div>
        </main>
      </UserProvider>
    </div>
  );
}