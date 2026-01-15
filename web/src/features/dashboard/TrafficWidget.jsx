import { useState, useEffect, useMemo } from "react";
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";
import { Skeleton } from "../../components/ui/Skeleton";
import { ArrowDown, ArrowUp, Activity } from "lucide-react";
import { api } from "../../api/client";
import useSWR from "swr";

const fetcher = (url) => api.get(url).then((res) => res.data);

// Custom Tooltip Component
const CustomTooltip = ({ active, payload }) => {
    if (!active || !payload || !payload.length) return null;

    const data = payload[0].payload;
    return (
        <div className="bg-white px-3 py-2 rounded-lg shadow-lg border border-slate-200">
            <p className="text-xs text-slate-500 mb-1">{data.time}</p>
            <div className="flex flex-col gap-1">
                <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full bg-blue-500"></div>
                    <span className="text-sm font-medium text-slate-700">Download: {data.rx} Mbps</span>
                </div>
                <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full bg-green-500"></div>
                    <span className="text-sm font-medium text-slate-700">Upload: {data.tx} Mbps</span>
                </div>
            </div>
        </div>
    );
};

export const TrafficWidget = ({ routerId = 1, selectedUser, onAutoSelect }) => {
    const [history, setHistory] = useState([]);
    const [peakRx, setPeakRx] = useState(0);
    const [peakTx, setPeakTx] = useState(0);

    // 1. Get first active user to monitor (Auto-Discovery)
    const { data: users, error: errorUsers } = useSWR("/monitoring/targets", fetcher);

    useEffect(() => {
        // If no user is selected externally, auto-select the first one
        if (users && users.length > 0 && !selectedUser && onAutoSelect) {
            onAutoSelect(users[0].name);
        }
    }, [users, selectedUser, onAutoSelect]);

    const target = selectedUser || (users && users.length > 0 ? users[0].name : null);

    // 2. Poll traffic for that user
    const { data: traffic, error: errorTraffic } = useSWR(
        target ? `/router/${routerId}/traffic?user=${target}` : null,
        fetcher,
        {
            refreshInterval: 1000,
            dedupingInterval: 500,
        }
    );

    useEffect(() => {
        if (traffic) {
            const rxMbps = parseFloat((parseInt(traffic.rx || "0") / 1000000).toFixed(2));
            const txMbps = parseFloat((parseInt(traffic.tx || "0") / 1000000).toFixed(2));

            setHistory((prev) => {
                const now = new Date().toLocaleTimeString();
                const point = {
                    time: now,
                    tx: txMbps,
                    rx: rxMbps,
                };
                return [...prev.slice(-59), point];
            });

            // Update peaks
            setPeakRx(prev => Math.max(prev, rxMbps));
            setPeakTx(prev => Math.max(prev, txMbps));
        }
    }, [traffic]);

    // Reset history and peaks when target changes
    useEffect(() => {
        setHistory([]);
        setPeakRx(0);
        setPeakTx(0);
    }, [target]);

    // Current speeds (last data point)
    const currentRx = useMemo(() => history.length > 0 ? history[history.length - 1].rx : 0, [history]);
    const currentTx = useMemo(() => history.length > 0 ? history[history.length - 1].tx : 0, [history]);

    // Treat network errors as loading
    const isTransient = (errorTraffic || errorUsers) && [500, 502, 503, 504].includes((errorTraffic || errorUsers)?.response?.status);
    const isLoading = (!users || (!traffic && !target)) && (!errorTraffic || isTransient);

    if (isLoading) {
        return (
            <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
                {/* Header skeleton */}
                <div className="mb-4">
                    <div className="flex items-center justify-between mb-3">
                        <Skeleton className="h-6 w-32" />
                        <Skeleton className="h-5 w-40" />
                    </div>

                    {/* Speed badges skeleton */}
                    <div className="grid grid-cols-2 gap-3 mb-4">
                        <div className="bg-gradient-to-br from-slate-50 to-slate-100 border border-slate-200 rounded-lg p-3">
                            <Skeleton className="h-4 w-20 mb-2" />
                            <Skeleton className="h-8 w-24 mb-1" />
                            <Skeleton className="h-3 w-16" />
                        </div>
                        <div className="bg-gradient-to-br from-slate-50 to-slate-100 border border-slate-200 rounded-lg p-3">
                            <Skeleton className="h-4 w-20 mb-2" />
                            <Skeleton className="h-8 w-24 mb-1" />
                            <Skeleton className="h-3 w-16" />
                        </div>
                    </div>
                </div>

                {/* Chart skeleton */}
                <Skeleton className="h-64 w-full rounded-lg" />
            </div>
        );
    }

    return (
        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
            {/* Header with current speeds */}
            <div className="mb-4">
                <div className="flex items-center justify-between mb-3">
                    <h3 className="text-lg font-semibold text-slate-800 flex items-center gap-2">
                        <Activity className="w-5 h-5 text-blue-500" />
                        Live Traffic
                    </h3>
                    {target && <span className="text-xs font-medium text-slate-500 bg-slate-100 px-2 py-1 rounded">Monitoring: {target}</span>}
                </div>

                {/* Real-time Speed Badges */}
                <div className="grid grid-cols-2 gap-3 mb-4">
                    <div className="bg-gradient-to-br from-blue-50 to-blue-100 border border-blue-200 rounded-lg p-3 flex flex-col">
                        <div className="flex items-center gap-2 mb-2">
                            <ArrowDown className="w-4 h-4 text-blue-600" />
                            <span className="text-xs font-medium text-blue-700">Download</span>
                        </div>
                        <div className="text-2xl font-bold text-blue-900 mb-1">{currentRx.toFixed(2)} <span className="text-sm font-normal text-blue-600">Mbps</span></div>
                        <div className="text-xs text-blue-600">Peak: {peakRx.toFixed(2)} Mbps</div>
                    </div>

                    <div className="bg-gradient-to-br from-green-50 to-green-100 border border-green-200 rounded-lg p-3 flex flex-col">
                        <div className="flex items-center gap-2 mb-2">
                            <ArrowUp className="w-4 h-4 text-green-600" />
                            <span className="text-xs font-medium text-green-700">Upload</span>
                        </div>
                        <div className="text-2xl font-bold text-green-900 mb-1">{currentTx.toFixed(2)} <span className="text-sm font-normal text-green-600">Mbps</span></div>
                        <div className="text-xs text-green-600">Peak: {peakTx.toFixed(2)} Mbps</div>
                    </div>
                </div>
            </div>

            {/* Chart */}
            <div className="h-64">
                {history.length === 0 ? (
                    <div className="flex items-center justify-center h-full">
                        <div className="text-center">
                            <div className="inline-flex items-center gap-2 px-4 py-2 bg-slate-100 rounded-lg mb-2">
                                <div className="w-2 h-2 rounded-full bg-blue-500 animate-pulse"></div>
                                <span className="text-sm font-medium text-slate-600">Collecting traffic data...</span>
                            </div>
                            <p className="text-xs text-slate-400">Graph will appear in ~1 second</p>
                        </div>
                    </div>
                ) : (
                    <ResponsiveContainer width="100%" height="100%">
                        <AreaChart data={history}>
                            <defs>
                                <linearGradient id="colorRx" x1="0" y1="0" x2="0" y2="1">
                                    <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3} />
                                    <stop offset="95%" stopColor="#3b82f6" stopOpacity={0.05} />
                                </linearGradient>
                                <linearGradient id="colorTx" x1="0" y1="0" x2="0" y2="1">
                                    <stop offset="5%" stopColor="#10b981" stopOpacity={0.3} />
                                    <stop offset="95%" stopColor="#10b981" stopOpacity={0.05} />
                                </linearGradient>
                            </defs>
                            <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" opacity={0.5} />
                            <XAxis
                                dataKey="time"
                                tick={{ fontSize: 11, fill: '#64748b' }}
                                tickLine={false}
                                axisLine={{ stroke: '#e2e8f0' }}
                            />
                            <YAxis
                                tick={{ fontSize: 11, fill: '#64748b' }}
                                tickLine={false}
                                axisLine={{ stroke: '#e2e8f0' }}
                                label={{ value: 'Mbps', angle: -90, position: 'insideLeft', style: { fontSize: 11, fill: '#64748b' } }}
                            />
                            <Tooltip content={<CustomTooltip />} />
                            <Area
                                type="monotone"
                                dataKey="rx"
                                stroke="#3b82f6"
                                strokeWidth={2}
                                fill="url(#colorRx)"
                                name="Download"
                                animationDuration={300}
                            />
                            <Area
                                type="monotone"
                                dataKey="tx"
                                stroke="#10b981"
                                strokeWidth={2}
                                fill="url(#colorTx)"
                                name="Upload"
                                animationDuration={300}
                            />
                        </AreaChart>
                    </ResponsiveContainer>
                )}
            </div>
        </div>
    );
};
