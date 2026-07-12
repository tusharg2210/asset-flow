export const mockResource = {
  id: 50,
  name: 'Conference room B2',
  date: 'Tue, 7 Jul',
};

export const mockTimeline = {
  hours: ['9:00', '10:00', '11:00', '12:00', '1:00'],
  // For the UI simulation, we use pixel offsets. 
  // Assuming each hour is 96px (h-24) high.
  bookings: [
    {
      id: 1,
      title: 'Booked - Procurement Team - 9 to 10',
      type: 'confirmed',
      topOffset: 0,      // 9:00 AM
      height: 96,        // 1 hour
    },
    {
      id: 2,
      title: 'Requested 9:30 to 10:30 - conflict - slot is unavailable',
      type: 'conflict',
      topOffset: 48,     // 9:30 AM (halfway through the first hour block)
      height: 96,        // 1 hour
    }
  ]
};