'use client';
import { LogEntry } from '@/types';

interface LogTableProps {
    logs: LogEntry[];
}

const LogTable = ({ logs }: LogTableProps) => {
    return (
        <div className="bg-white p-4 rounded-lg shadow-lg">
            <h2 className="font-bold text-lg text-dark-background mb-4">Logs</h2>
            <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                    <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
                    </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                    {logs.length === 0 ? (
                        <tr>
                            <td className="px-6 py-4 whitespace-nowrap" colSpan={3}>
                                No logs available.
                            </td>
                        </tr>
                    ) : (
                        logs.map((log) => (
                            <tr key={log.id}>
                                <td className="px-6 py-4 whitespace-nowrap">{log.timestamp}</td>
                                <td className="px-6 py-4 whitespace-nowrap">{log.username}</td>
                                <td className="px-6 py-4 whitespace-nowrap">{log.action}</td>
                            </tr>
                        ))
                    )}
                </tbody>
            </table>
        </div>
    );
};

export default LogTable;