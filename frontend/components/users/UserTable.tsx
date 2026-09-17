'use client';

import { User } from '@/types';

interface UserTableProps {
  users: User[];
}

const formatSalary = (salary: number) =>
  new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(salary);

const UserTable = ({ users }: UserTableProps) => {
  return (
    <div className="bg-white p-4 rounded-lg shadow-lg overflow-hidden">
      <h2 className="font-bold text-lg text-dark-background mb-4">Users</h2>
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Salary</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {users.length === 0 ? (
              <tr>
                <td className="px-6 py-4 whitespace-nowrap" colSpan={4}>
                  No users available.
                </td>
              </tr>
            ) : (
              users.map((user) => (
                <tr key={user.id}>
                  <td className="px-6 py-4 whitespace-nowrap font-medium text-gray-900">
                    {user.name} {user.last_name}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-700">{user.email}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-700">{user.work_position}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-700">{formatSalary(user.salary)}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default UserTable;