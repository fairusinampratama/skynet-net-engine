import { useState, useEffect } from "react";
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";
import { Skeleton } from "../../components/ui/Skeleton";
import { api } from "../../api/client";
import useSWR from "swr";

const fetcher = (url) => api.get(url).then((res) => res.data);

export const TrafficWidget = ({ routerId = 1, selectedUser, onAutoSelect }) => {
    const [history, setHistory] = useState([]);

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
            setHistory((prev) => {
                const now = new Date().toLocaleTimeString();
                const pointA = {
                    time: now,
                    tx: (parseInt(traffic.tx || "0") / 1000000).toFixed(2), // bits to Mbps
                    rx: (parseInt(traffic.rx || "0") / 1000000).toFixed(2), // bits to Mbps
                };
                return [...prev.slice(-59), pointA];
            });
        }
        // Reset history when user changes
    }, [traffic]);

    // Reset history when target changes
    useEffect(() => {
        setHistory([]);
    }, [target]);

    // Treat network errors as loading
    const isTransient = (errorTraffic || errorUsers) && [500, 502, 503, 504].includes((errorTraffic || errorUsers)?.response?.status);
    const isLoading = (!users || (!traffic && !target)) && (!errorTraffic || isTransient);

    if (isLoading) {
        return (
            <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
                <div className="flex justify-between mb-4">
                    <Skeleton className="h-6 w-40" />
                    <Skeleton className="h-4 w-32" />
                </div>
                <Skeleton className="h-64 w-full rounded-lg" />
            </div>
        );
    }

    return (
        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
            <h3 className="text-lg font-semibold text-slate-800 mb-4 flex justify-between">
                <span>Live Traffic (Mbps)</span>
                {target && <span className="text-sm font-normal text-slate-500">Monitoring: {target}</span>}
            </h3>
            <div className="h-64">
                {history.length === 0 ? (
                    <div className="space-y-3">
                        <Skeleton className="h-8 w-full" />
                        <Skeleton className="h-8 w-full" />
                        <Skeleton className="h-8 w-full" />
                        <Skeleton className="h-8 w-full" />
                        <Skeleton className="h-8 w-full" />
                        <Skeleton className="h-8 w-full" />
                        <div className="text-center mt-4">
                            <div className="text-xs text-slate-400">Loading traffic data...</div>
                        </div>
                    </div>
                ) : (
                    <ResponsiveContainer width="100%" height="100%">
                        <AreaChart data={history}>
                            <CartesianGrid strokeDasharray="3 3" opacity={0.3} />
                            <XAxis dataKey="time" hide />
                            <YAxis />
                            <Tooltip />
                            <Area type="monotone" dataKey="rx" stackId="1" stroke="#3b82f6" fill="#3b82f6" fillOpacity={0.2} name="Download" />
                            <Area type="monotone" dataKey="tx" stackId="1" stroke="#10b981" fill="#10b981" fillOpacity={0.2} name="Upload" />
                        </AreaChart>
                    </ResponsiveContainer>
                )}
            </div>
        </div>
    );
};
