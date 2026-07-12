import React, { useState } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { mockResource, mockTimeline } from '../mockData/bookings';
import { Search } from 'lucide-react';
// import { axiosClient } from '../api/axiosClient';
// import { ENDPOINTS } from '../api/endpoints';

export const BookingsPage = () => {
  const [resourceSearch, setResourceSearch] = useState(`${mockResource.name} - ${mockResource.date}`);
  const [isBooking, setIsBooking] = useState(false);

  /*
  // TODO: Uncomment when backend is ready
  useEffect(() => {
    // Fetch bookings for the selected date and resource to populate the timeline
    const fetchBookings = async () => {
      try {
        const { data } = await axiosClient.get(ENDPOINTS.BOOKINGS.ASSET_BOOKINGS(mockResource.id));
        // Map backend data to UI timeline format here
      } catch (error) {
        console.error('Failed to fetch bookings', error);
      }
    };
    fetchBookings();
  }, []);
  */

  const handleBookSlot = () => {
    setIsBooking(true);
    // Simulate booking attempt
    setTimeout(() => {
      setIsBooking(false);
      alert('Mock Booking Action Triggered!');
    }, 500);
  };

  return (
    <AppLayout>
      <div className="p-8 max-w-5xl mx-auto space-y-8">
        
        {/* Resource Selection */}
        <section className="flex flex-col gap-1.5 max-w-2xl">
          <label className="text-sm font-medium text-gray-400">Resource</label>
          <div className="relative">
            <input 
              type="text" 
              value={resourceSearch}
              onChange={(e) => setResourceSearch(e.target.value)}
              className="w-full px-4 py-2.5 bg-gray-900 border border-gray-700 rounded-lg text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all"
            />
          </div>
        </section>

        {/* Timeline View */}
        <section className="mt-8">
          <div className="relative h-[480px] mt-4">
            
            {/* Grid Background */}
            {mockTimeline.hours.map((hour, index) => (
              <div 
                key={hour} 
                className="absolute w-full flex items-start" 
                style={{ top: `${index * 96}px` }}
              >
                <span className="w-16 text-sm text-gray-400 font-mono -mt-2.5">
                  {hour}
                </span>
                <div className="flex-1 border-t border-gray-800/80 ml-4"></div>
              </div>
            ))}

            {/* Booking Blocks */}
            <div className="absolute left-20 right-0 bottom-0 top-0">
              {mockTimeline.bookings.map((block) => (
                <div
                  key={block.id}
                  className={`absolute left-0 right-4 rounded-lg p-3 text-sm font-medium shadow-sm flex items-center transition-all ${
                    block.type === 'confirmed'
                      ? 'bg-blue-500/20 border border-blue-500/50 text-blue-300 z-10'
                      : 'bg-red-500/10 border-2 border-dashed border-red-500/60 text-red-400 z-20 backdrop-blur-[1px]'
                  }`}
                  style={{
                    top: `${block.topOffset}px`,
                    height: `${block.height}px`,
                  }}
                >
                  {block.title}
                </div>
              ))}
            </div>

          </div>
        </section>

        {/* Action Button */}
        <section className="pt-4">
          <button 
            onClick={handleBookSlot}
            disabled={isBooking}
            className="px-6 py-2.5 bg-emerald-600/10 border border-emerald-600 text-emerald-500 hover:bg-emerald-600 hover:text-white rounded-lg font-medium transition-all shadow-sm disabled:opacity-50"
          >
            {isBooking ? 'Processing...' : 'Book a slot'}
          </button>
        </section>

      </div>
    </AppLayout>
  );
};