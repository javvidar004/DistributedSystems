// components/layout/Navbar.tsx
'use client';
import { Bars3Icon } from '@heroicons/react/24/solid';

// Props interface for the Navbar component
interface NavbarProps {
  onToggleLeftSidebar: () => void;
}

/**
 * The main navigation bar for authenticated users.
 * Contains toggles for opening and closing the sidebars.
 */
const Navbar = ({ onToggleLeftSidebar }: NavbarProps) => {
  return (
    <nav className="bg-dark-background text-white shadow-md w-full z-20 fixed top-0 left-0" >
      <div className="container mx-auto px-4 flex justify-between items-center h-16">
        {/* Left Section: Menu Toggle and Title */}
        <div className="flex items-center space-x-4">
          <button
            onClick={onToggleLeftSidebar}
            className="p-2 rounded-md hover:bg-gray-700 transition-colors"
            aria-label="Toggle menu"
          >
            <Bars3Icon className="h-6 w-6" />
          </button>
          <h1 className="text-xl font-bold">Access Control Manager</h1>
        </div>
        
      </div>
    </nav>
  );
};

export default Navbar;