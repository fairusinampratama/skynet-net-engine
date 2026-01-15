import { useState } from "react";
import useSWR from "swr";
import { api } from "../../api/client";
import { Wifi, ChevronLeft, ChevronRight, ArrowUpDown, ArrowUp, ArrowDown, Users } from "lucide-react";
import { Skeleton } from "../../components/ui/Skeleton";

const fetcher = (url) => api.get(url).then((res) => res.data);

export const CustomerTable = ({ routerId = 1, onSelectUser, selectedUser }) => {
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [sortColumn, setSortColumn] = useState('username');
    const [sortDirection, setSortDirection] = useState('asc');

    const { data, error, isLoading } = useSWR(`/router/${routerId}/users`, fetcher, {
        refreshInterval: 30000,
        revalidateOnFocus: false,
    });

    // Sorting Logic
    const safeData = data || [];
    const sortedData = [...safeData].sort((a, b) => {
        let aVal = a[sortColumn] || '';
        let bVal = b[sortColumn] || '';

        // Handle status sorting (connected first)
        if (sortColumn === 'status') {
            aVal = a.status === 'connected' ? '0' : '1';
            bVal = b.status === 'connected' ? '0' : '1';
        }

        if (aVal < bVal) return sortDirection === 'asc' ? -1 : 1;
        if (aVal > bVal) return sortDirection === 'asc' ? 1 : -1;
        return 0;
    });

    // Pagination Logic
    const totalPages = Math.ceil(sortedData.length / pageSize);
    const startIdx = (page - 1) * pageSize;
    const endIdx = Math.min(startIdx + pageSize, sortedData.length);
    const paginatedData = sortedData.slice(startIdx, endIdx);

    // Sort handler
    const handleSort = (column) => {
        if (sortColumn === column) {
            setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
        } else {
            setSortColumn(column);
            setSortDirection('asc');
        }
    };

    // Page size handler
    const handlePageSizeChange = (newSize) => {
        setPageSize(newSize);
        setPage(1); // Reset to first page
    };

    return (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
            <div className="p-6 border-b border-slate-100 flex items-center justify-between">
                <h3 className="text-lg font-semibold text-slate-800 flex items-center gap-2">
                    <Users className="w-5 h-5 text-slate-500" />
                    Active Sessions
                </h3>
                <span className="text-sm font-medium text-slate-500">
                    Total: {sortedData.length}
                </span>
            </div>

            <div className="overflow-x-auto">
                <table className="w-full text-left">
                    <thead className="bg-slate-50 border-b border-slate-200">
                        <tr>
                            <th
                                onClick={() => handleSort('username')}
                                className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider cursor-pointer hover:bg-slate-100 transition-colors"
                            >
                                <div className="flex items-center gap-1">
                                    User
                                    {sortColumn === 'username' ? (
                                        sortDirection === 'asc' ? <ArrowUp className="w-3 h-3" /> : <ArrowDown className="w-3 h-3" />
                                    ) : (
                                        <ArrowUpDown className="w-3 h-3 opacity-30" />
                                    )}
                                </div>
                            </th>
                            <th
                                onClick={() => handleSort('ip')}
                                className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider cursor-pointer hover:bg-slate-100 transition-colors"
                            >
                                <div className="flex items-center gap-1">
                                    IP Address
                                    {sortColumn === 'ip' ? (
                                        sortDirection === 'asc' ? <ArrowUp className="w-3 h-3" /> : <ArrowDown className="w-3 h-3" />
                                    ) : (
                                        <ArrowUpDown className="w-3 h-3 opacity-30" />
                                    )}
                                </div>
                            </th>
                            <th
                                onClick={() => handleSort('status')}
                                className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider cursor-pointer hover:bg-slate-100 transition-colors"
                            >
                                <div className="flex items-center gap-1">
                                    Status
                                    {sortColumn === 'status' ? (
                                        sortDirection === 'asc' ? <ArrowUp className="w-3 h-3" /> : <ArrowDown className="w-3 h-3" />
                                    ) : (
                                        <ArrowUpDown className="w-3 h-3 opacity-30" />
                                    )}
                                </div>
                            </th>
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
                                onClick={() => onSelectUser && onSelectUser(user.username)}
                                className={`cursor-pointer transition-colors ${selectedUser === user.username ? 'bg-blue-50' : 'hover:bg-slate-50'}`}
                            >
                                <td className="px-6 py-4 font-medium text-slate-900 flex items-center gap-2">
                                    {selectedUser === user.username && <div className="w-2 h-2 rounded-full bg-blue-500 animate-pulse" />}
                                    {user.username}
                                </td>
                                <td className="px-6 py-4 text-slate-600 font-mono text-sm">{user.ip || '-'}</td>
                                <td className="px-6 py-4">
                                    {user.status === 'connected' ? (
                                        <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 border border-green-200">
                                            <Wifi className="w-3 h-3" />
                                            Connected
                                        </span>
                                    ) : (
                                        <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-600 border border-slate-200">
                                            <div className="w-2 h-2 rounded-full bg-slate-400" />
                                            Offline
                                        </span>
                                    )}
                                </td>
                            </tr>
                        ))}

                        {!isLoading && sortedData.length === 0 && (
                            <tr><td colSpan="3" className="px-6 py-8 text-center text-slate-400">No active users found.</td></tr>
                        )}
                    </tbody>
                </table>
            </div>

            {/* Advanced Pagination Controls */}
            <div className="flex items-center justify-between px-6 py-4 border-t border-slate-200 bg-slate-50">
                <div className="flex items-center gap-3">
                    <span className="text-sm text-slate-600 font-medium">
                        Showing <span className="text-slate-900">{startIdx + 1}-{endIdx}</span> of <span className="text-slate-900">{sortedData.length}</span>
                    </span>
                    <div className="flex items-center gap-2">
                        <span className="text-sm text-slate-500">Show:</span>
                        <select
                            value={pageSize}
                            onChange={(e) => handlePageSizeChange(Number(e.target.value))}
                            className="px-3 py-1.5 text-sm font-medium border border-slate-300 rounded-lg bg-white hover:border-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent cursor-pointer transition-colors"
                        >
                            <option value={10}>10</option>
                            <option value={25}>25</option>
                            <option value={50}>50</option>
                            <option value={100}>100</option>
                        </select>
                        <span className="text-sm text-slate-500">entries</span>
                    </div>
                </div>

                <div className="flex items-center gap-2">
                    <button
                        onClick={() => setPage(1)}
                        disabled={page === 1}
                        className="px-2 py-1 text-sm font-medium text-slate-600 hover:text-slate-900 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        ««
                    </button>
                    <button
                        onClick={() => setPage(p => Math.max(1, p - 1))}
                        disabled={page === 1}
                        className="px-3 py-1 text-sm font-medium text-slate-600 hover:text-slate-900 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1"
                    >
                        <ChevronLeft className="w-4 h-4" />
                        Previous
                    </button>

                    <span className="text-sm text-slate-600 px-2">
                        Page {page} of {totalPages}
                    </span>

                    <button
                        onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                        disabled={page === totalPages}
                        className="px-3 py-1 text-sm font-medium text-slate-600 hover:text-slate-900 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1"
                    >
                        Next
                        <ChevronRight className="w-4 h-4" />
                    </button>
                    <button
                        onClick={() => setPage(totalPages)}
                        disabled={page === totalPages}
                        className="px-2 py-1 text-sm font-medium text-slate-600 hover:text-slate-900 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        »»
                    </button>
                </div>
            </div>
        </div>
    );
};
