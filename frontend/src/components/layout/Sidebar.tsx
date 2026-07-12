
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Building2, 
  Box, 
  ArrowRightLeft, 
  CalendarCheck, 
  Wrench, 
  ClipboardCheck, 
  BarChart3, 
  Bell,
  LogOut
} from 'lucide-react';
import { useAuthStore } from '../../store/useAuthStore';

const navItems = [
  { name: 'Dashboard', icon: LayoutDashboard, path: '/dashboard' },
  { name: 'Organization setup', icon: Building2, path: '/organization' },
  { name: 'Assets', icon: Box, path: '/assets' },
  { name: 'Allocation & Transfer', icon: ArrowRightLeft, path: '/allocations' },
  { name: 'Resource Booking', icon: CalendarCheck, path: '/bookings' },
  { name: 'Maintenance', icon: Wrench, path: '/maintenance' },
  { name: 'Audit', icon: ClipboardCheck, path: '/audit' },
  { name: 'Reports', icon: BarChart3, path: '/reports' },
  { name: 'Notifications', icon: Bell, path: '/notifications' },
];

export const Sidebar = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const logout = useAuthStore((state) => state.logout);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="w-64 bg-gray-900 border-r border-gray-800 h-screen flex flex-col text-gray-300">
      
      {/* Brand Header */}
      <div className="h-16 flex items-center px-6 border-b border-gray-800 shrink-0">
        <h1 className="text-2xl font-bold text-gray-50 tracking-tight">
          Asset<span className="text-orange-600">Flow</span>
        </h1>
      </div>

      {/* Navigation */}
      <nav className="flex-1 py-6 px-4 space-y-1.5 overflow-y-auto">
        {navItems.map((item) => {
          const Icon = item.icon;
          const isActive = location.pathname.startsWith(item.path);
          
          return (
            <Link
              key={item.name}
              to={item.path}
              className={`w-full flex items-center gap-3 px-4 py-2.5 rounded-lg font-medium transition-colors ${
                isActive 
                  ? 'bg-gray-800 text-orange-500 border border-gray-700 shadow-sm' 
                  : 'text-gray-400 hover:text-gray-100 hover:bg-gray-800/50'
              }`}
            >
              <Icon size={18} className={isActive ? 'text-orange-500' : 'text-gray-400'} />
              {item.name}
            </Link>
          );
        })}
      </nav>

      {/* Logout Button at Bottom */}
      <div className="p-4 border-t border-gray-800 shrink-0">
        <button 
          onClick={handleLogout}
          className="w-full flex items-center gap-3 px-4 py-2.5 rounded-lg font-medium text-gray-400 hover:text-red-400 hover:bg-gray-800/50 transition-colors"
        >
          <LogOut size={18} />
          Log Out
        </button>
      </div>

    </div>
  );
};