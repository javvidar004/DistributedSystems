// app/(main)/logs/page.tsx
'use client';

import { useQuery } from '@tanstack/react-query';
import LogTable from '@/components/logs/LogTable';
import { getLogs } from '@/lib/api';

/**
 * Página principal de logs.
 * Consume la lista del servicio de logs y la muestra en tabla.
 */
export default function LogsPage() {
  const { data: logs = [], isLoading, isError } = useQuery({
    queryKey: ['logs'],
    queryFn: getLogs,
  });

  return (
    <div className="container mx-auto">
      <h1 className="text-3xl font-bold text-dark-background mb-6">Logs</h1>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-3">
          {isLoading ? (
            <p className="text-gray-600">Loading logs...</p>
          ) : isError ? (
            <p className="text-red-600">Could not load logs.</p>
          ) : (
            <LogTable logs={logs} />
          )}
        </div>

      </div>
    </div>
  );
}