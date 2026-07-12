import { useState, useEffect } from 'react';
import { AppLayout } from '../components/layout/AppLayout';

import { Download } from 'lucide-react';
import { axiosClient } from '../api/axiosClient';
import { ENDPOINTS } from '../api/endpoints';

export const ReportsPage = () => {
  const [isExporting, setIsExporting] = useState(false);
  const [reportData, setReportData] = useState({
    utilization: [] as { label: string; value: number }[],
    mostUsed: [] as { asset: string; stat: string }[],
    idle: [] as { asset: string; stat: string }[],
    actionNeeded: [] as { asset: string; stat: string }[]
  });
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const fetchReports = async () => {
      setIsLoading(true);
      try {
        const { data } = await axiosClient.get(ENDPOINTS.REPORTS.ANALYTICS);
        if (data) {
          setReportData(data);
        }
      } catch (error) {
        console.error('Failed to fetch reports', error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchReports();
  }, []);

  const handleExport = () => {
    setIsExporting(true);
    setTimeout(() => {
      setIsExporting(false);
      alert('Report exported successfully!');
    }, 1000);
  };

  return (
    <AppLayout>
      <div className="p-8 max-w-6xl mx-auto space-y-10">
        
        {/* Charts Row */}
        {isLoading ? (
          <div className="flex justify-center items-center h-48">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          
          {/* Utilization Bar Chart */}
          <section className="bg-slate-800/40 border border-slate-700/60 rounded-2xl p-6 shadow-sm">
            <h3 className="text-gray-200 font-medium mb-8">Utilization by department</h3>
            {isLoading ? (
              <p className="text-gray-500">Loading chart...</p>
            ) : (
              <div className="h-40 flex items-end gap-4 border-b border-gray-600/50 pb-2">
                {reportData.utilization.map((item, i) => (
                  <div key={i} className="flex-1 flex flex-col items-center gap-2 group">
                    <div 
                      className="w-full bg-amber-600/80 border border-amber-500 rounded-t-sm transition-all duration-300 group-hover:bg-amber-500"
                      style={{ height: `${item.value}%` }}
                    ></div>
                    <span className="text-xs text-gray-400 font-medium">{item.label}</span>
                  </div>
                ))}
              </div>
            )}
          </section>

          {/* Maintenance Line Chart */}
          <section className="bg-slate-800/40 border border-slate-700/60 rounded-2xl p-6 shadow-sm flex flex-col">
            <h3 className="text-gray-200 font-medium mb-8">Maintenance Frequency</h3>
            <div className="flex-1 relative border-b border-gray-600/50 pb-2 flex items-end">
              <svg viewBox="0 0 100 50" className="w-full h-full overflow-visible preserve-3d" preserveAspectRatio="none">
                <polyline 
                  points="0,40 15,25 30,28 45,15 55,22 70,10 85,7" 
                  fill="none" 
                  stroke="#ef4444" 
                  strokeWidth="1.5" 
                  vectorEffect="non-scaling-stroke"
                />
              </svg>
            </div>
          </section>
        </div>

        {/* Lists Row */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 pt-4">
          
          {/* Most Used */}
          <section>
            <h3 className="text-xl font-semibold text-gray-200 mb-4">Most used assets</h3>
            <div className="space-y-3">
              {reportData.mostUsed.map((item, i) => (
                <div key={i} className="flex text-sm">
                  <span className="text-gray-300 font-medium w-36 shrink-0">{item.asset}</span>
                  <span className="text-gray-500">: {item.stat}</span>
                </div>
              ))}
            </div>
          </section>

          {/* Idle Assets */}
          <section>
            <h3 className="text-xl font-semibold text-gray-200 mb-4">Idle assets</h3>
            <div className="space-y-3">
              {reportData.idle.map((item, i) => (
                <div key={i} className="flex text-sm">
                  <span className="text-gray-300 font-medium w-36 shrink-0">{item.asset}</span>
                  <span className="text-gray-500">: {item.stat}</span>
                </div>
              ))}
            </div>
          </section>

        </div>

        <div className="border-t border-gray-800 my-4"></div>

        {/* Actionable Insights & Export */}
        <section className="space-y-6">
          <div>
            <h3 className="text-xl font-semibold text-gray-200 mb-4">Assets due for maintenance / nearing retirement</h3>
            <div className="space-y-3">
              {reportData.actionNeeded.map((item, i) => (
                <div key={i} className="flex text-sm">
                  <span className="text-gray-300 font-medium w-40 shrink-0">{item.asset}</span>
                  <span className="text-gray-500">: {item.stat}</span>
                </div>
              ))}
            </div>
          </div>

          <button 
            onClick={handleExport}
            disabled={isExporting}
            className="flex items-center gap-2 px-6 py-2.5 bg-gray-800/80 border border-gray-700 hover:border-gray-500 text-gray-300 hover:text-white rounded-lg font-medium transition-all shadow-sm disabled:opacity-50 mt-4"
          >
            <Download size={18} className={isExporting ? 'animate-bounce' : ''} />
            {isExporting ? 'Generating...' : 'Export report'}
          </button>
        </section>
        </>
        )}

      </div>
    </AppLayout>
  );
};