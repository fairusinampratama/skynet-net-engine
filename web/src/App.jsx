import { useState, useEffect } from "react";
import useSWR from "swr";
import { LayoutDashboard, Server, ChevronDown, RefreshCw, Database } from "lucide-react";
import { ResourceWidget } from "./features/dashboard/ResourceWidget";
import { TrafficWidget } from "./features/dashboard/TrafficWidget";
import { CustomerTable } from "./features/customers/CustomerTable";
import { api, syncRouter, backupRouter } from "./api/client";

const fetcher = (url) => api.get(url).then((res) => res.data);

function App() {
  const [selectedRouter, setSelectedRouter] = useState(null);
  const [monitoringUser, setMonitoringUser] = useState(null);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [isSyncing, setIsSyncing] = useState(false);
  const [isBackingUp, setIsBackingUp] = useState(false);

  const handleSync = async () => {
    if (!selectedRouter) return;
    setIsSyncing(true);
    try {
      await syncRouter(selectedRouter.id);
      // Ideally trigger a re-fetch of SWR in child components, but for now just wait a bit
      setTimeout(() => setIsSyncing(false), 1000);
    } catch (e) {
      console.error("Sync failed", e);
      setIsSyncing(false);
    }
  };

  const handleBackup = async () => {
    if (!selectedRouter) return;
    setIsBackingUp(true);
    try {
      const res = await backupRouter(selectedRouter.id);
      alert(`Backup created: ${res.data.file}`); // Simple alert for MVP
    } catch (e) {
      console.error("Backup failed", e);
      alert("Backup failed!");
    } finally {
      setIsBackingUp(false);
    }
  };

  // Fetch Routers
  const { data: routers } = useSWR("/routers", fetcher);

  // Auto-select first router
  useEffect(() => {
    if (routers && routers.length > 0 && !selectedRouter) {
      setSelectedRouter(routers[0]);
    }
  }, [routers]);

  const currentRouterId = selectedRouter ? selectedRouter.id : 1;

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900 pb-20">
      {/* Navbar */}
      <nav className="bg-white border-b border-slate-200 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center text-white shadow-lg shadow-blue-600/20">
              <LayoutDashboard className="w-6 h-6" />
            </div>
            <div>
              <h1 className="text-xl font-bold bg-gradient-to-r from-blue-700 to-indigo-700 bg-clip-text text-transparent">
                Skynet NetEngine
              </h1>
              <p className="text-xs text-slate-500 font-medium">ISP Control Center v1.0</p>
            </div>
          </div>

          {/* Router Selector */}
          <div className="flex items-center gap-4">
            <div className="relative">
              <button
                onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                className="flex items-center gap-2 px-4 py-2 bg-slate-100 hover:bg-slate-200 rounded-lg text-sm font-medium text-slate-700 transition-colors border border-slate-200"
              >
                <Server className="w-4 h-4 text-slate-500" />
                <span>{selectedRouter ? selectedRouter.name : "Loading Router..."}</span>
                <ChevronDown className={`w-4 h-4 text-slate-400 transition-transform ${isDropdownOpen ? 'rotate-180' : ''}`} />
              </button>

              {/* Dropdown */}
              {isDropdownOpen && (
                <div className="absolute right-0 top-full mt-2 w-56 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in fade-in slide-in-from-top-2 z-50">
                  <div className="py-1">
                    {routers && routers.map((router) => (
                      <button
                        key={router.id}
                        onClick={() => {
                          setSelectedRouter(router);
                          setIsDropdownOpen(false);
                        }}
                        className={`w-full text-left px-4 py-2.5 text-sm flex items-center gap-2 hover:bg-blue-50 hover:text-blue-700 transition-colors ${selectedRouter?.id === router.id ? 'bg-blue-50 text-blue-700 font-medium' : 'text-slate-600'}`}
                      >
                        <div className={`w-2 h-2 rounded-full ${selectedRouter?.id === router.id ? 'bg-blue-600' : 'bg-slate-300'}`} />
                        {router.name}
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {/* Action Buttons */}
            {selectedRouter && (
              <div className="flex items-center gap-2">
                <button
                  onClick={handleSync}
                  disabled={isSyncing}
                  className={`p-2 rounded-lg border border-slate-200 bg-white text-slate-500 hover:text-blue-600 hover:border-blue-200 transition-all ${isSyncing ? 'animate-spin text-blue-600' : ''}`}
                  title="Sync Users & Secrets"
                >
                  <RefreshCw className="w-5 h-5" />
                </button>
                <button
                  onClick={handleBackup}
                  disabled={isBackingUp}
                  className={`p-2 rounded-lg border border-slate-200 bg-white text-slate-500 hover:text-blue-600 hover:border-blue-200 transition-all ${isBackingUp ? 'opacity-50 cursor-not-allowed' : ''}`}
                  title="Create Backup"
                >
                  <Database className="w-5 h-5" />
                </button>
              </div>
            )}
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">

        {/* Top Grid: Health & Traffic */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <ResourceWidget routerId={currentRouterId} />
          <TrafficWidget
            routerId={currentRouterId}
            selectedUser={monitoringUser}
            onAutoSelect={setMonitoringUser}
          />
        </div>

        {/* Bottom Grid: Customers */}
        <div className="grid grid-cols-1">
          <CustomerTable
            routerId={currentRouterId}
            onSelectUser={setMonitoringUser}
            selectedUser={monitoringUser}
          />
        </div>

      </main>
    </div>
  );
}

export default App;
