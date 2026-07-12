import React, { useState, useEffect } from 'react';
import { AppLayout } from '../components/layout/AppLayout';

import { axiosClient } from '../api/axiosClient';
import { ENDPOINTS } from '../api/endpoints';
import { Modal } from '../components/common/Modal';
import { FormInput } from '../components/common/FormInput';
import { Button } from '../components/common/Button';

export const BookingsPage = () => {
  const [resources, setResources] = useState<any[]>([]);
  const [selectedAssetId, setSelectedAssetId] = useState('');
  const [timelineBookings, setTimelineBookings] = useState<any[]>([]);
  const [isBookingModalOpen, setIsBookingModalOpen] = useState(false);
  const [isBooking, setIsBooking] = useState(false);

  // Form states for booking
  const [bookStart, setBookStart] = useState('');
  const [bookEnd, setBookEnd] = useState('');

  const fetchResources = async () => {
    try {
      const { data } = await axiosClient.get(ENDPOINTS.ASSETS.DIRECTORY);
      const assetsList = data?.data || data || [];
      const bookable = assetsList.filter((a: any) => a.is_bookable);
      setResources(bookable);
      if (bookable.length > 0 && !selectedAssetId) {
        setSelectedAssetId(bookable[0].id.toString());
      }
    } catch (error) {
      console.error('Failed to fetch resources', error);
    }
  };

  const fetchBookings = async () => {
    if (!selectedAssetId) return;
    try {
      const { data } = await axiosClient.get(ENDPOINTS.BOOKINGS.ASSET_BOOKINGS(selectedAssetId));
      const slots = data || [];
      
      // Map backend slots to timeline items starting from 9:00 AM
      const mapped = slots.map((slot: any, idx: number) => {
        const start = new Date(slot.start_time);
        const end = new Date(slot.end_time);
        
        const startHours = start.getHours();
        const startMins = start.getMinutes();
        const durationMins = (end.getTime() - start.getTime()) / (60 * 1000);
        
        // Calculate offset (each hour is 96px high, i.e., 1.6px per minute)
        // Chart starts at 8:00 AM
        const topOffset = ((startHours - 8) * 60 + startMins) * 1.6;
        const height = durationMins * 1.6;
        
        return {
          id: slot.booking_id || idx,
          title: `Booked - Slot from ${start.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})} to ${end.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}`,
          type: slot.booking_status === 'CANCELLED' ? 'conflict' : 'confirmed',
          topOffset: topOffset < 0 ? 0 : topOffset,
          height: height <= 0 ? 96 : height,
        };
      });
      setTimelineBookings(mapped);
    } catch (error) {
      console.error('Failed to fetch bookings', error);
    }
  };

  useEffect(() => {
    fetchResources();
  }, []);

  useEffect(() => {
    fetchBookings();
  }, [selectedAssetId]);

  const handleBookSlotSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsBooking(true);
    try {
      const payload = {
        asset_id: parseInt(selectedAssetId, 10),
        start_time: new Date(bookStart).toISOString(),
        end_time: new Date(bookEnd).toISOString()
      };
      await axiosClient.post(ENDPOINTS.BOOKINGS.CREATE, payload);
      alert('Booking confirmed!');
      setIsBookingModalOpen(false);
      setBookStart('');
      setBookEnd('');
      fetchBookings();
    } catch (error) {
      console.error('Failed to book resource:', error);
      alert('Failed to book resource (overlap conflict or server error).');
    } finally {
      setIsBooking(false);
    }
  };

  const getSelectedResourceName = () => {
    const found = resources.find(r => r.id.toString() === selectedAssetId);
    return found ? `${found.tag} - ${found.name}` : 'Loading resource...';
  };

  return (
    <AppLayout>
      <div className="p-8 max-w-5xl mx-auto space-y-8">
        
        {/* Resource Selection */}
        <section className="flex flex-col gap-1.5 max-w-2xl">
          <label className="text-sm font-medium text-gray-400">Select Bookable Resource</label>
          <select 
            value={selectedAssetId}
            onChange={(e) => setSelectedAssetId(e.target.value)}
            className="w-full px-4 py-2.5 bg-gray-900 border border-gray-700 rounded-lg text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all appearance-none"
          >
            {resources.map(res => (
              <option key={res.id} value={res.id}>{res.tag} - {res.name}</option>
            ))}
            {resources.length === 0 && (
              <option disabled>No bookable resources configured.</option>
            )}
          </select>
        </section>

        {/* Timeline View */}
        <section className="mt-8">
          <h3 className="text-sm font-semibold text-gray-300">
            Schedule for {getSelectedResourceName()}
          </h3>
          <div className="relative h-[480px] mt-6 border border-gray-800 rounded-xl bg-gray-900/20 overflow-y-auto">
            <div className="relative min-h-[1056px] p-4">
            
            {/* Grid Background */}
            {['8:00 AM', '9:00 AM', '10:00 AM', '11:00 AM', '12:00 PM', '1:00 PM', '2:00 PM', '3:00 PM', '4:00 PM', '5:00 PM', '6:00 PM'].map((hour, index) => (
              <div 
                key={hour} 
                className="absolute w-full flex items-start" 
                style={{ top: `${index * 96 + 16}px` }}
              >
                <span className="w-20 text-sm text-gray-400 font-mono -mt-2.5">
                  {hour}
                </span>
                <div className="flex-1 border-t border-gray-800/80 ml-2"></div>
              </div>
            ))}

            {/* Booking Blocks */}
            <div className="absolute left-20 right-0 bottom-0 top-0 mt-4">
              {timelineBookings.length === 0 ? (
                <div className="text-gray-500 text-sm mt-4 italic">No active bookings scheduled in this window.</div>
              ) : (
                timelineBookings.map((block) => (
                  <div
                    key={block.id}
                    className={`absolute left-0 right-4 rounded-lg p-3 text-sm font-medium shadow-sm flex items-center transition-all ${
                      block.type === 'confirmed'
                        ? 'bg-blue-500/20 border border-blue-500/50 text-blue-300 z-10'
                        : 'bg-red-500/10 border border-red-500/30 text-red-400 z-20 backdrop-blur-[1px]'
                    }`}
                    style={{
                      top: `${block.topOffset}px`,
                      height: `${block.height}px`,
                    }}
                  >
                    {block.title}
                  </div>
                ))
              )}
            </div>

            </div>
          </div>
        </section>

        {/* Action Button */}
        <section className="pt-4">
          <button 
            onClick={() => setIsBookingModalOpen(true)}
            disabled={resources.length === 0}
            className="px-6 py-2.5 bg-emerald-600/10 border border-emerald-600 text-emerald-500 hover:bg-emerald-600 hover:text-white rounded-lg font-medium transition-all shadow-sm disabled:opacity-50"
          >
            Book a slot
          </button>
        </section>

      </div>

      <Modal isOpen={isBookingModalOpen} onClose={() => setIsBookingModalOpen(false)} title={`Book ${getSelectedResourceName()}`}>
        <form onSubmit={handleBookSlotSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <FormInput 
              label="Start Time" 
              type="datetime-local" 
              value={bookStart}
              onChange={(e) => setBookStart(e.target.value)}
              required 
            />
            <FormInput 
              label="End Time" 
              type="datetime-local" 
              value={bookEnd}
              onChange={(e) => setBookEnd(e.target.value)}
              required 
            />
          </div>
          <div className="pt-4">
            <Button type="submit" isLoading={isBooking}>Confirm Booking</Button>
          </div>
        </form>
      </Modal>
    </AppLayout>
  );
};