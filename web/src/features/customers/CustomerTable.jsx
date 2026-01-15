import { useState } from "react";
import useSWR from "swr";
import { api } from "../../api/client";
import { Users, Wifi, AlertTriangle, ChevronLeft, ChevronRight } from "lucide-react";
import { Skeleton } from "../../components/ui/Skeleton";

const fetcher = (url) => api.get(url).then((res) => res.data);

export const CustomerTable = ({ onSelectUser, selectedUser }) => {
    const [page, setPage] = useState(1);
    const PAGE_SIZE = 10;

    const { data, error } = useSWR("/monitoring/targets", fetcher, {
        refreshInterval: 10000,
    });

    const isTransient = error && [500, 502, 503, 504].includes(error.response?.status);
    const isLoading = !data && (!error || isTransient);

    // Pagination Logic
    const safeData = data || [];
    const totalPages = Math.ceil(safeData.length / PAGE_SIZE);
    const startIndex = (page - 1) * PAGE_SIZE;
    const paginatedData = safeData.slice(startIndex, startIndex + PAGE_SIZE);

    return (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
            <div className="p-6 border-b border-slate-100 flex items-center justify-between">
                <h3 className="text-lg font-semibold text-slate-800 flex items-center gap-2">
                    <Users className="w-5 h-5 text-slate-500" />
                    Active Sessions
                </h3>
                <span className="text-sm font-medium text-slate-500">
                    Total: {safeData.length}
                </span>
            </div>

            <div className="overflow-x-auto">
                <table className="w-full text-left">
                    <thead className="bg-slate-50 text-slate-500 text-xs uppercase tracking-wider">
                        <tr>
                            <th className="px-6 py-4 font-semibold">User</th>
                            <th className="px-6 py-4 font-semibold">IP Address</th>
                            <th className="px-6 py-4 font-semibold">Status</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                        {isLoading && Array.from({ length: 5 }).map((_, i) => (
                            <tr key={i}>
                                <td className="px-6 py-4"><Skeleton className="h-4 w-32" /></td>
                                <td className="px-6 py-4"><Skeleton className="h-4 w-24" /></td>
                                <td className="px-6 py-4"><Skeleton className="h-6 w-20 rounded-full" /></td>
                            </tr>
                        ))}

                        {paginatedData.map((user, i) => (
                            <tr
                                key={i}
                                onClick={() => onSelectUser && onSelectUser(user.name)}
                                className={`cursor-pointer transition-colors ${selectedUser === user.name ? 'bg-blue-50' : 'hover:bg-slate-50'}`}
                            >
                                <td className="px-6 py-4 font-medium text-slate-900 flex items-center gap-2">
                                    {selectedUser === user.name && <div className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />}
                                    {user.name}
                                </td>
                                <td className="px-6 py-4 text-slate-600 font-mono text-sm">{user.address}</td>
                                <td className="px-6 py-4">
                                    <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 border border-green-200">
                                        <Wifi className="w-3 h-3" />
                                        Connected
                                    </span>
                                </td>
                            </tr>
                        ))}

                        {!isLoading && safeData.length === 0 && (
                            <tr><td colSpan="3" className="px-6 py-8 text-center text-slate-400">No active users found.</td></tr>
                        )}
                    </tbody>
                </table>
            </div>

            {/* Pagination Controls */}
            {totalPages > 1 && (
                <div className="p-4 border-t border-slate-100 flex items-center justify-between bg-slate-50">
                    <button
                        onClick={() => setPage(p => Math.max(1, p - 1))}
                        disabled={page === 1}
                        className="p-2 rounded-lg hover:bg-white hover:shadow-sm disabled:opacity-50 disabled:hover:shadow-none transition-all"
                    >
                        <ChevronLeft className="w-5 h-5 text-slate-600" />
                    </button>
                    <span className="text-sm font-medium text-slate-600">
                        Page {page} of {totalPages}
                    </span>
                    <button
                        onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                        disabled={page === totalPages}
                        className="p-2 rounded-lg hover:bg-white hover:shadow-sm disabled:opacity-50 disabled:hover:shadow-none transition-all"
                    >
                        <ChevronRight className="w-5 h-5 text-slate-600" />
                    </button>
                </div>
            )}
        </div>
    );
};
