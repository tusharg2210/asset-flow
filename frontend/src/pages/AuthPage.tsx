import React, { useState } from 'react';
import { Input } from '../components/common/Input';
import { Button } from '../components/common/Button';
import { axiosClient } from '../api/axiosClient';
import { ENDPOINTS } from '../api/endpoints';
import { useAuthStore } from '../store/useAuthStore';
import type { AuthResponse } from '../types/api';

export const AuthPage = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState(''); // Only for signup
  const [isLoading, setIsLoading] = useState(false);
  
  const setAuth = useAuthStore((state) => state.setAuth);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    
    try {
      const endpoint = isLogin ? ENDPOINTS.AUTH.LOGIN : ENDPOINTS.AUTH.SIGNUP;
      const payload = isLogin ? { email, password } : { name, email, password };
      
      const { data } = await axiosClient.post<AuthResponse>(endpoint, payload);
      
      // If signup is successful, we might need to log them in immediately, 
      // but assuming the API returns the token on signup too based on standard practices
      if (data.token && data.user) {
         setAuth(data.user, data.token);
      } else if (!isLogin) {
         // If signup doesn't return token, switch to login view automatically
         setIsLogin(true);
         alert("Account created! Please log in.");
      }
    } catch (error) {
      console.error('Authentication failed:', error);
      alert('Authentication failed. Please check your credentials.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
      <div className="max-w-md w-full bg-white rounded-2xl shadow-xl border border-slate-100 p-8">
        
        {/* Header & Logo */}
        <div className="flex flex-col items-center mb-8">
          <h1 className="text-2xl font-bold text-slate-900 mb-6 tracking-tight">
            AssetFlow <span className="text-orange-600">– {isLogin ? 'login' : 'signup'}</span>
          </h1>
          <div className="w-16 h-16 rounded-full border-2 border-orange-600 flex items-center justify-center bg-orange-50 text-orange-600 font-bold text-xl shadow-inner">
            AF
          </div>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="flex flex-col gap-5">
          {!isLogin && (
            <Input 
              label="Full Name" 
              type="text" 
              placeholder="John Doe" 
              value={name}
              onChange={(e) => setName(e.target.value)}
              required 
            />
          )}
          <Input 
            label="Email" 
            type="email" 
            placeholder="name@company.com" 
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required 
          />
          <div className="flex flex-col gap-1">
            <Input 
              label="Password" 
              type="password" 
              placeholder="••••••••" 
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required 
            />
            {isLogin && (
              <div className="flex justify-end mt-1">
                <button type="button" className="text-sm font-medium text-orange-600 hover:text-orange-700">
                  Forgot password?
                </button>
              </div>
            )}
          </div>

          <Button type="submit" isLoading={isLoading}>
            {isLogin ? 'Log In' : 'Create Account'}
          </Button>
        </form>

        {/* Divider & Switcher */}
        <div className="mt-8 border-t border-slate-200 pt-6">
          <div className="text-sm text-slate-600 mb-4 font-medium">
            {isLogin ? 'New here?' : 'Already have an account?'}
          </div>
          
          {isLogin && (
            <div className="bg-slate-50 p-4 rounded-lg border border-slate-200 mb-4">
              <p className="text-sm text-slate-600">
                Sign up creates an <span className="font-semibold text-slate-800">employee account</span>. 
                Admin roles are assigned later.
              </p>
            </div>
          )}

          <Button 
            variant="outline" 
            type="button"
            onClick={() => setIsLogin(!isLogin)}
          >
            {isLogin ? 'Create Account' : 'Back to Login'}
          </Button>
        </div>

      </div>
    </div>
  );
};