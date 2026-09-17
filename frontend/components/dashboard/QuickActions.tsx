// components/dashboard/QuickActions.tsx
import Link from 'next/link';
import { UsersIcon, PlusCircleIcon, MagnifyingGlassIcon, ClipboardDocumentListIcon } from '@heroicons/react/24/solid';

const QuickActions = () => {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
      <Link href="/users" className="p-8 bg-primary text-white rounded-lg shadow-lg hover:bg-blue-500 transition-colors flex items-center gap-6">
        <UsersIcon className="h-10 w-10" />
        <div>
          <h2 className="font-bold text-2xl">View Users</h2>
          <p className="text-base opacity-90">Manage the users of the system.</p>
        </div>
      </Link>

      <Link href="/logs" className="p-8 bg-secondary text-white rounded-lg shadow-lg hover:bg-opacity-90 transition-colors flex items-center gap-6">
        <ClipboardDocumentListIcon className="h-10 w-10" />
        <div>
          <h2 className="font-bold text-2xl">View logs</h2>
          <p className="text-base opacity-90">Review the actions performed in the system.</p>
        </div>
      </Link>

      <Link href="/users/new" className="p-8 bg-green-600 text-white rounded-lg shadow-lg hover:bg-green-700 transition-colors flex items-center gap-6">
        <PlusCircleIcon className="h-10 w-10" />
        <div>
          <h2 className="font-bold text-2xl">Create new users</h2>
          <p className="text-base opacity-90">Create new users for the system.</p>
        </div>
      </Link>
    </div>
  );
};

export default QuickActions;