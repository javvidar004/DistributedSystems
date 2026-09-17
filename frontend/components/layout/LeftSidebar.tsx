// components/layout/LeftSidebar.tsx
'use client';
import Link from 'next/link';
import {
  HomeIcon,
  CalendarDaysIcon,
  ShoppingCartIcon,
  BookOpenIcon,
  MagnifyingGlassIcon,
  HeartIcon,
  Cog6ToothIcon,
  XMarkIcon,
} from '@heroicons/react/24/solid';

// Props interface defines the component's expected properties
interface LeftSidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

// Sidebar navigation items
const navItems = [
  { href: '/dashboard', label: 'Home', icon: HomeIcon },
  { href: '/logs', label: 'Logs', icon: ShoppingCartIcon },
  { href: '/users', label: 'Users', icon: BookOpenIcon },
];

/**
 * The left sidebar component for navigation.
 * It is collapsible and contains links to the main pages.
 */
const LeftSidebar = ({ isOpen, onClose }: LeftSidebarProps) => {
  return (
    <>
      {/* Overlay for mobile view */}
      <div
        className={`fixed inset-0 bg-black bg-opacity-50 z-30 transition-opacity md:hidden ${
          isOpen ? 'opacity-100' : 'opacity-0 pointer-events-none'
        }`}
        onClick={onClose}
      ></div>

      {/* Sidebar */}
      <aside
        className={`fixed top-0 left-0 h-full w-64 bg-dark-background text-white transform transition-transform z-40 ${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        <div className="flex flex-col h-full">
          {/* Sidebar Header */}
          <div className="flex items-center justify-between p-4 h-16 border-b border-gray-700">
            <h2 className="text-lg font-bold">Options</h2>
            <button onClick={onClose} className="p-1 rounded-md hover:bg-gray-700"> 
              <XMarkIcon className="h-6 w-6" />
            </button>
          </div>

          {/* Navigation Links */}
          <nav className="flex-grow p-4 space-y-2">
            {navItems.map(({ href, label, icon: Icon }) => (
              <Link
                key={href}
                href={href}
                onClick={onClose}
                className="flex items-center space-x-3 px-3 py-2 rounded-md text-gray-300 hover:bg-secondary hover:text-white transition-colors"
              >
                <Icon className="h-5 w-5" />
                <span>{label}</span>
              </Link>
            ))}
          </nav>
          
          {/* Footer Link (Configurations) */}
          <div className="p-4 border-t border-gray-700">
            <Link
              href="/configurations"
              onClick={onClose}
              className="flex items-center space-x-3 px-3 py-2 rounded-md text-gray-300 hover:bg-secondary hover:text-white transition-colors"
            >
              <Cog6ToothIcon className="h-5 w-5" />
              <span>Configurations</span>
            </Link>
          </div>
        </div>
      </aside>
    </>
  );
};

export default LeftSidebar;