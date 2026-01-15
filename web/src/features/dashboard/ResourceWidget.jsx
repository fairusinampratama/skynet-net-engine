import useSWR from "swr";
import { api } from "../../api/client";
import { Cpu, MemoryStick, Activity, Server } from "lucide-react";
import { Skeleton } from "../../components/ui/Skeleton";

const fetcher = (url) => api.get(url).then((res) => res.data);

const StatCard = ({ icon: Icon, label, value, sub, color }) => (
    <div className="flex items-center p-4 bg-slate-50 rounded-lg border border-slate-100">
        <div className={`p-3 rounded-full mr-4 ${color}`}>
            <Icon className="w-6 h-6 text-white" />
        </div>
        <div>
            <p className="text-sm text-slate-500 font-medium">{label}</p>
            <p className="text-xl font-bold text-slate-800">{value}</p>
            {sub && <p className="text-xs text-slate-400">{sub}</p>}
        </div>
    </div>
);

export const ResourceWidget = ({ routerId = 1 }) => {
    const { data, error } = useSWR(`/router/${routerId}/health`, fetcher, {
        refreshInterval: 5000,
    });

    // Treat network errors (500, 502, 503, 504) as loading/reconnecting
    // 503 = Service Unavailable (Warmup)
    // 504 = Gateway Timeout (Slow router)
    // 502 = Bad Gateway (Proxy issue)
    const isTransient = error && [500, 502, 503, 504].includes(error.response?.status);
    const isLoading = !data && (!error || isTransient);
    const isError = error && !isTransient;

    if (isError) return <div className="p-4 bg-red-50 text-red-500 rounded">Failed to load health data</div>;

    if (isLoading) return (
        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
            <div className="flex justify-between mb-4">
                <Skeleton className="h-6 w-32" />
                <Skeleton className="h-6 w-16 rounded-full" />
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {[1, 2, 3, 4].map((i) => (
                    <Skeleton key={i} className="h-24 w-full rounded-lg" />
                ))}
            </div>
        </div>
    );

    const totalMem = (data.total_memory / 1024 / 1024).toFixed(0);
    const freeMem = (data.free_memory / 1024 / 1024).toFixed(0);
    const usedMem = totalMem - freeMem;
    const memPercent = Math.round((usedMem / totalMem) * 100);

    return (
        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
            <h3 className="text-lg font-semibold text-slate-800 mb-4 flex items-center justify-between">
                <span>System Health</span>
                <span className="text-xs bg-green-100 text-green-700 px-2 py-1 rounded-full font-bold">ONLINE</span>
            </h3>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <StatCard
                    icon={Cpu}
                    label="CPU Load"
                    value={`${data.cpu}%`}
                    color="bg-purple-500"
                />
                <StatCard
                    icon={MemoryStick}
                    label="Memory Usage"
                    value={`${memPercent}%`}
                    sub={`${usedMem}MB / ${totalMem}MB`}
                    color="bg-orange-500"
                />
                <StatCard
                    icon={Activity}
                    label="Uptime"
                    value={data.uptime}
                    color="bg-blue-500"
                />
                <StatCard
                    icon={Server}
                    label="Device"
                    value={data.board_name}
                    sub={data.version}
                    color="bg-slate-700"
                />
            </div>
        </div>
    );
};
